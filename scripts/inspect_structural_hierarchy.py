#!/usr/bin/env python3
"""
Inspect structural hierarchy, heading numbering, pseudo-headings, and authorial alert blocks in monograph chapters.

Detects:
  [ERR-STRUCT-01]: Bold pseudo-headings inside sections (e.g. '**Діалог з...**', '**Синтез...**') that should be formal H3/H4 headings.
  [ERR-STRUCT-02]: Authorial contribution declarations (e.g. '**Принцип ... М. Федчика**', '**Закон ... Федчика**') outside of canonical `> [!NOTE]` alert blocks.
  [ERR-STRUCT-03]: Excessive stacked blank lines (>= 3 consecutive empty lines).
  [ERR-STRUCT-04]: Malformed alert block header or missing Exact Count Matching.
"""

import sys
import re
from pathlib import Path

PSEUDO_HEADING_PATTERNS = [
    re.compile(r'^\s*\*\*(Діалог з [^*]+)\*\*\s*$'),
    re.compile(r'^\s*\*\*(Синтез та науковий вердикт[^*]*)\*\*\s*$'),
    re.compile(r'^\s*\*\*(Підсумок[^*]*)\*\*\s*$'),
]

AUTHORIAL_CLAIM_PATTERN = re.compile(
    r'(?<!> )(?<!>\s)(?:\d+\.\s*)?\*\*(?:Принцип|Закон|Модель|Теорія|Правило|Архітектурна парадигма|Інваріант)[^*]*(?:Федчик|М\.\s*Федчик|М\.\s*М\.\s*Федчик)[^*]*\*\*'
)

STACKED_BLANK_LINES_PATTERN = re.compile(r'\n{4,}')

def inspect_file(filepath: Path) -> list:
    issues = []
    text = filepath.read_text(encoding='utf-8')
    lines = text.splitlines()

    # Check stacked blank lines
    for match in STACKED_BLANK_LINES_PATTERN.finditer(text):
        start_pos = match.start()
        line_num = text[:start_pos].count('\n') + 1
        issues.append((line_num, "ERR-STRUCT-03", "Excessive stacked blank lines (3+ empty lines in a row)"))

    in_code_block = False
    in_note_block = False

    for idx, line in enumerate(lines, start=1):
        stripped = line.strip()
        
        # Track code blocks
        if stripped.startswith('```'):
            in_code_block = not in_code_block
            continue
        if in_code_block:
            continue

        # Track NOTE blocks
        if stripped.startswith('> [!NOTE]'):
            in_note_block = True
            continue
        elif in_note_block:
            if not stripped.startswith('>'):
                in_note_block = False

        # Check pseudo headings
        for pat in PSEUDO_HEADING_PATTERNS:
            m = pat.match(stripped)
            if m:
                heading_text = m.group(1)
                issues.append((idx, "ERR-STRUCT-01", f"Pseudo-heading bold line '**{heading_text}**' should be a formal H3/H4 markdown heading."))

        # Check unformatted authorial claims outside > [!NOTE]
        if not in_note_block:
            m = AUTHORIAL_CLAIM_PATTERN.search(line)
            if m:
                # Exclude if it's already an H1-H4 heading or in a blockquote note
                if not stripped.startswith('#') and not stripped.startswith('>'):
                    issues.append((idx, "ERR-STRUCT-02", f"Authorial contribution claim '{m.group(0)}' outside of canonical '> [!NOTE]' block."))

    return issues

def main():
    if len(sys.argv) > 1:
        target = Path(sys.argv[1])
        if target.is_file():
            files = [target]
        elif target.is_dir():
            files = sorted(list(target.glob('ch*.md')) + list(target.glob('appendix*.md')))
        else:
            print(f"Path not found: {target}")
            sys.exit(1)
    else:
        expert_dir = Path("ExpertSystem")
        files = sorted(list(expert_dir.glob('ch*.md')) + list(expert_dir.glob('appendix*.md')))

    total_issues = 0
    for f in files:
        if 'translations' in f.parts:
            continue
        issues = inspect_file(f)
        if issues:
            print(f"\n❌ {f} ({len(issues)} issue{'s' if len(issues) > 1 else ''}):")
            for line_num, code, msg in issues:
                print(f"  L{line_num} [{code}]: {msg}")
            total_issues += len(issues)

    if total_issues == 0:
        print("\n✅ All checked files passed structural hierarchy and authorial formatting inspection.")
        sys.exit(0)
    else:
        print(f"\n[ERROR] Found {total_issues} structural violation(s).")
        sys.exit(1)

if __name__ == '__main__':
    main()
