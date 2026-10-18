#!/usr/bin/env python3
"""
Confidentiality & Internal Project Artifact Leakage Inspector for Monograph.

Ensures that no proprietary codebase artifacts, internal file extensions,
private repository names, or implementation-specific file paths leak into
the public scientific/engineering monograph text.

Enforces project invariant:
  Monograph text must use canonical, generalized architectural and engineering
  terminology (e.g. 'бінарний граф онтології', 'черга дослідницьких прогалин')
  rather than internal implementation artifact names ('ontology.zkpg', 'spool/gaps', 'vida.*').

Usage:
  python3 scripts/inspect_project_leakage.py [path ...]
"""

import os
import re
import sys
from pathlib import Path
from typing import List, Tuple

# Patterns of proprietary/internal project artifacts that must NOT appear in the public monograph
LEAKAGE_PATTERNS = [
    (re.compile(r"\bontology\.zkpg\b", re.IGNORECASE), "Replace with 'бінарний граф онтології'"),
    (re.compile(r"\bzkpg\b", re.IGNORECASE), "Replace with canonical ontology/graph terminology"),
    (re.compile(r"spool/gaps\b", re.IGNORECASE), "Replace with 'черга дослідницьких прогалин (Knowledge Gaps Queue)'"),
    (re.compile(r"\bvida\.(tasks|gap|events|bus)\b", re.IGNORECASE), "Replace with general domain prefix like 'epistemic.*' or 'es.*'"),
    (re.compile(r"~/vida\b", re.IGNORECASE), "Internal local path leaked; replace with generalized repository/architecture path"),
    (re.compile(r"/home/\w+/vida\b", re.IGNORECASE), "Internal absolute path leaked; replace with generalized architecture path"),
]

# Files or directories to skip (e.g. this inspection script itself, git directories)
IGNORED_PATHS = {
    "inspect_project_leakage.py",
    ".git",
    "node_modules",
}


def find_files(targets: List[str]) -> List[Path]:
    result = []
    for target in targets:
        p = Path(target)
        if p.is_file() and p.suffix == ".md":
            result.append(p)
        elif p.is_dir():
            for root, dirs, files in os.walk(p):
                # Skip translations if needed or git directories
                if any(ignored in root for ignored in IGNORED_PATHS):
                    continue
                for f in sorted(files):
                    if f.endswith(".md"):
                        result.append(Path(root) / f)
    return sorted(set(result))


def inspect_file(filepath: Path) -> List[Tuple[int, str, str, str]]:
    issues = []
    try:
        content = filepath.read_text(encoding="utf-8")
    except Exception as e:
        print(f"[ERROR] Could not read {filepath}: {e}", file=sys.stderr)
        return []

    lines = content.splitlines()
    for lineno, line in enumerate(lines, start=1):
        for pattern, recommendation in LEAKAGE_PATTERNS:
            match = pattern.search(line)
            if match:
                matched_text = match.group(0)
                issues.append((lineno, matched_text, line.strip(), recommendation))

    return issues


def main():
    targets = sys.argv[1:] if len(sys.argv) > 1 else ["ExpertSystem"]
    files = find_files(targets)

    if not files:
        print(f"[WARN] No markdown files found in target(s): {targets}")
        sys.exit(0)

    total_violations = 0
    print(f"=== Inspecting {len(files)} Markdown files for internal project artifact leaks ===")

    for fpath in files:
        issues = inspect_file(fpath)
        if issues:
            print(f"\n[LEAK] {fpath}:")
            for lineno, matched, line, rec in issues:
                total_violations += 1
                print(f"  L {lineno}: '{matched}' -> {rec}")
                print(f"    Line snippet: {line[:120]}")

    if total_violations > 0:
        print(f"\n[ERROR] Found {total_violations} proprietary artifact leakage violation(s)!")
        sys.exit(1)
    else:
        print("\n[OK] Clean! Zero proprietary artifact leaks detected across all files.")
        sys.exit(0)


if __name__ == "__main__":
    main()
