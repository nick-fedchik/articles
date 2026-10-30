#!/usr/bin/env python3
"""
GitHub Flavored Markdown (GFM) Mathematical Expressions & Advanced Formatting Validator.

Enforces official GitHub Docs rules:
- https://docs.github.com/ru/get-started/writing-on-github/working-with-advanced-formatting/writing-mathematical-expressions
- https://docs.github.com/ru/get-started/writing-on-github/working-with-advanced-formatting
and project-specific MathJax/GFM invariants from AGENTS.md.

Checks performed:
  [ERR-MATH-01] Code-wrapped math (`$...$` instead of `$`...`$` or `$...$`)
  [ERR-MATH-02] Indented or inline-mixed block math `$$` (not at column 0, or inside `> ` blockquotes)
  [ERR-MATH-03] Multi-line `$$` block where `$$` shares a line with LaTeX content instead of being on its own `$$` line
  [ERR-MATH-04] Unisolated `$$...$$` block (missing blank line before or after)
  [ERR-MATH-05] Space-padded inline `$ ... $` delimiters (space after opening `$` or before closing `$`)
  [ERR-MATH-06] Unescaped vertical pipe `|` inside math within a Markdown table row
  [ERR-MATH-07] Unbalanced inline `$` or block `$$` delimiters
  [ERR-MATH-08] Corrupted LaTeX escape sequences (control chars \x07, \x08, \x0b, \x0c, or split \r/\n/\t macros)
  [ERR-MATH-09] Unprotected inline math containing Markdown formatting characters ('_' or '*')
  [ERR-MATH-10] Unbalanced curly braces or broken subscript/superscript patterns (missing '_' or '^') in math
  [ERR-ADV-01]  Missing blank line after `</summary>` inside `<details>` block
  [ERR-ADV-02]  Text on the same line after GitHub Alert tag `> [!NOTE|TIP|IMPORTANT|WARNING|CAUTION]`

Usage:
  python3 scripts/inspect_math_formatting.py [--fix] [path ...]
"""

import argparse
import os
import re
import sys
from typing import List, Tuple


def find_markdown_files(paths: List[str]) -> List[str]:
    if not paths:
        paths = ["ExpertSystem"]
    result = []
    for p in paths:
        if os.path.isfile(p) and p.endswith(".md"):
            result.append(p)
        elif os.path.isdir(p):
            for root, dirs, files in os.walk(p):
                if "registry" in root.split(os.sep) or "translations" in root.split(os.sep):
                    continue
                for f in sorted(files):
                    if f.endswith(".md"):
                        result.append(os.path.join(root, f))
    return sorted(set(result))


def strip_code_and_valid_backtick_math(line: str) -> str:
    """Replace valid $`...`$, escaped \\$, and regular `...` code spans with placeholders."""
    s = re.sub(r"\$`[^`]+`\$", "MATH_BACKTICK_OK", line)

    def _repl_code(m: re.Match) -> str:
        inner = m.group(1)
        if inner.startswith("$") and not inner.startswith("$$") and inner.endswith("$") and len(inner) >= 2:
            return m.group(0)
        return "CODE_SPAN_OK"

    s = re.sub(r"`([^`]+)`", _repl_code, s)
    s = s.replace(r"\$", "ESCAPED_DOLLAR_OK")
    return s


def protect_valid_math(text: str) -> Tuple[str, List[str]]:
    tokens = []

    def _rep(m: re.Match) -> str:
        tokens.append(m.group(0))
        return f"__VALID_MATH_BACKTICK_{len(tokens)-1}__"

    return re.sub(r"\$`[^`]+`\$", _rep, text), tokens


def restore_valid_math(text: str, tokens: List[str]) -> str:
    for idx, t in enumerate(tokens):
        text = text.replace(f"__VALID_MATH_BACKTICK_{idx}__", t)
    return text


