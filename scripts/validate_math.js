#!/usr/bin/env node

/**
 * validate_math.js
 * 
 * Scans markdown files for KaTeX syntax errors, legacy GitLab-style math delimiters ($`...`$),
 * invalid spacing before superscripts (\\,^\\circ), and unclosed math blocks.
 * 
 * Usage:
 *   node scripts/validate_math.js [--dir <dir>] [--file <file>] [--fix]
 */

const fs = require('fs');
const path = require('path');

// Resolve KaTeX library
let katex;
try {
  katex = require('katex');
} catch (e) {
  try {
    katex = require('/usr/local/lib/node_modules/@deepseek-ai/dsh/node_modules/katex');
  } catch (e2) {
    console.error('Error: Could not load KaTeX library. Please ensure katex is installed.');
    process.exit(1);
  }
}

// Parse command line arguments
const args = process.argv.slice(2);
let targetDir = path.resolve(__dirname, '..', 'ExpertSystem');
let targetFile = null;
let autoFix = false;

for (let i = 0; i < args.length; i++) {
  if (args[i] === '--dir' && args[i + 1]) {
    targetDir = path.resolve(args[i + 1]);
    i++;
  } else if (args[i] === '--file' && args[i + 1]) {
    targetFile = path.resolve(args[i + 1]);
    i++;
  } else if (args[i] === '--fix') {
    autoFix = true;
  }
}

// Find all .md files recursively
function getMarkdownFiles(dir) {
  let results = [];
  const list = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of list) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== '.git' && entry.name !== 'node_modules') {
        results = results.concat(getMarkdownFiles(fullPath));
      }
    } else if (entry.isFile() && entry.name.endsWith('.md')) {
      results.push(fullPath);
    }
  }
  return results;
}

const filesToScan = targetFile ? [targetFile] : getMarkdownFiles(targetDir);

let totalFiles = 0;
let filesWithErrors = 0;
let totalLegacyCount = 0;
let totalSpacingCount = 0;
let totalKatexErrors = 0;

console.log(`\n======================================================`);
console.log(` KaTeX & Math Syntax Validator`);
console.log(` Scanning ${filesToScan.length} file(s)...`);
if (autoFix) console.log(` Mode: AUTO-FIX enabled`);
console.log(`======================================================\n`);

