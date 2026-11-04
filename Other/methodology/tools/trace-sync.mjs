// Audit tool: check (and with --write synchronize) bidirectional links between ADR "Satisfies:" headers and SHR/SWR "Implemented by ADR-..." lines, plus SHR <-> SWR matrix symmetry in requirements/README.md.
// Usage: node trace-sync.mjs [--req-dir requirements] [--adr-dir adr] [--write]
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const opt = (name, def) => {
  const i = args.indexOf(name);
  return i >= 0 && args[i + 1] !== undefined ? args[i + 1] : def;
};
const reqDir = path.resolve(opt('--req-dir', 'requirements'));
const adrDir = path.resolve(opt('--adr-dir', 'adr'));
const write = args.includes('--write');

const adrFiles = fs.existsSync(adrDir)
  ? fs.readdirSync(adrDir).filter((f) => /^ADR-\d+.*\.md$/.test(f)).sort()
  : [];
const shrDir = path.join(reqDir, 'shr');
const swrDir = path.join(reqDir, 'swr');
const reqFiles = [
  ...(fs.existsSync(shrDir) ? fs.readdirSync(shrDir).filter((f) => /^SHR-\d+.*\.md$/.test(f)).map((f) => path.join(shrDir, f)) : []),
  ...(fs.existsSync(swrDir) ? fs.readdirSync(swrDir).filter((f) => /^SWR-\d+.*\.md$/.test(f)).map((f) => path.join(swrDir, f)) : []),
].sort();

const adrMap = new Map();
const reqMap = new Map();
for (const f of adrFiles) {
  const id = /^(ADR-\d+)/.exec(f)[1];
  const text = fs.readFileSync(path.join(adrDir, f), 'utf8');
  const m = /^\*\*Satisfies:\*\*\s*(.+)$/m.exec(text);
  const set = new Set();
  if (m) for (const r of m[1].matchAll(/\b(SHR-\d{3}|SWR-\d{3})\b/g)) set.add(r[1]);
  adrMap.set(id, { file: path.join(adrDir, f), set });
}
for (const f of reqFiles) {
  const id = /^(SHR-\d+|SWR-\d+)/.exec(path.basename(f))[1];
  const text = fs.readFileSync(f, 'utf8');
  const m = /^Implemented by (ADR-[^\r\n]+)$/m.exec(text);
  const set = new Set();
  if (m) for (const a of m[1].matchAll(/\b(ADR-\d{3})\b/g)) set.add(a[1]);
  reqMap.set(id, { file: f, set });
}

const asymmetries = [];
for (const [adr, { set }] of adrMap) {
  for (const req of set) {
    if (reqMap.has(req) && !reqMap.get(req).set.has(adr)) {
      asymmetries.push(`${adr} -> ${req}: є в шапці ${adr}, але відсутнє в Implemented by у ${req}`);
      if (write) reqMap.get(req).set.add(adr);
    }
  }
}
for (const [req, { set }] of reqMap) {
  for (const adr of set) {
    if (adrMap.has(adr) && !adrMap.get(adr).set.has(req)) {
      asymmetries.push(`${req} -> ${adr}: є в Implemented by у ${req}, але відсутнє в шапці ${adr}`);
      if (write) adrMap.get(adr).set.add(req);
    }
  }
}

if (write && asymmetries.length > 0) {
  for (const [, { file, set }] of adrMap) {
    if (!set.size) continue;
    const text = fs.readFileSync(file, 'utf8').replace(
      /^\*\*Satisfies:\*\*\s*.+$/m,
      `**Satisfies:** ${[...set].sort().join(', ')}`,
    );
    fs.writeFileSync(file, text);
  }
  for (const [, { file, set }] of reqMap) {
    if (!set.size) continue;
    let text = fs.readFileSync(file, 'utf8');
    const line = `Implemented by ${[...set].sort().join(', ')}.`;
    if (/^Implemented by ADR-[^\r\n]+$/m.test(text)) {
      text = text.replace(/^Implemented by ADR-[^\r\n]+$/m, line);
    } else {
      text = `${text.trimEnd()}\n\n${line}\n`;
    }
    fs.writeFileSync(file, text);
  }
}

console.log(`ADR=${adrMap.size} REQ=${reqMap.size} asymmetries=${asymmetries.length} mode=${write ? 'write' : 'check'}`);
for (const a of asymmetries) console.log(`- ${a}`);
process.exitCode = !write && asymmetries.length > 0 ? 1 : 0;
