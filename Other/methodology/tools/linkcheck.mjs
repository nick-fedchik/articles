// Audit tool: relative Markdown links and anchors; --tracked also flags targets ignored by Git (they break in a clean clone).
// Usage: node linkcheck.mjs [root=.] [--tracked] [--exclude dir/sub,...] [--json]
import fs from 'node:fs';
import path from 'node:path';
import { execSync } from 'node:child_process';

const args = process.argv.slice(2);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const root = path.resolve(args.find((a, i) => !a.startsWith('--') && args[i - 1] !== '--exclude') ?? '.');
const trackedOnly = args.includes('--tracked');
const userExcludes = opt('--exclude', '').split(',').filter(Boolean).map((e) => e.replace(/\\/g, '/').replace(/^\/+|\/+$/g, ''));
const skipSegment = (s) => ['node_modules', '.git', 'dist', 'out', 'third_party', 'audit'].includes(s) || /^build/.test(s);
const relOf = (p) => path.relative(root, p).replace(/\\/g, '/');
const excluded = (p) => relOf(p).split('/').some(skipSegment) || userExcludes.some((e) => `/${relOf(p)}/`.includes(`/${e}/`));

function walk(dir, acc) {
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name);
    if (excluded(p)) continue;
    if (e.isDirectory()) walk(p, acc);
    else if (e.name.endsWith('.md')) acc.push(p);
  }
  return acc;
}

const git = (cmd) => execSync(cmd, { cwd: root, encoding: 'utf8' }).split(/\r?\n/).filter(Boolean);
const files = trackedOnly
  ? git('git ls-files "*.md"').map((f) => path.join(root, f)).filter((f) => !excluded(f))
  : walk(root, []);
const tracked = trackedOnly ? new Set(git('git ls-files').map((f) => path.normalize(path.join(root, f)).toLowerCase())) : null;

const slug = (h) => h.trim().toLowerCase().replace(/<[^>]+>/g, '').replace(/[`*~]/g, '')
  .replace(/[^\p{L}\p{N}\s_-]/gu, '').replace(/\s/g, '-');
const anchorCache = new Map();
function anchors(file) {
  if (anchorCache.has(file)) return anchorCache.get(file);
  const set = new Set();
  const seen = new Map();
  let fence = false;
  for (const line of fs.readFileSync(file, 'utf8').split(/\r?\n/)) {
    if (/^\s*(```|~~~)/.test(line)) { fence = !fence; continue; }
    if (fence) continue;
    const m = /^#{1,6}\s+(.*?)\s*#*\s*$/.exec(line);
    if (!m) continue;
    const s = slug(m[1]);
    const n = seen.get(s) ?? 0;
    seen.set(s, n + 1);
    set.add(n ? `${s}-${n}` : s);
  }
  anchorCache.set(file, set);
  return set;
}

const broken = [];
let total = 0;
for (const file of files) {
  let fence = false;
  fs.readFileSync(file, 'utf8').split(/\r?\n/).forEach((line, i) => {
    if (/^\s*(```|~~~)/.test(line)) { fence = !fence; return; }
    if (fence) return;
    const re = /\[[^\]]*\]\(([^)\s]+)(?:\s+"[^"]*")?\)/g;
    const text = line.replace(/`[^`]*`/g, '');
    let m;
    while ((m = re.exec(text))) {
      const target = m[1];
      if (/^[a-z][a-z0-9+.-]*:/i.test(target)) continue;
      total++;
      const [rawPath, frag] = target.split('#');
      const p = decodeURIComponent(rawPath);
      const abs = p ? path.resolve(path.dirname(file), p) : file;
      const where = { file: relOf(file), line: i + 1, target };
      if (p && !fs.existsSync(abs)) { broken.push({ ...where, reason: 'файл не існує' }); continue; }
      const isFile = fs.statSync(abs).isFile();
      if (tracked && p && isFile && !tracked.has(path.normalize(abs).toLowerCase())) {
        broken.push({ ...where, reason: 'ціль не відстежується Git' });
        continue;
      }
      if (frag && isFile && abs.endsWith('.md') && !/^L\d+/.test(frag) && !anchors(abs).has(decodeURIComponent(frag).toLowerCase())) {
        broken.push({ ...where, reason: 'немає якоря' });
      }
    }
  });
}

if (args.includes('--json')) {
  console.log(JSON.stringify({ files: files.length, links: total, broken }, null, 1));
} else {
  console.log(`files=${files.length} links=${total} broken=${broken.length}`);
  for (const b of broken) console.log(`${b.file}:${b.line}\t${b.target}\t${b.reason}`);
}
process.exitCode = broken.length ? 1 : 0;