for (const filePath of filesToScan) {
  totalFiles++;
  const relPath = path.relative(path.resolve(__dirname, '..'), filePath);
  let content = fs.readFileSync(filePath, 'utf8');
  let lines = content.split('\n');
  let fileIssues = [];
  let fileModified = false;

  let inFencedCode = false;
  let codeFenceType = '';
  let inDisplayMath = false;
  let displayMathStartLine = 0;
  let displayMathBuffer = [];

  for (let idx = 0; idx < lines.length; idx++) {
    const lineNo = idx + 1;
    let line = lines[idx];

    // Check for fenced code block toggle
    const fenceMatch = line.match(/^(\s*)(```+|~~~+)(.*)$/);
    if (fenceMatch) {
      const fenceTag = fenceMatch[3].trim().toLowerCase();
      if (!inFencedCode) {
        inFencedCode = true;
        codeFenceType = fenceTag;
        if (codeFenceType === 'math') {
          inDisplayMath = true;
          displayMathStartLine = lineNo;
          displayMathBuffer = [];
          continue;
        }
      } else {
        if (codeFenceType === 'math') {
          inDisplayMath = false;
          const fullMath = displayMathBuffer.join('\n');
          try {
            katex.renderToString(fullMath, { throwOnError: true, displayMode: true });
          } catch (e) {
            fileIssues.push({
              line: displayMathStartLine,
              type: 'KATEX_PARSE_ERROR (```math)',
              detail: e.message,
              snippet: displayMathBuffer[0] || ''
            });
            totalKatexErrors++;
          }
        }
        inFencedCode = false;
        codeFenceType = '';
        continue;
      }
    }

    if (inFencedCode) {
      if (codeFenceType === 'math') {
        displayMathBuffer.push(line);
      }
      continue; // Skip non-math code blocks
    }

    // Check for $$ display math blocks
    const trimmed = line.trim();
    if (trimmed === '$$') {
      if (!inDisplayMath) {
        inDisplayMath = true;
        displayMathStartLine = lineNo;
        displayMathBuffer = [];
        continue;
      } else {
        inDisplayMath = false;
        const fullMath = displayMathBuffer.join('\n');
        try {
          katex.renderToString(fullMath, { throwOnError: true, displayMode: true });
        } catch (e) {
          fileIssues.push({
            line: displayMathStartLine,
            type: 'KATEX_PARSE_ERROR ($$)',
            detail: e.message,
            snippet: displayMathBuffer[0] || ''
          });
          totalKatexErrors++;
        }
        continue;
      }
    }

    if (inDisplayMath) {
      displayMathBuffer.push(line);
      continue;
    }

    // 1. Detect legacy GitLab math: $`...`$
    const legacyMatches = line.match(/\$`([^`\n]+)`\$/g);
    if (legacyMatches) {
      for (const m of legacyMatches) {
        fileIssues.push({
          line: lineNo,
          type: 'LEGACY_GITLAB_MATH',
          detail: 'Enclosed in backticks ($`...`$) which crashes KaTeX with "unknown type: internal"',
          snippet: m
        });
        totalLegacyCount++;
      }
      if (autoFix) {
        line = line.replace(/\$`([^`\n]+)`\$/g, '$$$1$$');
        lines[idx] = line;
        fileModified = true;
      }
    }

    // 2. Detect invalid thin space before superscript: \,^\circ or \,^
    const spacingMatches = line.match(/\\,\s*\^\s*(\\?circ|[a-zA-Z0-9])/g);
    if (spacingMatches) {
      for (const m of spacingMatches) {
        fileIssues.push({
          line: lineNo,
          type: 'INVALID_SPACING_EXPONENT',
          detail: 'Superscript attached to thin space (\\,^) crashes KaTeX with "unknown type: internal"',
          snippet: m
        });
        totalSpacingCount++;
      }
      if (autoFix) {
        line = line.replace(/\\,\s*\^/g, '^');
        lines[idx] = line;
        fileModified = true;
      }
    }

    // 3. Check inline math: $...$
    // Match inline math not preceded or followed by $
    const inlineRegex = /(?<!\$)\$(?!\$)(.*?)(?<!\$)\$(?!\$)/g;
    let match;
    while ((match = inlineRegex.exec(line)) !== null) {
      const mathExpr = match[1];
      if (!mathExpr.trim()) continue;
      // Skip if it contains legacy backticks already handled
      if (mathExpr.startsWith('`') && mathExpr.endsWith('`')) continue;

      try {
        katex.renderToString(mathExpr, { throwOnError: true, displayMode: false });
      } catch (e) {
        fileIssues.push({
          line: lineNo,
          type: 'KATEX_PARSE_ERROR (inline $)',
          detail: e.message,
          snippet: `$${mathExpr}$`
        });
        totalKatexErrors++;
      }
    }
  }

  if (fileIssues.length > 0) {
    filesWithErrors++;
    console.log(`\n❌ ${relPath} (${fileIssues.length} issue${fileIssues.length > 1 ? 's' : ''}):`);
    for (const issue of fileIssues) {
      console.log(`   [L${issue.line}] [${issue.type}] ${issue.snippet}`);
      console.log(`         -> ${issue.detail}`);
    }
  }

  if (autoFix && fileModified) {
    fs.writeFileSync(filePath, lines.join('\n'), 'utf8');
    console.log(`   ✅ Automatically fixed issues in ${relPath}`);
  }
}

console.log(`\n======================================================`);
console.log(` Summary:`);
console.log(` Total files checked:        ${totalFiles}`);
console.log(` Files with issues:          ${filesWithErrors}`);
console.log(` Legacy GitLab math ($` + '`):  ' + totalLegacyCount);
console.log(` Invalid \\,^ spacing:        ${totalSpacingCount}`);
console.log(` KaTeX parse errors:         ${totalKatexErrors}`);
console.log(`======================================================\n`);