def repair_corrupted_latex_escapes(text: str) -> Tuple[str, bool]:
    """Repair Python string-escape corruptions of LaTeX commands (\frac, \approx, \rightarrow, \text, etc.)."""
    orig = text
    # Control character replacements
    text = text.replace("\x0crac", r"\frac")
    text = text.replace("\x07pprox", r"\approx")
    text = text.replace("\x08igl", r"\bigl")
    text = text.replace("\x08igr", r"\bigr")
    text = text.replace("\x08egin", r"\begin")
    text = text.replace("\x0barnothing", r"\varnothing")
    # Split \right + \rvert / \rceil across two newlines
    text = text.replace("\night\nvert", r"\right\rvert")
    text = text.replace("\night\nceil", r"\right\rceil")
    text = text.replace("\night\\}", r"\right\}")
    text = text.replace("\night)", r"\right)")
    text = text.replace("\nightarrow", r"\rightarrow")
    text = text.replace("\nvert", r"\rvert")
    text = text.replace("\nVert_2", r"\rVert_2")
    text = text.replace("\nVert", r"\rVert")
    text = text.replace("\nceil", r"\rceil")
    text = text.replace("\nangle", r"\rangle")
    text = text.replace("\neg H", r"\neg H")
    text = text.replace("B\\cap A\ne\\varnothing", r"B\cap A\ne\varnothing")
    text = text.replace("B\\cap A\n\\varnothing", r"B\cap A\ne\varnothing")
    # Tab-corrupted LaTeX macros
    text = text.replace("\text{", r"\text{")
    text = text.replace("\tau", r"\tau")
    text = re.sub(r"\to(?=\s|\[|\||\\)", r"\\to", text)
    return text, (text != orig)


