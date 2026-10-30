// Audit tool: instant crash-recovery and progress report from {{AUDIT}} (STATUS.md, RUNLOG.md, PROGRESS.md, artifacts).
// Usage: node status-report.mjs [audit-dir=audit/work] [--json]
import fs from 'node:fs';
import path from 'node:path';

const args = process.argv.slice(2);
const auditDir = path.resolve(args.find((a) => !a.startsWith('--')) ?? 'audit/work');
const asJson = args.includes('--json');

if (!fs.existsSync(auditDir)) {
  console.error(`Каталог аудиту не знайдено: ${auditDir}`);
  process.exit(2);
}

const walk = (dir, filter = () => true, acc = []) => {
  if (!fs.existsSync(dir)) return acc;
  for (const e of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, e.name);
    if (e.isDirectory()) walk(p, filter, acc);
    else if (filter(p)) acc.push(p);
  }
  return acc;
};

// Parse RUNLOG.md for total runs, completed runs, and interrupted ("в роботі") runs
const runlogPath = path.join(auditDir, 'RUNLOG.md');
const runs = [];
const activeRuns = [];
if (fs.existsSync(runlogPath)) {
  for (const line of fs.readFileSync(runlogPath, 'utf8').split(/\r?\n/)) {
    const m = /^\|\s*(RUN-\d+[a-z]?)\s*\|\s*([^|]+)\|\s*([^|]+)\|\s*([^|]+)\|\s*([^|]+)\|\s*([^|]+)\|\s*([^|]+)\|\s*([^|]+)\|/.exec(line);
    if (!m) continue;
    const item = {
      id: m[1].trim(),
      time: m[2].trim(),
      phasePkg: m[3].trim(),
      actor: m[4].trim(),
      calls: m[6].trim(),
      output: m[7].trim(),
      verified: m[8].trim(),
    };
    runs.push(item);
    if (/в роботі/i.test(item.time) || /в роботі/i.test(item.verified)) activeRuns.push(item);
  }
}

