// Audit tool: check that every extracted claim quotes its source verbatim at the stated lines; lists sources still "в роботі".
// Usage: node verify-quotes.mjs <audit-requirements-dir> [--repo .] [--slack 3]
// Expects the per-source format of ARTIFACTS.md §3.1-3.2: header "# Джерело: `path`", sections "## RQ-NNN",
// fields "- **Джерело:** ... рядки N–M" and "- **Формулювання:** «...»".
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const dir = args.find((a, i) => !a.startsWith('--') && !['--repo', '--slack'].includes(args[i - 1]));
if (!dir) {
  console.error('usage: node verify-quotes.mjs <audit-requirements-dir> [--repo .] [--slack 3]');
  process.exit(2);
}
const repo = path.resolve(opt('--repo', '.'));
const slack = Number(opt('--slack', '3'));
const norm = (s) => s.replace(/\s+/g, ' ').trim();

function walk(p, acc) {
  if (fs.statSync(p).isDirectory()) for (const n of fs.readdirSync(p)) walk(path.join(p, n), acc);
  else if (p.endsWith('.md')) acc.push(p);
  return acc;
}

// Normalized text of the source with the line number of every character offset.
function indexSource(file) {
  const lines = fs.readFileSync(file, 'utf8').split(/\r?\n/).map(norm);
  const starts = [];
  let text = '';
  lines.forEach((l, i) => { starts.push(text.length); text += (i ? ' ' : '') + l; });
  const lineAt = (offset) => { let i = starts.length - 1; while (i > 0 && starts[i] > offset) i--; return i + 1; };
  return { lines, text, starts, lineAt };
}

function findFragments(src, quote, from, to) {
  const frags = quote.split(/\s*(?:…|\.\.\.)\s*/).map(norm).filter(Boolean);
  const lo = from ? src.starts[Math.max(0, from - 1 - slack)] : 0;
  const hiLine = to ? Math.min(src.lines.length, to + slack) : src.lines.length;
  const hi = hiLine < src.lines.length ? src.starts[hiLine] : src.text.length;
  let pos = lo;
  let first = null;
  for (const f of frags) {
    const at = src.text.indexOf(f, pos);
    if (at < 0 || at + f.length > hi) return { ok: false, first };
    first ??= at;
    pos = at + f.length;
  }
  return { ok: frags.length > 0, first };
}

const files = walk(path.resolve(dir), []);
const problems = [];
const unfinished = [];
let rqTotal = 0, ok = 0, shifted = 0;
for (const file of files) {
  const lines = fs.readFileSync(file, 'utf8').split(/\r?\n/);
  const head = /^# Джерело: `([^`]+)`/.exec(lines[0] ?? '');
  if (!head) continue;
  const rel = path.relative(process.cwd(), file).replace(/\\/g, '/');
  const state = lines.find((l) => /^- \*\*Стан:\*\*/.test(l));
  if (state && /в роботі/.test(state)) unfinished.push(rel);
  const srcPath = path.join(repo, head[1]);
  if (!fs.existsSync(srcPath)) { problems.push(`${rel}: джерело ${head[1]} не існує`); continue; }
  const src = indexSource(srcPath);
  let id = null, from = null, to = null;
  for (let i = 0; i < lines.length; i++) {
    const sec = /^## (RQ-\d+)/.exec(lines[i]);
    if (sec) { id = sec[1]; from = to = null; continue; }
    const loc = /^- \*\*Джерело:\*\*.*?рядки?\s+(\d+)(?:\s*[–-]\s*(\d+))?/.exec(lines[i]);
    if (loc) { from = Number(loc[1]); to = Number(loc[2] ?? loc[1]); continue; }
    if (!id || !/^- \*\*Формулювання:\*\*/.test(lines[i])) continue;
    let raw = lines[i].replace(/^- \*\*Формулювання:\*\*\s*/, '');
    while (!/»\s*$/.test(raw) && i + 1 < lines.length && !/^- \*\*|^## /.test(lines[i + 1])) raw += ` ${lines[++i]}`;
    const quote = raw.replace(/^«/, '').replace(/»\s*$/, '');
    rqTotal++;
    if (findFragments(src, quote, from, to).ok) { ok++; continue; }
    const anywhere = findFragments(src, quote, null, null);
    if (anywhere.ok) {
      shifted++;
      problems.push(`${rel} ${id}: цитата на рядку ${src.lineAt(anywhere.first)}, а не ${from}–${to}`);
    } else {
      problems.push(`${rel} ${id}: цитату не знайдено в ${head[1]}`);
    }
  }
}

console.log(`джерел=${files.length} незавершених=${unfinished.length} RQ=${rqTotal} збіг=${ok} зсув=${shifted} не_знайдено=${rqTotal - ok - shifted}`);
for (const u of unfinished) console.log(`в роботі: ${u}`);
for (const p of problems) console.log(p);
process.exitCode = rqTotal - ok - shifted > 0 ? 1 : 0;
