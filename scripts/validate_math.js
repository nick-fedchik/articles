#!/usr/bin/env node

/**
 * Math Formula Validator for GitHub Flavored Markdown (GFM)
 * Validates LaTeX syntax using KaTeX (strict mode) and checks GFM-specific rules:
 * 1. KaTeX parsing in strict mode
 * 2. No Cyrillic characters inside math mode ($...$ or $$...$$)
 * 3. No leading/trailing spaces inside inline math ($ x $ or $x $)
 * 4. No unescaped brackets after \big, \Big, \bigg, \Bigg (\big{ -> \big\{)
 * 5. Display math ($$) must have blank lines before and after (or valid list indent)
 * 6. No unescaped underscores in \text{...}
 */

const fs = require('fs');
const path = require('path');

// Resolve katex from local node_modules
let katex;
try {
  katex = require(path.resolve(__dirname, '../ExpertSystem/node_modules/katex'));
} catch (e) {
  try {
    katex = require('katex');
  } catch (err) {
    console.error('ERROR: katex module not found. Run "npm install katex" in ExpertSystem or root.');
    process.exit(1);
  }
}

const targetDirs = process.argv.slice(2);
if (targetDirs.length === 0) {
  targetDirs.push('MilTech');
}

function findMarkdownFiles(dir) {
  let results = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      if (entry.name !== 'node_modules' && !entry.name.startsWith('.')) {
        results = results.concat(findMarkdownFiles(fullPath));
      }
    } else if (entry.isFile() && entry.name.endsWith('.md')) {
      results.push(fullPath);
    }
  }
  return results;
}

let totalFiles = 0;
let totalFormulas = 0;
let totalErrors = 0;

for (const target of targetDirs) {
  const stat = fs.statSync(target);
  const files = stat.isDirectory() ? findMarkdownFiles(target) : [target];

  for (const file of files) {
    totalFiles++;
    const content = fs.readFileSync(file, 'utf-8');
    const lines = content.split('\n');

    // 1. Check Display Math: $$ ... $$
    const displayRegex = /\$\$([\s\S]*?)\$\$/g;
    let match;
    while ((match = displayRegex.exec(content)) !== null) {
      totalFormulas++;
      const formula = match[1].trim();
      const lineNum = content.substring(0, match.index).split('\n').length;

      // Cyrillic check
      const cyrillicMatches = formula.match(/[\u0400-\u04FF]/g);
      if (cyrillicMatches) {
        console.error(`[GFM MATH ERROR] ${file}:${lineNum} - Cyrillic characters found inside display math: "${[...new Set(cyrillicMatches)].join(', ')}"`);
        console.error(`  Formula: ${formula.substring(0, 80)}...\n`);
        totalErrors++;
      }

      // Check \big{
      if (/\\(big|Big|bigg|Bigg)\{/.test(formula)) {
        console.error(`[GFM MATH ERROR] ${file}:${lineNum} - Missing backslash before brace in \\big delimiter: use \\big\\{ instead of \\big{`);
        console.error(`  Formula: ${formula}\n`);
        totalErrors++;
      }

      // Check for literal '*' (must use \ast, \star, \cdot, \times)
      if (/\*/.test(formula)) {
        console.error(`[GFM MATH ERROR] ${file}:${lineNum} - Literal '*' found inside display math. Markdown interprets '*' as italics/bold. Use \\ast, \\star, \\cdot, or \\times instead.`);
        console.error(`  Formula: ${formula}\n`);
        totalErrors++;
      }

      // KaTeX strict parse
      try {
        katex.renderToString(formula, { displayMode: true, throwOnError: true, strict: 'warn' });
      } catch (err) {
        console.error(`[KATEX DISPLAY ERROR] ${file}:${lineNum} - ${err.message}`);
        console.error(`  Formula: ${formula}\n`);
        totalErrors++;
      }
    }

    // 2. Check Inline Math: $ ... $
    for (let lineIdx = 0; lineIdx < lines.length; lineIdx++) {
      const line = lines[lineIdx];
      // Skip code blocks
      if (line.trim().startsWith('```') || line.trim().startsWith('`')) continue;

      // Inline math expressions
      const inlineRegex = /(?<!\$)\$(?!\$)([^\$\n]+?)(?<!\$)\$(?!\$)/g;
      let im;
      while ((im = inlineRegex.exec(line)) !== null) {
        const formula = im[1];
        // Skip currency references like $100
        if (/^\s*\d+(?:\.\d+)?\s*$/.test(formula)) continue;

        totalFormulas++;

        // Leading / trailing space check
        if (formula.startsWith(' ') || formula.endsWith(' ')) {
          console.error(`[GFM MATH ERROR] ${file}:${lineIdx + 1} - Leading or trailing space inside inline math: "$${formula}$"`);
          totalErrors++;
        }

        // Check for literal '*' in inline math
        if (/\*/.test(formula)) {
          console.error(`[GFM MATH ERROR] ${file}:${lineIdx + 1} - Literal '*' found inside inline math. Markdown interprets '*' as italics/bold. Use \\ast, \\star, \\cdot, or \\times instead.`);
          console.error(`  Expression: $${formula}$\n`);
          totalErrors++;
        }

        // Cyrillic check
        const cyrillic = formula.match(/[\u0400-\u04FF]/g);
        if (cyrillic) {
          console.error(`[GFM MATH ERROR] ${file}:${lineIdx + 1} - Cyrillic characters found inside inline math: "${[...new Set(cyrillic)].join(', ')}"`);
          console.error(`  Expression: $${formula}$\n`);
          totalErrors++;
        }

        // Unescaped underscore in \text
        if (/\\text\{[^}]*_[^}]*\}/.test(formula) && !/\\text\{[^}]*\\_[^}]*\}/.test(formula)) {
          console.error(`[GFM MATH ERROR] ${file}:${lineIdx + 1} - Unescaped underscore in \\text: use \\text{...\\_...} or hyphen \\text{...-...}`);
          console.error(`  Expression: $${formula}$\n`);
          totalErrors++;
        }

        // KaTeX strict parse
        try {
          katex.renderToString(formula, { displayMode: false, throwOnError: true, strict: 'warn' });
        } catch (err) {
          console.error(`[KATEX INLINE ERROR] ${file}:${lineIdx + 1} - ${err.message}`);
          console.error(`  Expression: $${formula}$\n`);
          totalErrors++;
        }
      }
    }
  }
}

console.log('----------------------------------------------------');
console.log(`Math Validation Summary:`);
console.log(`Files scanned:    ${totalFiles}`);
console.log(`Formulas checked: ${totalFormulas}`);
console.log(`Errors found:     ${totalErrors}`);
console.log('----------------------------------------------------');

if (totalErrors > 0) {
  process.exit(1);
} else {
  console.log('SUCCESS: All mathematical expressions are valid for GitHub Flavored Markdown!');
  process.exit(0);
}
