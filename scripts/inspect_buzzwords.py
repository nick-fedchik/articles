#!/usr/bin/env python3
"""
Style & AI Marker Inspection Tool for 'Architecture of Evidence-Governed Expert Systems'.
Inspects text files (Markdown) for unnatural/machine buzzwords, specifically the word 'контур' (kontur)
which AI models frequently overuse where natural Ukrainian requires 'система', 'середовище',
'шлюз', 'цикл керування', 'тракт', 'ланцюг', etc.
"""

import sys
import re
from pathlib import Path

# Patterns of "контур" that are suspicious / machine markers in Ukrainian tech writing
SUSPICIOUS_PATTERNS = [
    re.compile(r"контур\w*\s+безпеки", re.IGNORECASE),
    re.compile(r"безпеков\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"двоконтурн\w*", re.IGNORECASE),
    re.compile(r"триконтурн\w*", re.IGNORECASE),
    re.compile(r"апаратн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"контур\w*\s+керування", re.IGNORECASE),
    re.compile(r"контур\w*\s+регулювання", re.IGNORECASE),
    re.compile(r"в\s+контурі\b", re.IGNORECASE),
    re.compile(r"у\s+контурі\b", re.IGNORECASE),
    re.compile(r"людин\w*\s+в\s+контурі", re.IGNORECASE),
    re.compile(r"контур\w*\s+допуску", re.IGNORECASE),
    re.compile(r"контур\w*\s+(?:до)?навчання", re.IGNORECASE),
    re.compile(r"оперативн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"виконавч\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"контур\w*\s+інженерн\w*", re.IGNORECASE),
    re.compile(r"контур\w*\s+сприйняття", re.IGNORECASE),
    re.compile(r"контур\w*\s+міркування", re.IGNORECASE),
    re.compile(r"контур\w*\s+арбітражу", re.IGNORECASE),
    re.compile(r"контур\w*\s+W3C", re.IGNORECASE),
    re.compile(r"обладнанням\s+у\s+контурі", re.IGNORECASE),
    re.compile(r"апаратурою\s+в\s+контурі", re.IGNORECASE),
    re.compile(r"контур\w*\s+зворотного\s+зв'язку", re.IGNORECASE),
    re.compile(r"карантинн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"контур\w*\s+вектор\w*", re.IGNORECASE),
    re.compile(r"лексичн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"семантичн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"контур\s+\d+", re.IGNORECASE),
    re.compile(r"цифров\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"нейроморфн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"сенсорн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"технологічн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"інженерн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"сертифікаційн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"аналітичн\w*\s+контур\w*", re.IGNORECASE),
    re.compile(r"кіберфізичн\w*\s+контур\w*", re.IGNORECASE),
]

# Patterns where "контур" is legitimate in Ukrainian technical writing:
# 1. Geometry, outlines, topographic/contour maps: контурна карта, дескриптор контуру, контур фігури
# 2. Electrical ground loops / physical grounding: контур заземлення
LEGITIMATE_CONTEXTS = [
    re.compile(r"контурн\w*\s+карт\w*", re.IGNORECASE),
    re.compile(r"контурн\w*\s+дескриптор\w*", re.IGNORECASE),
    re.compile(r"контур\w*\s+заземлення", re.IGNORECASE),
    re.compile(r"контур\s+або\s+подію", re.IGNORECASE),  # image segmentation
]

def check_file(filepath: Path):
    issues = []
    try:
        text = filepath.read_text(encoding="utf-8")
    except Exception as e:
        return [f"Could not read {filepath}: {e}"]

    # We also do a general catch for any "контур" word
    kontur_word = re.compile(r"\bконтур\w*", re.IGNORECASE)

    for lineno, line in enumerate(text.splitlines(), start=1):
        if not kontur_word.search(line):
            continue

        # Check if it matches legitimate context
        is_legit = any(legit.search(line) for legit in LEGITIMATE_CONTEXTS)
        
        # Check suspicious patterns
        matched_suspicious = [pat.pattern for pat in SUSPICIOUS_PATTERNS if pat.search(line)]
        
        # If it's in AGENTS.md or the script itself, skip checking rules definition
        if "AGENTS.md" in filepath.name or "inspect_buzzwords.py" in filepath.name:
            continue

        if matched_suspicious:
            issues.append((lineno, line.strip(), f"Suspicious marker: {', '.join(matched_suspicious)}"))
        elif not is_legit:
            issues.append((lineno, line.strip(), "Generic 'контур' used; verify if geometric/contour map or machine buzzword"))

    return issues

def main():
    root = Path(__file__).resolve().parent.parent
    target_dir = root / "ExpertSystem"
    if len(sys.argv) > 1:
        target_dir = Path(sys.argv[1])

    files = sorted(target_dir.glob("**/*.md")) if target_dir.is_dir() else [target_dir]
    total_issues = 0

    print(f"=== Inspecting {len(files)} files in {target_dir} for 'контур' AI markers ===")

    for f in files:
        issues = check_file(f)
        if issues:
            rel_path = f.relative_to(root) if f.is_relative_to(root) else f
            print(f"\n[!] {rel_path} ({len(issues)} occurrences):")
            for lineno, line, reason in issues:
                total_issues += 1
                print(f"  Line {lineno}: {line[:120]}")
                print(f"    -> {reason}")

    if total_issues == 0:
        print("\n[OK] No unnatural 'контур' buzzwords found!")
        sys.exit(0)
    else:
        print(f"\n[FAIL] Found {total_issues} occurrences to inspect/fix.")
        sys.exit(1)

if __name__ == "__main__":
    main()
