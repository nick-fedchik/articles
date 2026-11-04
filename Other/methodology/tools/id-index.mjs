// Audit tool: ID definitions (by heading, not file name), broken references, numbering gaps, unreferenced IDs.
// Usage: node id-index.mjs --prefixes SHR,SWR,ADR [path ...] [--exclude dir/sub,...] [--out file.md]
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const valueFlags = new Set(['--prefixes', '--exclude', '--out']);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const prefixes = opt('--prefixes', '').split(',').filter(Boolean);
if (prefixes.length === 0) {
  console.error('--prefixes is required, e.g. --prefixes SHR,SWR,ADR');
  process.exit(2);
}
const roots = args.filter((a, i) => !a.startsWith('--') && !valueFlags.has(args[i - 1]));
if (roots.length === 0) roots.push('.');
const defaultExcludes = roots.some((r) => /^audit([\\/]|$)/.test(r)) ? ['node_modules', '.git'] : ['node_modules', '.git', 'audit'];
const excludes = [...defaultExcludes, ...opt('--exclude', '').split(',').filter(Boolean)]
  .map((e) => e.replace(/\\/g, '/').replace(/^\/+|\/+$/g, ''));

const cwd = process.cwd();
const rel = (p) => path.relative(cwd, p).replace(/\\/g, '/');
const excluded = (p) => excludes.some((e) => `/${rel(p)}/`.includes(`/${e}/`));
function collect(p, acc) {
  if (excluded(p)) return acc;
  if (fs.statSync(p).isDirectory()) for (const n of fs.readdirSync(p)) collect(path.join(p, n), acc);
  else if (p.endsWith('.md')) acc.push(p);
  return acc;
}

const alt = prefixes.map((p) => p.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')).join('|');
const refRe = new RegExp(`\\b(${alt})-(\\d+)\\b`, 'g');
const defRe = new RegExp(`^#{1,6}\\s+(?:\\d+(?:\\.\\d+)*\\.?\\s+)?(${alt})-(\\d+)\\b`);
const key = (p, n) => `${p}-${Number(n)}`;

const defs = new Map();
const refs = new Map();
const nameMismatch = [];
const files = roots.flatMap((r) => collect(path.resolve(r), []));
for (const file of files) {
  const lines = fs.readFileSync(file, 'utf8').split(/\r?\n/);
  let fence = false;
  let firstDef = null;
  lines.forEach((line, i) => {
    if (/^\s*(```|~~~)/.test(line)) { fence = !fence; return; }
    if (fence) return;
    const loc = `${rel(file)}:${i + 1}`;
    const d = defRe.exec(line);
    if (d) {
      const k = key(d[1], d[2]);
      if (!defs.has(k)) defs.set(k, { shown: `${d[1]}-${d[2]}`, at: [] });
      defs.get(k).at.push(loc);
      firstDef ??= k;
    }
    for (const m of line.matchAll(refRe)) {
      const k = key(m[1], m[2]);
      if (d && k === key(d[1], d[2])) continue;
      if (!refs.has(k)) refs.set(k, []);
      refs.get(k).push({ loc, file: rel(file) });
    }
  });
  const fn = new RegExp(`^(${alt})-(\\d+)\\b`).exec(path.basename(file));
  if (fn && firstDef && key(fn[1], fn[2]) !== firstDef) nameMismatch.push(`${rel(file)}: файл ${fn[1]}-${fn[2]}, заголовок ${defs.get(firstDef).shown}`);
}

const out = [];
out.push(`# Індекс ID (${prefixes.join(', ')})`, '', `Файлів: ${files.length}. Визначень: ${defs.size}. Згаданих ID: ${refs.size}.`, '');
out.push('## Визначення за префіксами', '', '| Префікс | Визначено | Діапазон | Пропуски в нумерації | Дублікати |', '| --- | ---: | --- | --- | --- |');
for (const p of prefixes) {
  const nums = [...defs.keys()].filter((k) => k.startsWith(`${p}-`)).map((k) => Number(k.slice(p.length + 1))).sort((a, b) => a - b);
  const gaps = [];
  for (let n = nums[0]; nums.length && n <= nums[nums.length - 1]; n++) if (!nums.includes(n)) gaps.push(n);
  const dups = [...defs].filter(([k, v]) => k.startsWith(`${p}-`) && v.at.length > 1).map(([, v]) => `${v.shown} (${v.at.join(', ')})`);
  out.push(`| ${p} | ${nums.length} | ${nums.length ? `${nums[0]}..${nums[nums.length - 1]}` : '-'} | ${gaps.join(', ') || 'немає'} | ${dups.join('; ') || 'немає'} |`);
}
out.push('', '## Биті посилання (ID згадано, але не визначено заголовком)', '');
const broken = [...refs].filter(([k]) => !defs.has(k)).sort();
if (broken.length === 0) out.push('Немає.');
for (const [k, list] of broken) out.push(`- ${k}: ${list.length} згадок; ${list.slice(0, 5).map((r) => `\`${r.loc}\``).join(', ')}${list.length > 5 ? ', …' : ''}`);
out.push('', '## Визначені ID без згадок поза власним файлом', '');
const lonely = [...defs].filter(([k, v]) => !(refs.get(k) ?? []).some((r) => !v.at.some((a) => a.startsWith(`${r.file}:`)))).map(([, v]) => `${v.shown} (\`${v.at[0]}\`)`);
out.push(lonely.length ? lonely.map((s) => `- ${s}`).join('\n') : 'Немає.');
out.push('', '## Назва файлу і заголовок визначають різні ID', '', nameMismatch.length ? nameMismatch.map((s) => `- ${s}`).join('\n') : 'Немає.');

const text = `${out.join('\n')}\n`;
const target = opt('--out', null);
if (target) fs.writeFileSync(target, text);
else process.stdout.write(text);
process.exitCode = broken.length ? 1 : 0;
