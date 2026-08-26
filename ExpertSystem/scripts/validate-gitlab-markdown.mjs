#!/usr/bin/env node

import { access, readdir, readFile } from "node:fs/promises";
import { constants } from "node:fs";
import path from "node:path";
import process from "node:process";
import { JSDOM } from "jsdom";
import katex from "katex";
import remarkGfm from "remark-gfm";
import remarkMath from "remark-math";
import remarkParse from "remark-parse";
import { unified } from "unified";

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

function gitlabAnchorBase(value) {
  return value
    .toLocaleLowerCase()
    .replace(/[^\p{L}\p{N}\s_-]/gu, "")
    .replace(/\s/gu, "-");
}

function anchorsFor(tree) {
  const counts = new Map();
  const anchors = new Set();
  walk(tree, (node) => {
    if (node.type !== "heading") return;
    const base = gitlabAnchorBase(plainText(node));
    const duplicate = counts.get(base) ?? 0;
    const anchor = duplicate === 0 ? base : `${base}-${duplicate}`;
    counts.set(base, duplicate + 1);
    anchors.add(anchor);
  });
  return anchors;
}

async function exists(file) {
  try {
    await access(file, constants.F_OK);
    return true;
  } catch {
    return false;
  }
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

function validateHeadingStructure(file, tree) {
  const headings = [];
  walk(tree, (node) => {
    if (node.type === "heading") headings.push(node);
  });

  const levelOne = headings.filter((node) => node.depth === 1);
  if (levelOne.length !== 1) {
    report(file, levelOne[1] ?? headings[0], "headings", `очікується рівно один H1, отримано ${levelOne.length}`);
  }
  if (headings[0]?.depth !== 1) {
    report(file, headings[0], "headings", "перший заголовок має бути H1");
  }

  for (let index = 1; index < headings.length; index += 1) {
    const previous = headings[index - 1];
    const current = headings[index];
    if (current.depth > previous.depth + 1) {
      report(file, current, "headings", `пропущено рівень після H${previous.depth}: H${current.depth}`);
    }
  }
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
  if (!target.anchors.has(decodedFragment)) {
    report(file, node, "anchor", `GitLab anchor не існує: ${url}`);
  }
}

async function main() {
  const files = (await readdir(root))
    .filter((name) => name.toLowerCase().endsWith(".md"))
    .sort();

  if (files.length === 0) {
    throw new Error(`Markdown-файли не знайдено у ${root}`);
  }

  const dom = new JSDOM("<!doctype html><html><body></body></html>");
  globalThis.window = dom.window;
  globalThis.document = dom.window.document;
  globalThis.Node = dom.window.Node;
  const { default: mermaid } = await import("mermaid");
  const cache = new Map();
  let mermaidCount = 0;
  let accessibleMermaidCount = 0;
  let mathCount = 0;
  let linkCount = 0;

  for (const file of files) {
    const { tree } = await loadMarkdown(path.join(root, file), cache);
    validateHeadingStructure(file, tree);
    const nodes = [];
    walk(tree, (node) => nodes.push(node));

    for (const node of nodes) {
      if (node.type === "code" && node.lang?.toLowerCase() === "mermaid") {
        mermaidCount += 1;
        const hasTitle = /^\s*accTitle\s*:/imu.test(node.value);
        const hasDescription = /^\s*accDescr\s*:/imu.test(node.value);
        if (hasTitle !== hasDescription) {
          report(file, node, "mermaid-accessibility", "accTitle і accDescr мають використовуватися разом");
        }
        if (hasTitle && hasDescription) accessibleMermaidCount += 1;
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
          katex.renderToString(node.value, {
            displayMode: node.type !== "inlineMath",
            output: "htmlAndMathml",
            strict: "error",
            throwOnError: true
          });
        } catch (error) {
          report(file, node, "katex", error);
        }
      }

      if (node.type === "link" || node.type === "image") {
        linkCount += 1;
        await validateLink(file, node, cache);
      }
    }
  }

  if (failures.length > 0) {
    console.error(`GitLab Markdown validation: ${failures.length} error(s)`);
    for (const failure of failures) console.error(`- ${failure}`);
    process.exitCode = 1;
    return;
  }

  console.log(
    `GitLab Markdown validation passed: ${files.length} files, ` +
      `${mermaidCount} Mermaid diagrams (${accessibleMermaidCount} with accessibility metadata), ` +
      `${mathCount} KaTeX formulas, ${linkCount} links/images.`
  );
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
