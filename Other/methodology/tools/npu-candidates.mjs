#!/usr/bin/env node
/**
 * npu-candidates.mjs (SWR-034, RM-006)
 *
 * Queries the local NPU Forge engine (`POST /v1/audit/candidates` on loopback)
 * to generate a bounded top-k candidate matrix for audit lenses C (conflicts),
 * S/T (cascades), or Phase 10 (requirement-to-code coverage), and optionally
 * writes the JSON + Markdown candidate report into `{{AUDIT}}/scan/`.
 *
 * Usage:
 *   node audit/methodology/tools/npu-candidates.mjs \
 *     --mode <conflicts|cascades|coverage> \
 *     --sources <dir-or-file[,...]> \
 *     --targets <dir-or-file[,...]> \
 *     [--top-k 5] [--min-score 0.35] [--url http://127.0.0.1:3939] [--out audit/work/scan/candidates.json]
 */

import fs from "node:fs";
import path from "node:path";
import http from "node:http";

function parseArgs(argv) {
  const opts = {
    mode: "coverage",
    sources: [],
    targets: [],
    topK: 5,
    minScore: 0.35,
    excludeSameFile: true,
    url: "http://127.0.0.1:3939",
    workspaceRoot: process.cwd(),
    out: "",
  };
  for (let i = 2; i < argv.length; i++) {
    const arg = argv[i];
    if (arg === "--mode" && argv[i + 1]) opts.mode = argv[++i];
    else if (arg === "--sources" && argv[i + 1]) opts.sources = argv[++i].split(",").map((s) => s.trim()).filter(Boolean);
    else if (arg === "--targets" && argv[i + 1]) opts.targets = argv[++i].split(",").map((s) => s.trim()).filter(Boolean);
    else if (arg === "--top-k" && argv[i + 1]) opts.topK = Number(argv[++i]);
    else if (arg === "--min-score" && argv[i + 1]) opts.minScore = Number(argv[++i]);
    else if (arg === "--include-same-file") opts.excludeSameFile = false;
    else if (arg === "--url" && argv[i + 1]) opts.url = argv[++i];
    else if (arg === "--root" && argv[i + 1]) opts.workspaceRoot = path.resolve(argv[++i]);
    else if (arg === "--out" && argv[i + 1]) opts.out = argv[++i];
  }
  return opts;
}

const ALLOWED_EXTS = new Set([".md", ".cpp", ".hpp", ".c", ".cc", ".h", ".ts", ".js", ".py"]);
const SKIP_DIRS = new Set([".git", "node_modules", "build", "build-release-gate", "out", ".npu-forge", "dist"]);

function collectRelativeFiles(rootDir, specs) {
  const results = new Set();
  function walk(absPath) {
    let stat;
    try {
      stat = fs.statSync(absPath);
    } catch {
      return;
    }
    if (stat.isFile()) {
      const ext = path.extname(absPath).toLowerCase();
      if (ALLOWED_EXTS.has(ext)) {
        results.add(path.relative(rootDir, absPath).replace(/\\/g, "/"));
      }
      return;
    }
    if (stat.isDirectory()) {
      for (const entry of fs.readdirSync(absPath, { withFileTypes: true })) {
        if (entry.isDirectory() && SKIP_DIRS.has(entry.name)) continue;
        walk(path.join(absPath, entry.name));
      }
    }
  }
  for (const spec of specs) {
    walk(path.resolve(rootDir, spec));
  }
  return [...results].sort();
}

function postJson(rawUrl, endpoint, payload) {
  const parsed = new URL(endpoint, rawUrl);
  if (!["127.0.0.1", "localhost", "[::1]", "::1"].includes(parsed.hostname.toLowerCase())) {
    return Promise.reject(new Error(`Refusing non-loopback engine URL: ${rawUrl} (SHR-003)`));
  }
  const body = JSON.stringify(payload);
  return new Promise((resolve, reject) => {
    const req = http.request(
      parsed,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Content-Length": Buffer.byteLength(body),
        },
      },
      (res) => {
        let data = "";
        res.on("data", (chunk) => (data += chunk));
        res.on("end", () => resolve({ statusCode: res.statusCode ?? 0, body: data }));
      },
    );
    req.on("error", reject);
    req.write(body);
    req.end();
  });
}

async function main() {
  const opts = parseArgs(process.argv);
  if (opts.sources.length === 0 || opts.targets.length === 0) {
    console.error(
      "Usage: node audit/methodology/tools/npu-candidates.mjs --mode <conflicts|cascades|coverage> --sources <dir,...> --targets <dir,...> [--top-k 5] [--min-score 0.35] [--out <path>]",
    );
    process.exit(1);
  }

  const sourcePaths = collectRelativeFiles(opts.workspaceRoot, opts.sources);
  const targetPaths = collectRelativeFiles(opts.workspaceRoot, opts.targets);
  if (sourcePaths.length === 0 || targetPaths.length === 0) {
    console.error(`No files found for sources (${sourcePaths.length}) or targets (${targetPaths.length}).`);
    process.exit(1);
  }

  const response = await postJson(opts.url, "/v1/audit/candidates", {
    workspace_root: opts.workspaceRoot,
    mode: opts.mode,
    source_paths: sourcePaths,
    target_paths: targetPaths,
    top_k: opts.topK,
    min_score: opts.minScore,
    exclude_same_file: opts.excludeSameFile,
  });

  if (response.statusCode !== 200) {
    console.error(`HTTP ${response.statusCode}: ${response.body}`);
    process.exit(1);
  }

  const report = JSON.parse(response.body);
  console.log(
    `mode=${report.mode} device=${report.device} source_chunks=${report.source_chunk_count} target_chunks=${report.target_chunk_count} rows=${report.rows.length} uncovered=${report.uncovered_sources.length} elapsed_ms=${Number(report.elapsed_ms).toFixed(1)} digest=${report.report_digest}`,
  );

  if (opts.out) {
    const outPath = path.resolve(opts.workspaceRoot, opts.out);
    fs.mkdirSync(path.dirname(outPath), { recursive: true });
    fs.writeFileSync(outPath, JSON.stringify(report, null, 2) + "\n", "utf8");
    console.log(`Saved candidate matrix to ${outPath}`);
  }
}

main().catch((err) => {
  console.error(`npu-candidates failed: ${err.message || err}`);
  process.exit(1);
});
