// Audit tool: files and lines per directory for scoping and estimates (blank lines included).
// Usage: node inventory.mjs [path ...] [--ext .md,.yaml] [--exclude dir/sub,...] [--per-file]
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const valueFlags = new Set(['--ext', '--exclude']);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const exts = opt('--ext', '.md').split(',').filter(Boolean);
const perFile = args.includes('--per-file');
const roots = args.filter((a, i) => !a.startsWith('--') && !valueFlags.has(args[i - 1]));
if (roots.length === 0) roots.push('.');
const defaultExcludes = roots.some((r) => /^audit([\\/]|$)/.test(r)) ? ['node_modules', '.git'] : ['node_modules', '.git', 'audit'];
const excludes = [...defaultExcludes, ...opt('--exclude', '').split(',').filter(Boolean)]
  .map((e) => e.replace(/\\/g, '/').replace(/^\/+|\/+$/g, ''));

const cwd = process.cwd();
const rel = (p) => path.relative(cwd, p).replace(/\\/g, '/') || '.';
const excluded = (p) => excludes.some((e) => `/${rel(p)}/`.includes(`/${e}/`));

function collect(p, acc) {
  if (excluded(p)) return acc;
  const st = fs.statSync(p);
  if (st.isDirectory()) {
    for (const name of fs.readdirSync(p)) collect(path.join(p, name), acc);
  } else if (exts.some((x) => p.endsWith(x))) {
    acc.push(p);
  }
  return acc;
}

function count(file) {
  const text = fs.readFileSync(file, 'utf8');
  if (text === '') return { lines: 0, nonBlank: 0 };
  const parts = text.split(/\r?\n/);
  if (text.endsWith('\n')) parts.pop();
  return { lines: parts.length, nonBlank: parts.filter((l) => l.trim() !== '').length };
}

const files = roots.flatMap((r) => collect(path.resolve(r), [])).sort();
const byDir = new Map();
const rows = [];
let total = { files: 0, lines: 0, nonBlank: 0 };
for (const f of files) {
  const c = count(f);
  rows.push([rel(f), c.lines, c.nonBlank]);
  const d = rel(path.dirname(f));
  const agg = byDir.get(d) ?? { files: 0, lines: 0, nonBlank: 0 };
  agg.files++; agg.lines += c.lines; agg.nonBlank += c.nonBlank;
  byDir.set(d, agg);
  total.files++; total.lines += c.lines; total.nonBlank += c.nonBlank;
}

if (perFile) {
  console.log('| Файл | Рядків | Непорожніх |\n| --- | ---: | ---: |');
  for (const [f, l, n] of rows) console.log(`| \`${f}\` | ${l} | ${n} |`);
} else {
  console.log('| Каталог | Файлів | Рядків | Непорожніх |\n| --- | ---: | ---: | ---: |');
  for (const [d, a] of [...byDir].sort()) console.log(`| \`${d}\` | ${a.files} | ${a.lines} | ${a.nonBlank} |`);
}
console.log(`| Разом | ${perFile ? '' : `${total.files} | `}${total.lines} | ${total.nonBlank} |`);