def inspect_file(filepath: str, fix: bool = False) -> Tuple[List[Tuple[int, str, str]], bool]:
    with open(filepath, "r", encoding="utf-8") as f:
        raw_text = f.read()

    modified = False
    if fix:
        raw_text, esc_mod = repair_corrupted_latex_escapes(raw_text)
        if esc_mod:
            modified = True

    lines = raw_text.splitlines(keepends=True)
    errors: List[Tuple[int, str, str]] = []

    if fix:
        new_lines: List[str] = []
        in_fence = False
        in_math_block = False
        i = 0
        n = len(lines)
        while i < n:
            line = lines[i]
            stripped = line.strip()

            if stripped.startswith("```"):
                in_fence = not in_fence
                new_lines.append(line)
                i += 1
                continue

            if in_fence:
                new_lines.append(line)
                i += 1
                continue

            # Track multi-line $$ blocks
            if in_math_block:
                line_no_esc = line.replace(r"\$", "ESCAPED_DOLLAR_OK")
                stripped_no_esc = line_no_esc.strip()
                if stripped_no_esc.endswith("$$") and stripped_no_esc != "$$":
                    content_part = stripped[:-2].rstrip()
                    if content_part:
                        new_lines.append(content_part + "\n")
                    new_lines.append("$$\n")
                    in_math_block = False
                    modified = True
                    if i + 1 < n and lines[i + 1].strip() != "":
                        new_lines.append("\n")
                    i += 1
                    continue
                elif stripped_no_esc == "$$":
                    if not line.startswith("$$"):
                        line = "$$\n"
                        modified = True
                    new_lines.append(line)
                    in_math_block = False
                    if i + 1 < n and lines[i + 1].strip() != "":
                        new_lines.append("\n")
                        modified = True
                    i += 1
                    continue
                else:
                    new_lines.append(line)
                    i += 1
                    continue

            line_no_esc = strip_code_and_valid_backtick_math(line)
            if "$$" in line_no_esc:
                # Case A: blockquote `> $$...$$` -> convert to `> $...$`
                m_bq = re.match(r"^(\s*>\s*)\$\$(.+)\$\$\s*$", line)
                if m_bq:
                    prefix, expr = m_bq.group(1), m_bq.group(2).strip()
                    line = f"{prefix}${expr}$\n"
                    modified = True
                else:
                    # Case B: indented single-line `   $$...$$`
                    m_single = re.match(r"^(\s+)\$\$(.+)\$\$\s*$", line)
                    if m_single:
                        expr = m_single.group(2).strip()
                        if new_lines and new_lines[-1].strip() != "":
                            new_lines.append("\n")
                        new_lines.append(f"$${expr}$$\n")
                        if i + 1 < n and lines[i + 1].strip() != "":
                            new_lines.append("\n")
                        modified = True
                        i += 1
                        continue

                    # Case C: unisolated single-line `$$...$$` at col 0
                    m_col0_single = re.match(r"^\$\$(.+)\$\$\s*$", line)
                    if m_col0_single:
                        if new_lines and new_lines[-1].strip() != "" and not new_lines[-1].rstrip().endswith("\\"):
                            new_lines.append("\n")
                            modified = True
                        new_lines.append(line)
                        if i + 1 < n and lines[i + 1].strip() != "":
                            new_lines.append("\n")
                            modified = True
                        i += 1
                        continue

                    # Case D: multi-line opening `$$...` where LaTeX starts on same line as `$$`
                    m_open_inline = re.match(r"^\s*\$\$(.+)$", line)
                    if m_open_inline and not stripped.endswith("$$"):
                        expr_start = m_open_inline.group(1).rstrip()
                        if new_lines and new_lines[-1].strip() != "" and not new_lines[-1].rstrip().endswith("\\"):
                            new_lines.append("\n")
                        new_lines.append("$$\n")
                        new_lines.append(expr_start + "\n")
                        in_math_block = True
                        modified = True
                        i += 1
                        continue

                    # Case E: normal multi-line opening `$$`
                    if stripped == "$$":
                        if not line.startswith("$$"):
                            line = "$$\n"
                            modified = True
                        if new_lines and new_lines[-1].strip() != "" and not new_lines[-1].rstrip().endswith("\\"):
                            new_lines.append("\n")
                            modified = True
                        new_lines.append(line)
                        in_math_block = True
                        i += 1
                        continue

            # Fix [ERR-MATH-01]: `$...$` -> $`...`$ (when complex) or $...$
            prot_line, prot_tokens = protect_valid_math(line)
            def _fix_code_wrapped_math(m: re.Match) -> str:
                nonlocal modified
                expr = m.group(1).strip()
                modified = True
                if any(ch in expr for ch in ("_", "|", "*", "{", "}")):
                    return f"$`{expr}`$"
                return f"${expr}$"

            prot_line = re.sub(r"(?<!\$)`\$([^`\$]+)\$`(?!\$)", _fix_code_wrapped_math, prot_line)
            line = restore_valid_math(prot_line, prot_tokens)

            # Fix [ERR-MATH-09]: Protect inline $expr$ containing Markdown special chars (_, *) with $`expr`$
            prot_line, prot_tokens = protect_valid_math(line)
            def _fix_unprotected_markdown_math(m: re.Match) -> str:
                nonlocal modified
                expr = m.group(1).strip()
                if any(ch in expr for ch in ("_", "*", "|", "{", "}")):
                    modified = True
                    return f"$`{expr}`$"
                return m.group(0)

            prot_line = re.sub(r"(?<![\$`\\])\$([^\$\n]+)\$(?![\$`])", _fix_unprotected_markdown_math, prot_line)
            line = restore_valid_math(prot_line, prot_tokens)

            # Fix [ERR-ADV-01]: missing blank line after </summary>
            if stripped.endswith("</summary>") and i + 1 < n and lines[i + 1].strip() != "":
                new_lines.append(line)
                new_lines.append("\n")
                modified = True
                i += 1
                continue

            # Fix [ERR-ADV-02]: GitHub alert tag with trailing text on the same line
            m_alert = re.match(r"^(\s*>\s*)\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\]\s+(.+)$", line)
            if m_alert:
                prefix, tag, rest = m_alert.group(1), m_alert.group(2), m_alert.group(3).strip()
                if not rest.startswith("**"):
                    rest = f"**{rest}**"
                new_lines.append(f"{prefix}[!{tag}]\n")
                new_lines.append(f"{prefix}{rest}\n")
                modified = True
                i += 1
                continue

            new_lines.append(line)
            i += 1

        if modified:
            lines = new_lines
            with open(filepath, "w", encoding="utf-8") as f:
                f.writelines(lines)

    # Validation pass
    full_text = "".join(lines)
    # Check multi-line split \r/\n LaTeX macros across lines
    for m in re.finditer(r"\n(ightarrow|ightleftharpoons|ight|vert|Vert|ceil|floor|angle|eg\b|eq\b)", full_text):
        line_no = full_text[: m.start()].count("\n") + 2
        errors.append((
            line_no,
            "ERR-MATH-08",
            f"Corrupted newline-split LaTeX macro '\\r{m.group(1)}' or '\\n{m.group(1)}'"
        ))

    in_fence = False
    in_math_fence = False
    in_math_block = False
    math_block_open_line = 0

    for idx, raw_line in enumerate(lines, 1):
        line = raw_line.rstrip("\r\n")
        stripped = line.strip()

        # Check control characters and tab-corrupted LaTeX macros (except inside Go code blocks where \x02\x03 test fixture exists)
        if not in_fence or in_math_fence:
            if re.search(r"[\x00-\x08\x0b\x0c\x0e-\x1f]", line):
                errors.append((
                    idx,
                    "ERR-MATH-08",
                    f"Corrupted ASCII control character in LaTeX/text: {repr(line[:60])}"
                ))
            if re.search(r"\t(ext\{|extrm\{|extbf\{|extit\{|au\b|o\b|imes\b|heta\b|op\b|riangle|ilde|anh\b)", line):
                errors.append((
                    idx,
                    "ERR-MATH-08",
                    f"Corrupted tab-escaped LaTeX macro (\\text, \\tau, \\to, etc.): {repr(line[:60])}"
                ))

        if stripped.startswith("```"):
            if not in_fence:
                in_fence = True
                in_math_fence = stripped.startswith("```math")
            else:
                in_fence = False
                in_math_fence = False
            continue
        if in_fence:
            continue

        line_no_esc = strip_code_and_valid_backtick_math(line)
        stripped_no_esc = line_no_esc.strip()

        # Multi-line $$ tracking
        if in_math_block:
            if "$$" in stripped_no_esc:
                if stripped_no_esc != "$$":
                    errors.append((
                        idx,
                        "ERR-MATH-03",
                        f"Closing '$$' in multi-line block must be on its own line at column 0: {stripped[:60]}"
                    ))
                elif not line.startswith("$$"):
                    errors.append((
                        idx,
                        "ERR-MATH-02",
                        f"Closing '$$' is indented (must be at column 0): {repr(line[:20])}"
                    ))
                if idx < len(lines) and lines[idx].strip() != "":
                    errors.append((
                        idx,
                        "ERR-MATH-04",
                        f"Missing blank line after closing '$$': next line is '{lines[idx].strip()[:40]}'"
                    ))
                in_math_block = False
            continue

        # Check $$ outside multi-line block
        if "$$" in line_no_esc:
            if stripped_no_esc == "$$":
                if not line.startswith("$$"):
                    errors.append((
                        idx,
                        "ERR-MATH-02",
                        f"Opening '$$' is indented (must be at column 0): {repr(line[:20])}"
                    ))
                if idx > 1 and lines[idx - 2].strip() != "" and not lines[idx - 2].rstrip().endswith("\\"):
                    errors.append((
                        idx,
                        "ERR-MATH-04",
                        f"Missing blank line before opening '$$': prev line is '{lines[idx - 2].strip()[:40]}'"
                    ))
                in_math_block = True
                math_block_open_line = idx
                continue
            elif (
                stripped_no_esc.startswith("$$")
                and stripped_no_esc.endswith("$$")
                and len(stripped_no_esc) > 4
                and stripped_no_esc.count("$$") == 2
            ):
                if not line.startswith("$$"):
                    errors.append((
                        idx,
                        "ERR-MATH-02",
                        f"Single-line '$$...$$' is indented (must be at column 0): {stripped[:60]}"
                    ))
                if idx > 1 and lines[idx - 2].strip() != "" and not lines[idx - 2].rstrip().endswith("\\"):
                    errors.append((
                        idx,
                        "ERR-MATH-04",
                        f"Missing blank line before '$$...$$': prev line is '{lines[idx - 2].strip()[:40]}'"
                    ))
                if idx < len(lines) and lines[idx].strip() != "":
                    errors.append((
                        idx,
                        "ERR-MATH-04",
                        f"Missing blank line after '$$...$$': next line is '{lines[idx].strip()[:40]}'"
                    ))
                continue
            else:
                if stripped_no_esc.startswith("$$") and not stripped_no_esc.endswith("$$"):
                    errors.append((
                        idx,
                        "ERR-MATH-03",
                        f"Multi-line '$$' block has LaTeX on the opening '$$' line (put '$$' alone on its line): {stripped[:60]}"
                    ))
                    in_math_block = True
                    math_block_open_line = idx
                    continue
                else:
                    errors.append((
                        idx,
                        "ERR-MATH-02",
                        f"Block delimiter '$$' used inline or inside blockquote/list: {stripped[:70]}"
                    ))
                    continue

        # Check [ERR-MATH-01]: `$...$` (backtick outside dollar sign)
        prot_line, _ = protect_valid_math(line)
        for m in re.finditer(r"(?<!\$)`\$([^`\$]+)\$`(?!\$)", prot_line):
            errors.append((
                idx,
                "ERR-MATH-01",
                f"Code-wrapped math `{m.group(0)}` disables MathJax! Use GitHub `$`...`$` or `$...$`."
            ))

        # Check [ERR-MATH-05] & [ERR-MATH-06] & [ERR-MATH-07] on remaining $...$
        if "$" in line_no_esc and "$$" not in line_no_esc:
            dollar_positions = [m.start() for m in re.finditer(r"(?<!\\)\$", line_no_esc)]
            if len(dollar_positions) % 2 != 0:
                errors.append((
                    idx,
                    "ERR-MATH-07",
                    f"Unbalanced inline '$' delimiters ({len(dollar_positions)} found): {stripped[:70]}"
                ))
            else:
                is_table_row = stripped.startswith("|") and stripped.endswith("|")
                for k in range(0, len(dollar_positions), 2):
                    start_p = dollar_positions[k]
                    end_p = dollar_positions[k + 1]
                    expr = line_no_esc[start_p + 1 : end_p]
                    if not expr:
                        errors.append((idx, "ERR-MATH-05", "Empty inline math '$$'"))
                        continue
                    if expr[0] == " " or expr[-1] == " ":
                        errors.append((
                            idx,
                            "ERR-MATH-05",
                            f"Space-padded inline math '${expr}$' fails to render in GitHub MathJax"
                        ))
                    if is_table_row and "|" in expr and r"\|" not in expr:
                        errors.append((
                            idx,
                            "ERR-MATH-06",
                            f"Unescaped pipe '|' inside inline math '${expr}$' in table row splits table cell (use \\lvert, \\rvert, or \\mid)"
                        ))
                    if any(ch in expr for ch in ("_", "*")):
                        errors.append((
                            idx,
                            "ERR-MATH-09",
                            f"Unprotected inline math '${expr}$' contains Markdown conflicting chars ('_' or '*'). Use GitHub protected syntax $`{expr}`$ to prevent italics/bold corruption."
                        ))
                    if expr.count("{") != expr.count("}"):
                        errors.append((
                            idx,
                            "ERR-MATH-10",
                            f"Unbalanced curly braces in inline math '${expr}$' ({expr.count('{')} opening vs {expr.count('}')} closing)"
                        ))
                    if re.search(r"\\(mathrm|mathbf|mathit|text|textbf|textit|operatorname|mathcal|mathbb|boldsymbol|bm)\{[^}]+\}\s*\{", expr):
                        errors.append((
                            idx,
                            "ERR-MATH-10",
                            f"Potentially broken subscript/superscript (missing '_' or '^') in inline math: '${expr}$'"
                        ))

        # Check [ERR-ADV-01]: </summary> followed immediately by non-blank line
        if stripped.endswith("</summary>") and idx < len(lines) and lines[idx].strip() != "":
            errors.append((
                idx,
                "ERR-ADV-01",
                "Missing blank line after </summary> inside <details> (breaks Markdown/code rendering on GitHub)"
            ))

        # Check [ERR-ADV-02]: GitHub alert tag with trailing text on the same line
        m_alert = re.match(r"^\s*>\s*\[!(NOTE|TIP|IMPORTANT|WARNING|CAUTION)\](.+)$", line)
        if m_alert and m_alert.group(2).strip() != "":
            errors.append((
                idx,
                "ERR-ADV-02",
                f"GitHub Alert '> [!{m_alert.group(1)}]' has text on the same line: '{m_alert.group(2).strip()[:40]}'"
            ))

    if in_math_block:
        errors.append((
            math_block_open_line,
            "ERR-MATH-07",
            f"Unclosed multi-line '$$' block opened at line {math_block_open_line}"
        ))

    return errors, modified


