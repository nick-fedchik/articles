// Audit tool (quick scan): names that documentation mentions (commands, settings, endpoints, tables) but code does not contain.
// Usage: node names-check.mjs --pattern "<regex>" --docs <path,...> --code <path,...> [--exclude dir/sub,...] [--code-ext .ts,.cpp,...]
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const pattern = opt('--pattern', null);
const docRoots = opt('--docs', '').split(',').filter(Boolean);
const codeRoots = opt('--code', '').split(',').filter(Boolean);
if (!pattern || docRoots.length === 0 || codeRoots.length === 0) {
  console.error('usage: node names-check.mjs --pattern "<regex>" --docs <path,...> --code <path,...>');
  process.exit(2);
}
const codeExts = opt('--code-ext', '.ts,.js,.mjs,.json,.cpp,.hpp,.h,.c,.go,.py,.ps1,.sql,.vue,.yaml,.yml,.txt,.cmake').split(',');
const excludes = ['node_modules', '.git', 'out', 'dist', 'audit', ...opt('--exclude', '').split(',').filter(Boolean)]
  .map((e) => e.replace(/\\/g, '/').replace(/^\/+|\/+$/g, ''));
const cwd = process.cwd();
const rel = (p) => path.relative(cwd, p).replace(/\\/g, '/');
const excluded = (p) => excludes.some((e) => `/${rel(p)}/`.includes(`/${e}/`)) || rel(p).split('/').some((s) => /^build/.test(s));

function collect(p, exts, acc) {
  if (!fs.existsSync(p) || excluded(p)) return acc;
  if (fs.statSync(p).isDirectory()) for (const n of fs.readdirSync(p)) collect(path.join(p, n), exts, acc);
  else if (exts.some((x) => p.endsWith(x))) acc.push(p);
  return acc;
}

const re = new RegExp(pattern, 'g');
const mentions = new Map();
for (const f of docRoots.flatMap((r) => collect(path.resolve(r), ['.md'], []))) {
  fs.readFileSync(f, 'utf8').split(/\r?\n/).forEach((line, i) => {
    for (const m of line.matchAll(re)) {
      const name = m[0].replace(/[.,;:)]+$/, '');
      if (!mentions.has(name)) mentions.set(name, []);
      mentions.get(name).push(`${rel(f)}:${i + 1}`);
    }
  });
}
const code = codeRoots.flatMap((r) => collect(path.resolve(r), codeExts, [])).map((f) => fs.readFileSync(f, 'utf8')).join('\n');

const missing = [...mentions].filter(([name]) => !code.includes(name)).sort();
console.log(`Шаблон: \`${pattern}\`. Різних назв у документації: ${mentions.size}. Немає в коді: ${missing.length}.\n`);
if (missing.length) console.log('| Назва | Згадок | Де (до 5) |\n| --- | ---: | --- |');
for (const [name, at] of missing) console.log(`| \`${name}\` | ${at.length} | ${at.slice(0, 5).map((a) => `\`${a}\``).join(', ')} |`);
process.exitCode = missing.length ? 1 : 0;
