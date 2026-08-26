#!/usr/bin/env node

import { access, readdir, readFile } from "node:fs/promises";
import { constants } from "node:fs";
import { createRequire } from "node:module";
import path from "node:path";
import process from "node:process";
import GithubSlugger from "github-slugger";
import { JSDOM } from "jsdom";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkParse from "remark-parse";
import { unified } from "unified";

const require = createRequire(import.meta.url);
const mathjaxLoader = require("mathjax");
const root = process.cwd();
const parser = unified().use(remarkParse).use(remarkGfm).use(remarkMath);
const failures = [];

function location(file, node) {
  const line = node?.position?.start?.line;
  return line ? `${file}:${line}` : file;
}

function report(file, node, kind, error) {
  const message = error instanceof Error ? error.message : String(error);
  failures.push(`${location(file, node)} [${kind}] ${message.replaceAll("\n", " ")}`);
}

function walk(node, visit) {
  visit(node);
  if (!Array.isArray(node.children)) return;
  for (const child of node.children) walk(child, visit);
}

function plainText(node) {
  if (typeof node.value === "string") return node.value;
  if (!Array.isArray(node.children)) return "";
  return node.children.map(plainText).join("");
}

function splitTarget(url) {
  const hashIndex = url.indexOf("#");
  return hashIndex < 0
    ? { pathname: url, fragment: "" }
    : { pathname: url.slice(0, hashIndex), fragment: url.slice(hashIndex + 1) };
}

function isExternal(url) {
  return /^(?:[a-z][a-z\d+.-]*:|\/\/)/iu.test(url);
}

async function exists(file) {
  try {
    await access(file, constants.F_OK);
    return true;
  } catch {
    return false;
  }
}

function anchorsFor(tree) {
  const slugger = new GithubSlugger();
  const anchors = new Set();
  walk(tree, (node) => {
    if (node.type === "heading") anchors.add(slugger.slug(plainText(node)));
  });
  return anchors;
}

async function loadMarkdown(target, cache) {
  const absolute = path.resolve(target);
  if (cache.has(absolute)) return cache.get(absolute);
  const source = await readFile(absolute, "utf8");
  const tree = parser.parse(source);
  const parsed = { source, tree, anchors: anchorsFor(tree) };
  cache.set(absolute, parsed);
  return parsed;
}

async function validateLink(file, node, cache) {
  const url = node.url?.trim();
  if (!url || isExternal(url)) return;

  const { pathname, fragment } = splitTarget(url);
  let decodedPath;
  let decodedFragment;
  try {
    decodedPath = decodeURIComponent(pathname);
    decodedFragment = decodeURIComponent(fragment);
  } catch {
    report(file, node, "link", `некоректне percent-encoding: ${url}`);
    return;
  }

  const sourcePath = path.join(root, file);
  const targetPath = decodedPath
    ? path.resolve(path.dirname(sourcePath), decodedPath)
    : sourcePath;

  if (!(await exists(targetPath))) {
    report(file, node, "link", `ціль не існує: ${url}`);
    return;
  }

  if (!decodedFragment || path.extname(targetPath).toLowerCase() !== ".md") return;
  const target = await loadMarkdown(targetPath, cache);
  const normalized = new GithubSlugger().slug(decodedFragment);
  if (!target.anchors.has(decodedFragment) && !target.anchors.has(normalized)) {
    report(file, node, "anchor", `розділ не існує: ${url}`);
  }
}

async function main() {
  const files = (await readdir(root))
    .filter((name) => name.toLowerCase().endsWith(".md"))
    .sort();

  if (files.length === 0) {
    throw new Error(`Markdown-файли не знайдено у ${root}`);
  }

  const mathjax = await mathjaxLoader.init({
    loader: { load: ["input/tex", "output/svg"] }
  });
  const dom = new JSDOM("<!doctype html><html><body></body></html>");
  globalThis.window = dom.window;
  globalThis.document = dom.window.document;
  globalThis.Node = dom.window.Node;
  const { default: mermaid } = await import("mermaid");
  const cache = new Map();
  let mermaidCount = 0;
  let mathCount = 0;
  let linkCount = 0;

  for (const file of files) {
    const { tree } = await loadMarkdown(path.join(root, file), cache);
    const nodes = [];
    walk(tree, (node) => nodes.push(node));

    for (const node of nodes) {
      if (node.type === "code" && node.lang?.toLowerCase() === "mermaid") {
        mermaidCount += 1;
        try {
          await mermaid.parse(node.value);
        } catch (error) {
          report(file, node, "mermaid", error);
        }
      }

      const isMathFence = node.type === "code" && node.lang?.toLowerCase() === "math";
      if (node.type === "math" || node.type === "inlineMath" || isMathFence) {
        mathCount += 1;
        try {
          const rendered = await mathjax.tex2svgPromise(node.value, {
            display: node.type !== "inlineMath"
          });
          const html = mathjax.startup.adaptor.outerHTML(rendered);
          const match = html.match(/data-mjx-error="([^"]+)"/u);
          if (match) throw new Error(match[1]);
        } catch (error) {
          report(file, node, "math", error);
        }
      }

      if (node.type === "link" || node.type === "image") {
        linkCount += 1;
        await validateLink(file, node, cache);
      }
    }
  }

  if (failures.length > 0) {
    console.error(`GitHub Markdown validation: ${failures.length} error(s)`);
    for (const failure of failures) console.error(`- ${failure}`);
    process.exitCode = 1;
    return;
  }

  console.log(
    `GitHub Markdown validation passed: ${files.length} files, ` +
      `${mermaidCount} Mermaid diagrams, ${mathCount} formulas, ${linkCount} links/images.`
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