// Parse STATUS.md for phase table
const statusPath = path.join(auditDir, 'STATUS.md');
const phases = [];
if (fs.existsSync(statusPath)) {
  let inPhases = false;
  for (const line of fs.readFileSync(statusPath, 'utf8').split(/\r?\n/)) {
    if (/^##\s+Фази/.test(line)) { inPhases = true; continue; }
    if (inPhases && /^##\s+/.test(line)) break;
    if (!inPhases || !line.startsWith('|') || /^\|[\s:-]+\|/.test(line) || /\|\s*Фаза\s*\|/.test(line)) continue;
    const cols = line.split('|').slice(1, -1).map((c) => c.trim());
    if (cols.length >= 6) {
      phases.push({
        phase: cols[0],
        state: cols[1].replace(/\*\*/g, ''),
        units: cols[2],
        passes: cols[3],
        time: cols[4],
        checkpoint: cols[5].replace(/\*\*/g, ''),
      });
    }
  }
}

// Count artifacts on disk
const reqFiles = walk(path.join(auditDir, 'requirements'), (p) => p.endsWith('.md'));
let rqCount = 0;
const unfinishedSources = [];
for (const f of reqFiles) {
  const txt = fs.readFileSync(f, 'utf8');
  if (/^- \*\*Стан:\*\*.*в роботі/m.test(txt)) unfinishedSources.push(path.relative(auditDir, f).replace(/\\/g, '/'));
  rqCount += (txt.match(/^## RQ-\d+/gm) ?? []).length;
}
const findingFiles = walk(path.join(auditDir, 'findings'), (p) => p.endsWith('.md'));
const lensFiles = walk(path.join(auditDir, 'lenses'), (p) => p.endsWith('.md'));
const covFiles = walk(path.join(auditDir, 'coverage'), (p) => p.endsWith('.md'));
const taskFiles = walk(path.join(auditDir, 'tasks'), (p) => p.endsWith('.md'));
const tmpDirs = fs.existsSync(path.join(auditDir, 'tmp'))
  ? fs.readdirSync(path.join(auditDir, 'tmp')).map((n) => `tmp/${n}`)
  : [];

const donePhases = phases.filter((p) => /перевірено|завершено|не потрібен/i.test(p.state)).length;
const pctPhases = phases.length ? Math.round((donePhases / phases.length) * 100) : 0;

const summary = {
  auditDir: path.relative(process.cwd(), auditDir).replace(/\\/g, '/') || '.',
  phasesTotal: phases.length,
  phasesDone: donePhases,
  phaseProgressPct: pctPhases,
  runsTotal: runs.length,
  runsInProgressOrInterrupted: activeRuns.map((r) => `${r.id} (${r.phasePkg})`),
  lastCompletedRun: [...runs].reverse().find((r) => !activeRuns.includes(r))?.id ?? null,
  artifacts: {
    extractedSources: reqFiles.length,
    extractedClaimsRQ: rqCount,
    unfinishedSources,
    lensLogs: lensFiles.length,
    findings: findingFiles.length,
    coverageLogs: covFiles.length,
    backlogTasks: taskFiles.length,
    hasConcept: fs.existsSync(path.join(auditDir, 'CONCEPT.md')),
    hasImplementation: fs.existsSync(path.join(auditDir, 'IMPLEMENTATION.md')),
    hasVerification: fs.existsSync(path.join(auditDir, 'VERIFICATION.md')),
    hasRegister: fs.existsSync(path.join(auditDir, 'REGISTER.md')),
    hasReport: fs.existsSync(path.join(auditDir, 'REPORT.md')),
    hasCoverageSummary: fs.existsSync(path.join(auditDir, 'COVERAGE-CODE.md')),
    hasProgressFeed: fs.existsSync(path.join(auditDir, 'PROGRESS.md')),
  },
  leftoverTmp: tmpDirs,
  phases,
};

if (asJson) {
  console.log(JSON.stringify(summary, null, 2));
} else {
  console.log(`# Діагностика стану аудиту (${summary.auditDir})\n`);
  console.log(`- Прогрес фаз: ${summary.phasesDone} / ${summary.phasesTotal} (${summary.phaseProgressPct}%)`);
  console.log(`- Запусків у RUNLOG.md: ${summary.runsTotal} (останній завершений: ${summary.lastCompletedRun ?? 'немає'})`);
  console.log(`- Перервані / активні запуски («в роботі»): ${summary.runsInProgressOrInterrupted.length ? summary.runsInProgressOrInterrupted.join(', ') : 'немає'}`);
  console.log(`- Залишки у tmp/: ${summary.leftoverTmp.length ? summary.leftoverTmp.join(', ') : 'немає'}`);
  console.log(`- Видобуто: ${summary.artifacts.extractedSources} джерел, ${summary.artifacts.extractedClaimsRQ} RQ (незавершених джерел: ${unfinishedSources.length})`);
  console.log(`- Артефакти: журналів лінз=${lensFiles.length}, знахідок F-*=${findingFiles.length}, журналів покриття=${covFiles.length}, задач=${taskFiles.length}`);
  console.log(`- Ключові звіти: CONCEPT=${summary.artifacts.hasConcept}, IMPLEMENTATION=${summary.artifacts.hasImplementation}, REGISTER=${summary.artifacts.hasRegister}, REPORT=${summary.artifacts.hasReport}, COVERAGE-CODE=${summary.artifacts.hasCoverageSummary}\n`);
  if (phases.length) {
    console.log('| Фаза | Стан | Одиниці | Проходи | Час | Контрольна точка |');
    console.log('| --- | --- | --- | --- | --- | --- |');
    for (const p of phases) console.log(`| ${p.phase} | ${p.state} | ${p.units} | ${p.passes} | ${p.time} | ${p.checkpoint} |`);
  }
}
process.exitCode = activeRuns.length || unfinishedSources.length ? 1 : 0;