def main() -> int:
    parser = argparse.ArgumentParser(description="Inspect and optionally fix GitHub math and advanced formatting.")
    parser.add_argument("--fix", action="store_true", help="Automatically fix common mechanical formatting errors")
    parser.add_argument("paths", nargs="*", help="Files or directories to inspect (default: ExpertSystem)")
    args = parser.parse_args()

    files = find_markdown_files(args.paths)
    print(f"=== Inspecting {len(files)} Markdown files for GitHub MathJax & Advanced Formatting ===")

    total_errors = 0
    fixed_files = 0
    for f in files:
        errs, mod = inspect_file(f, fix=args.fix)
        if mod:
            fixed_files += 1
        if errs:
            total_errors += len(errs)
            print(f"\n[FAIL] {f} ({len(errs)} issues):")
            for line_no, code, msg in errs:
                print(f"  L{line_no:4d} [{code}]: {msg}")

    if fixed_files:
        print(f"\n[FIXED] Automatically repaired formatting in {fixed_files} file(s).")

    if total_errors == 0:
        print("\n[OK] All mathematical expressions (inline & block) and advanced GFM formatting are 100% valid!")
        return 0
    else:
        print(f"\n[ERROR] Found {total_errors} formatting violation(s). Run with --fix or repair manually.")
        return 1


if __name__ == "__main__":
    sys.exit(main())
