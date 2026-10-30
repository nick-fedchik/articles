import assert from 'node:assert/strict';
import crypto from 'node:crypto';
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const toolRoot = path.join(root, 'review/tmp/technical-tools');
if (process.argv.includes('--setup')) {
  fs.mkdirSync(toolRoot, { recursive: true });
  for (const name of ['package.json', 'package-lock.json']) {
    fs.copyFileSync(path.join(root, 'scripts/render-validation', name), path.join(toolRoot, name));
  }
  console.log('Dependency manifests prepared under review/tmp/technical-tools.');
  console.log('npm --prefix review/tmp/technical-tools --cache review/tmp/npm-cache ci --engine-strict --ignore-scripts --no-fund --no-audit');
  process.exit(0);
}
const local = createRequire(path.join(toolRoot, 'package.json'));
const load = name => import(local.resolve(name));
const { JSDOM } = local('jsdom');

const { unified } = await load('unified');
const { default: remarkParse } = await load('remark-parse');
const { default: remarkGfm } = await load('remark-gfm');
const { default: remarkMath } = await load('remark-math');
const { default: GithubSlugger } = await load('github-slugger');
const mathjax = local('mathjax');
const mathjaxPackageConfig = {
  removed: ['noundefined', 'autoload', 'require'],
  added: ['boldsymbol', 'color'],
};
const mermaidSecurity = 'strict';
const observedGithubMacroRejections = [{
  macro: 'operatorname',
  observedAt: '2026-10-30',
  evidence: 'Author-reported GitHub warning on ExpertSystem/ch29-neuro-symbolic-architecture.md',
}];
await mathjax.init({
  loader: { load: ['input/tex', 'output/svg', '[tex]/boldsymbol', '[tex]/color'] },
  tex: { packages: { '[-]': mathjaxPackageConfig.removed, '[+]': mathjaxPackageConfig.added } },
});
const dom = new JSDOM('<!doctype html><html><body></body></html>');
globalThis.window = dom.window;
globalThis.document = dom.window.document;
globalThis.Node = dom.window.Node;
globalThis.HTMLElement = dom.window.HTMLElement;
globalThis.SVGElement = dom.window.SVGElement;
const { default: mermaid } = await load('mermaid');
mermaid.initialize({ startOnLoad: false, securityLevel: mermaidSecurity });
const parser = unified().use(remarkParse).use(remarkGfm).use(remarkMath);
const sha = text => crypto.createHash('sha256').update(text).digest('hex');

function observedGithubMacroViolations(object) {
  if (object.type !== 'math') return [];
  return observedGithubMacroRejections
    .filter(({ macro }) => new RegExp(`\\\\${macro}(?![A-Za-z])`).test(object.source))
    .map(({ macro, observedAt, evidence }) => ({ file: object.file, line: object.line, macro, observedAt, evidence }));
}

function walk(node, visit) {
  visit(node);
  for (const child of node.children ?? []) walk(child, visit);
}

function plain(node) {
  if (node.type === 'html') return '';
  return node.value ?? (node.children ?? []).map(plain).join('');
}

