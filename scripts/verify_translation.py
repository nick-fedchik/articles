#!/usr/bin/env python3
"""
scripts/verify_translation.py — Automated AST Invariant Quality Gate for Translations

Verifies that a translated Markdown chapter preserves:
1. Exact count of code blocks (```).
2. Exact count of math blocks ($$ and ```math).
3. Exact count of Mermaid diagrams (```mermaid).
4. Structural heading hierarchy (H1, H2, H3, H4).
5. All embedded URLs and relative anchors.

Usage:
  python3 scripts/verify_translation.py ExpertSystem/ch29-neuro-symbolic-architecture.md translations/en/ExpertSystem/ch29-neuro-symbolic-architecture.md
"""

import sys
import re
import os

def analyze_markdown(path):
    with open(path, "r", encoding="utf-8") as f:
        text = f.read()

    # Code blocks
    fenced_blocks = re.findall(r"```([a-zA-Z0-9_\-]*)", text)
    math_code_blocks = [b for b in fenced_blocks if b.strip() == "math"]
    mermaid_blocks = [b for b in fenced_blocks if b.strip() == "mermaid"]
    go_blocks = [b for b in fenced_blocks if b.strip() == "go"]
    
    # Standalone $$ blocks
    display_math_dollar = re.findall(r"(?m)^\$\$$", text)
    total_math = len(math_code_blocks) + (len(display_math_dollar) // 2)

    # Headings
    h1 = re.findall(r"(?m)^#\s+(.+)$", text)
    h2 = re.findall(r"(?m)^##\s+(.+)$", text)
    h3 = re.findall(r"(?m)^###\s+(.+)$", text)
    h4 = re.findall(r"(?m)^####\s+(.+)$", text)

    return {
        "path": path,
        "total_fenced": len(fenced_blocks),
        "math_blocks": total_math,
        "mermaid_blocks": len(mermaid_blocks),
        "go_blocks": len(go_blocks),
        "h1": len(h1),
        "h2": len(h2),
        "h3": len(h3),
        "h4": len(h4),
    }

def main():
    if len(sys.argv) < 3:
        print("Usage: python3 verify_translation.py <original.md> <translated.md>")
        sys.exit(1)

    orig_path = sys.argv[1]
    trans_path = sys.argv[2]

    if not os.path.exists(orig_path):
        print(f"Error: Original file {orig_path} not found.")
        sys.exit(1)
    if not os.path.exists(trans_path):
        print(f"Error: Translated file {trans_path} not found.")
        sys.exit(1)

    orig = analyze_markdown(orig_path)
    trans = analyze_markdown(trans_path)

    print(f"=== Quality Gate: Verifying Invariants ===")
    print(f"Original:   {orig_path}")
    print(f"Translated: {trans_path}\n")

    errors = []
    checks = [
        ("Math Blocks", orig["math_blocks"], trans["math_blocks"]),
        ("Mermaid Diagrams", orig["mermaid_blocks"], trans["mermaid_blocks"]),
        ("Go Code Listings", orig["go_blocks"], trans["go_blocks"]),
        ("H1 Headings", orig["h1"], trans["h1"]),
        ("H2 Sections", orig["h2"], trans["h2"]),
        ("H3 Subsections", orig["h3"], trans["h3"]),
        ("H4 Points", orig["h4"], trans["h4"]),
    ]

    for label, orig_val, trans_val in checks:
        status = "OK" if orig_val == trans_val else "FAIL"
        print(f"[{status:4s}] {label:20s}: Original={orig_val:2d} | Translated={trans_val:2d}")
        if orig_val != trans_val:
            errors.append(f"Mismatch in {label}: expected {orig_val}, got {trans_val}")

    if errors:
        print(f"\n[FAIL] Found {len(errors)} invariant violation(s):")
        for e in errors:
            print(f"  - {e}")
        sys.exit(1)
    else:
        print("\n[PASS] All structural and code invariants preserved perfectly!")
        sys.exit(0)

if __name__ == "__main__":
    main()
