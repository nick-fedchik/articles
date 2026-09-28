#!/usr/bin/env bash
set -euo pipefail

echo "===================================================="
echo "Running Full Validation Pipeline for MilTech & Repo"
echo "===================================================="

echo "[1/4] Running markdownlint-cli2..."
npx markdownlint-cli2 --config ExpertSystem/.markdownlint-cli2.jsonc "MilTech/*.md"

echo "[2/4] Running KaTeX GFM Math syntax validator..."
node scripts/validate_math.js MilTech

echo "[3/4] Running GitHub Markdown AST & Mermaid Parser..."
(cd MilTech && node ../ExpertSystem/scripts/validate-github-markdown.mjs)

echo "[4/4] Running GitLab Markdown AST & KaTeX Parser..."
(cd MilTech && node ../ExpertSystem/scripts/validate-gitlab-markdown.mjs)

echo "===================================================="
echo "SUCCESS: All Markdown, Math, and Mermaid checks passed!"
echo "===================================================="