function extract(text, file = '<control>') {
  const prepared = text.replace(/\$`[^`]*`\$/g, match => match.replace(/\\\$/g, 'ZZ'));
  const tree = parser.parse(prepared);
  const objects = [];
  const links = [];
  const listings = [];
  const codeBlocks = [];
  const chapterSections = [];
  const anchors = new Set();
  const slugger = new GithubSlugger();
  walk(tree, node => {
    if (node.type === 'heading') {
      const title = plain(node);
      anchors.add(slugger.slug(title));
      const match = node.depth === 2 && title.match(/^(\d+)\. Діалог до глави/);
      if (match) chapterSections.push({ number: Number(match[1]), offset: node.position.start.offset });
    }
    if (node.type === 'html') {
      const fragment = JSDOM.fragment(node.value);
      for (const element of fragment.querySelectorAll('[id], a[name]')) {
        if (element.id) anchors.add(element.id);
        if (element.getAttribute('name')) anchors.add(element.getAttribute('name'));
      }
    }
    if (['link', 'image', 'definition'].includes(node.type)) {
      links.push({ file, line: node.position.start.line, url: node.url,
        start: node.position.start.offset, end: node.position.end.offset,
        raw: text.slice(node.position.start.offset, node.position.end.offset) });
    }
    if (node.type === 'code' && ['go', 'golang'].includes(node.lang)) {
      listings.push({ file, line: node.position.start.line, source_sha256: sha(text), source: node.value,
        skip: /skip-validate/.test(node.value) || /<!--\s*skip-validate\s*-->\s*$/.test(text.slice(Math.max(0, node.position.start.offset - 100), node.position.start.offset)) });
    }
    if (node.type === 'code') codeBlocks.push({ file, line: node.position.start.line, language: node.lang ?? 'unlabelled',
      source_sha256: sha(text), source: node.value });
    const isMath = ['math', 'inlineMath'].includes(node.type) || (node.type === 'code' && node.lang === 'math');
    const isMermaid = node.type === 'code' && node.lang === 'mermaid';
    if (!isMath && !isMermaid) return;
    const raw = text.slice(node.position.start.offset, node.position.end.offset);
    let value = node.value;
    if (node.type === 'inlineMath' && raw.startsWith('$`') && raw.endsWith('`$')) value = raw.slice(2, -2);
    objects.push({
      file, line: node.position.start.line, column: node.position.start.column,
      source_sha256: sha(text), type: isMath ? 'math' : 'mermaid',
      display: node.type !== 'inlineMath', source: value,
      raw,
    });
  });
  return { objects, links, anchors, listings, chapterSections, codeBlocks };
}

async function validate(object) {
  try {
    if (object.type === 'math') {
      mathjax.texReset();
      const result = await mathjax.tex2svgPromise(object.source, { display: object.display });
      const output = mathjax.startup.adaptor.serializeXML(result);
      if (/data-mjx-error|<merror\b/i.test(output)) throw new Error('MathJax error node in output');
    } else {
      const result = await mermaid.parse(object.source);
      if (!result) throw new Error('Mermaid parser returned a negative result');
    }
    return { ...object, status: 'PASS' };
  } catch (error) {
    return { ...object, status: 'FAIL', error: String(error.message ?? error) };
  }
}

async function controls() {
  const positiveCases = [
    { type: 'math', source: '\\frac{1}{x^2+1}', display: true },
    { type: 'math', source: '\\mathrm{Canon}(x)', display: false },
    { type: 'math', source: '\\hat{w}_{1:N}=\\arg\\max_{w_{1:N}}P(w_{1:N})', display: true },
    { type: 'math', source: '\\operatorname{Canon}(x)', display: false },
    { type: 'mermaid', source: 'flowchart LR\nStart["Ready"] --> End["Done"]' },
  ];
  const negativeCases = [
    { type: 'math', source: '\\frac{1}{', display: true },
    { type: 'math', source: '\\BookUnknownCommand{x}', display: false },
    { type: 'mermaid', source: 'flowchart LR\nStart["Unclosed' },
  ];
  const engineResults = [];
  for (const object of positiveCases) {
    const actual = (await validate(object)).status;
    assert.equal(actual, 'PASS', JSON.stringify(object));
    engineResults.push({ expected: 'PASS', actual, type: object.type });
  }
  for (const object of negativeCases) {
    const actual = (await validate(object)).status;
    assert.equal(actual, 'FAIL', JSON.stringify(object));
    engineResults.push({ expected: 'FAIL', actual, type: object.type });
  }
  const sample = '# Heading\n\n$`x_1`$ and $y^2$.\n\n```go\nvalue := "$not_math$"\n```\n\n```math\nz=3\n```\n\n```mermaid\nflowchart LR\nA-->B\n```\n';
  const extracted = extract(sample);
  const extractionChecks = [
    extracted.objects.filter(item => item.type === 'math').length === 3,
    extracted.objects.filter(item => item.type === 'mermaid').length === 1,
    extracted.objects[0].source === 'x_1',
    extracted.objects[0].line === 3,
    extracted.anchors.has('heading'),
    extracted.listings.length === 1,
    extracted.codeBlocks.length === 3,
  ];
  const currency = extract('$`\\text{\\$4}`$');
  extractionChecks.push(currency.objects.length === 1);
  extractionChecks.push(currency.objects[0]?.source === '\\text{\\$4}');
  extractionChecks.push((await validate(currency.objects[0])).status === 'PASS');
  assert(extractionChecks.every(Boolean));
  const githubPositive = observedGithubMacroViolations({ type: 'math', source: '\\mathrm{Canon}(x)' });
  const githubNegative = observedGithubMacroViolations({ type: 'math', source: '\\operatorname{Canon}(x)' });
  assert.equal(githubPositive.length, 0);
  assert.equal(githubNegative[0]?.macro, 'operatorname');
  console.log('PASS: MathJax/Mermaid controls, observed GitHub-macro control, and Markdown extraction controls');
  return {
    engine: { passed: engineResults.filter(item => item.expected === item.actual).length, total: engineResults.length, cases: engineResults },
    extraction: { passed: extractionChecks.filter(Boolean).length, total: extractionChecks.length },
    observed_github_macro: { positive_control: 'PASS', reported_rejection_detected: 'PASS' },
  };
}

const controlResults = await controls();
if (process.argv.includes('--self-test')) process.exit(0);

const book = path.join(root, 'ExpertSystem');
const files = fs.readdirSync(book).filter(name => name.endsWith('.md') && name !== 'AGENTS.md').map(name => path.join(book, name));
for (const directory of ['registry', 'requirements/L1-stakeholder', 'requirements/L2-system', 'requirements/L3-component']) {
  const folder = path.join(book, directory);
  files.push(...fs.readdirSync(folder).filter(name => name.endsWith('.md')).map(name => path.join(folder, name)));
}
files.sort();
assert.equal(files.filter(file => path.dirname(file) === book && /^ch\d\d-/.test(path.basename(file))).length, 40);
const inventory = new Map();
for (const file of files) {
  const text = fs.readFileSync(file, 'utf8');
  inventory.set(file, { text, ...extract(text, path.relative(root, file)) });
}
if (process.argv.includes('--fix-links')) {
  const aliases = {
    'ch16-symbolic-reasoning-and-rule-engines.md': 'ch16-expert-systems-architecture.md',
    'ch35-hardware-accelerators-fpga-asic.md': 'ch35-self-organizing-expert-systems-synergetics-and-npu-runtime.md',
    'ch36-testing-and-Popperian-falsification.md': 'ch36-knowledge-testing-pyramid-and-variational-calibration.md',
  };
  const revisions = [];
  for (const [file, record] of inventory) {
    const patches = [];
    for (const link of record.links) {
      let replacement = '';
      if (aliases[link.url]) replacement = aliases[link.url];
      else if (/^#src-\d+\]$/.test(link.url)) replacement = link.url.slice(0, -1);
      else if (/^#src-\d+$/.test(link.url)) {
        let number = null;
        if (path.basename(file) === 'PRO_ET_CONTRA.md') {
          number = record.chapterSections.filter(section => section.offset < link.start).at(-1)?.number;
        } else if (path.dirname(file) === path.join(book, 'registry')) {
          number = Number(path.basename(file).match(/^ch(\d\d)-/)?.[1]);
        }
        if (number) {
          const chapter = files.find(candidate => path.dirname(candidate) === book && path.basename(candidate).startsWith(`ch${String(number).padStart(2, '0')}-`));
          if (chapter) replacement = path.relative(path.dirname(file), chapter) + link.url;
        }
      }
      if (!replacement || replacement === link.url) continue;
      const [base, fragment] = replacement.split('#');
      const destination = base ? path.resolve(path.dirname(file), base) : file;
      if (!fs.existsSync(destination) || (fragment && !inventory.get(destination)?.anchors.has(fragment))) continue;
      const position = link.raw.lastIndexOf('](') + 2;
      if (position < 2 || !link.raw.slice(position).startsWith(link.url)) continue;
      patches.push({ start: link.start + position, end: link.start + position + link.url.length, replacement,
        line: link.line, before: link.url });
    }
    if (!patches.length) continue;
    let updated = record.text;
    for (const patch of patches.sort((left, right) => right.start - left.start)) updated = updated.slice(0, patch.start) + patch.replacement + updated.slice(patch.end);
    if (fs.readFileSync(file, 'utf8') !== record.text) throw new Error(`Concurrent modification: ${file}`);
    revisions.push({ file: path.relative(root, file), before_sha256: sha(record.text), after_sha256: sha(updated), patches, original: record.text });
    fs.writeFileSync(file, updated);
  }
  const output = path.join(root, 'review/tmp/technical-preparation/link-fixes.json');
  fs.mkdirSync(path.dirname(output), { recursive: true });
  fs.writeFileSync(output, JSON.stringify(revisions, null, 2) + '\n');
  console.log(JSON.stringify({ files_changed: revisions.length, links_fixed: revisions.reduce((sum, item) => sum + item.patches.length, 0), report: path.relative(root, output) }));
  process.exit(0);
}
const objectResults = [];
for (const record of inventory.values()) {
  for (const object of record.objects) objectResults.push(await validate(object));
}
const linkResults = [];
for (const [file, record] of inventory) {
  for (const link of record.links) {
    if (/^[a-z][a-z\d+.-]*:/i.test(link.url) || link.url.startsWith('//')) continue;
    const separator = link.url.indexOf('#');
    const base = separator < 0 ? link.url : link.url.slice(0, separator);
    const fragment = separator < 0 ? '' : decodeURIComponent(link.url.slice(separator + 1));
    const destination = base ? path.resolve(path.dirname(file), decodeURIComponent(base)) : file;
    if (destination.includes(`${path.sep}translations${path.sep}`)) {
      linkResults.push({ ...link, status: 'SKIP', reason: 'Translations frozen and outside scope' });
      continue;
    }
    let error = '';
    if (!fs.existsSync(destination)) error = 'Missing local target';
    else if (fragment && fs.statSync(destination).isFile() && destination.endsWith('.md')) {
      let linked = inventory.get(destination);
      if (!linked) {
        linked = extract(fs.readFileSync(destination, 'utf8'), path.relative(root, destination));
        inventory.set(destination, linked);
      }
      if (!linked.anchors.has(fragment)) error = 'Missing local anchor';
    }
    linkResults.push({ ...link, status: error ? 'FAIL' : 'PASS', ...(error ? { error } : {}) });
  }
}
const summary = {};
for (const type of ['math', 'mermaid']) {
  const selected = objectResults.filter(item => item.type === type);
  summary[type] = { total: selected.length, passed: selected.filter(item => item.status === 'PASS').length, failed: selected.filter(item => item.status === 'FAIL').length };
  assert(selected.length > 0);
}
summary.links = { total: linkResults.length, passed: linkResults.filter(item => item.status === 'PASS').length,
  failed: linkResults.filter(item => item.status === 'FAIL').length, skipped: linkResults.filter(item => item.status === 'SKIP').length };
const observedMacroFailures = objectResults.flatMap(object => observedGithubMacroViolations(object)
  .map(violation => ({ ...violation, status: 'FAIL', type: 'github-observed-macro',
    error: `Author-reported GitHub rejection: \\${violation.macro}` })));
summary.github_observed_macro_restrictions = {
  checked_math_expressions: objectResults.filter(item => item.type === 'math').length,
  failed: observedMacroFailures.length,
  observations: observedGithubMacroRejections,
  complete_github_policy: false,
};
const output = path.join(root, 'review/tmp/technical-preparation/render-syntax.json');
fs.mkdirSync(path.dirname(output), { recursive: true });
fs.writeFileSync(output, JSON.stringify({
  at: new Date().toISOString(), node: process.version, dependencies: local('./package.json').dependencies,
  mathjax_package_config: mathjaxPackageConfig, mermaid_security: mermaidSecurity,
  observed_github_macro_restrictions: { complete_policy: false, observations: observedGithubMacroRejections },
  scope: files.map(file => path.relative(root, file)),
  input_hashes: Object.fromEntries(files.map(file => [path.relative(root, file), sha(fs.readFileSync(file, 'utf8'))])),
  controls: controlResults, visual_rendering: 'DEFERRED_S7', summary, objects: objectResults, links: linkResults,
  listings: [...inventory.values()].flatMap(record => record.listings),
  code_blocks: [...inventory.values()].flatMap(record => record.codeBlocks),
}, null, 2) + '\n');
console.log(JSON.stringify({ files: files.length, summary, report: path.relative(root, output) }, null, 2));
for (const failed of [...objectResults, ...linkResults, ...observedMacroFailures].filter(item => item.status === 'FAIL').slice(0, 12)) {
  console.log(`${failed.file}:${failed.line}: ${failed.type ?? 'link'}: ${failed.error}`);
}
process.exitCode = Object.values(summary).some(item => item.failed) ? 1 : 0;