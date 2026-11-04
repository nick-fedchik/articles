#!/usr/bin/env python3
"""Перевірка українських документів усієї книги або одного тому."""

from __future__ import annotations

import argparse
import ast
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile
from collections import Counter
from dataclasses import dataclass
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent
BOOK = ROOT / "Expert_Systems"
TEMP_ROOT = ROOT / "review" / "tmp"
SCRIPT_ROOT = ROOT / "scripts"
sys.path.insert(0, str(SCRIPT_ROOT))
import inspect_structural_hierarchy  # noqa: E402


@dataclass(frozen=True)
class Finding:
    path: Path
    line: int
    code: str
    message: str
    warning: bool = False


@dataclass
class Details:
    line: int
    caption: str | None
    opened: bool
    summary: str | None = None


@dataclass(frozen=True)
class CodeBlock:
    path: Path
    line: int
    language: str
    source: str
    line_count: int
    details: Details | None
    skip: bool
    caption: str | None = None


FENCE = re.compile(r"^ {0,3}(`{3,}|~{3,})(.*)$")
HEADING = re.compile(r"^(#{1,6})\s+(.+?)\s*$")
CHAPTER = re.compile(r"^ch(\d+)-.*\.md$")
CAPTION = re.compile(
    r"\*\*(Формула|Рівняння|Таблиця|Схема|Діаграма|Графік|Рисунок|Лістинг)"
    r" (\d+)\.(\d+)(?:\.\s*|\s+[—–-]\s+).*?\*\*"
)
LISTING_CAPTION = re.compile(
    r"^\*\*Лістинг (\d+)\.(\d+)(?:\.\s*|\s+[—–-]\s+).*\(([^()]*)\)\.?\*\*"
)
LISTING_SUMMARY = re.compile(
    r"^Показати лістинг (\d+)\.(\d+) \((\d+) "
    r"(рядок|рядки|рядків), ([^)]+)\)$"
)
HTML_TAG = re.compile(r"<details\b[^>]*>|</details>|<summary>(.*?)</summary>", re.I)
ID = re.compile(r'<[a-z][\w:-]*\b[^>]*\bid=["\']([^"\']+)["\'][^>]*>', re.I)
INLINE_CODE = re.compile(r"(?<!`)(`+)(?!`)(.*?)\1(?!`)")
ANCHOR_LINE = re.compile(r'<a\b[^>]*(?:id|name)=["\'][^"\']+["\'][^>]*></a>', re.I)
TABLE_SEPARATOR = re.compile(r"^\s*\|?\s*:?-{2,}:?\s*(?:\|\s*:?-{2,}:?\s*)*\|?\s*$")
VISIBLE_BLOCKS = {"text", "plaintext", "output", "console", "mermaid", "math"}
LANGUAGE_NAMES = {
    "python": "python",
    "python3": "python",
    "py": "python",
    "prolog": "prolog",
    "go": "go",
    "golang": "go",
    "bash": "bash",
    "sh": "shell",
    "clips": "clips",
    "minizinc": "minizinc",
    "lisp": "lisp",
    "turtle": "turtle",
    "json": "json",
    "smt2": "smt",
}
META_HEADINGS = {
    "Анотація", "Зміст глави", "Висновки", "Глосарій",
    "Абревіатури", "Скорочення", "Скорочення та терміни",
    "Джерела", "Література", "Посилання",
}


def display_path(path: Path) -> str:
    try:
        return str(path.relative_to(ROOT))
    except ValueError:
        return str(path)


def discover_documents(volume: int | None) -> tuple[Path, list[Path]]:
    scope = BOOK / f"Volume-{volume}" if volume is not None else BOOK
    if BOOK.is_symlink() or scope.is_symlink():
        raise ValueError(f"Символічний каталог книги не підтримується: {scope}")
    if not scope.is_dir():
        raise FileNotFoundError(f"Каталог не знайдено: {scope}")
    files = []

    def raise_walk_error(error: OSError) -> None:
        raise error

    for directory, subdirectories, names in os.walk(scope, topdown=True, onerror=raise_walk_error):
        subdirectories[:] = sorted(name for name in subdirectories if name not in {"translations", ".tmp"})
        for name in subdirectories:
            if (Path(directory) / name).is_symlink():
                raise ValueError(f"Символічний каталог не перевіряється: {Path(directory) / name}")
        for name in names:
            if not name.endswith(".md"):
                continue
            path = Path(directory) / name
            if path.is_symlink():
                raise ValueError(f"Символічний документ не перевіряється: {path}")
            if path.is_file():
                files.append(path)
    files.sort()
    if not files:
        raise ValueError(f"У каталозі немає документів Markdown: {scope}")
    return scope, files


def plural_for(number: int) -> str:
    if number % 10 == 1 and number % 100 != 11:
        return "рядок"
    if 2 <= number % 10 <= 4 and not 12 <= number % 100 <= 14:
        return "рядки"
    return "рядків"


def following_caption(lines: list[str], start: int, kinds: set[str]) -> bool:
    for line in lines[start:]:
        stripped = line.strip()
        if not stripped or ANCHOR_LINE.fullmatch(stripped):
            continue
        caption = CAPTION.search(stripped)
        return bool(caption and caption.group(1) in kinds)
    return False


def following_math_caption(lines: list[str], start: int) -> bool:
    for position in range(start, len(lines)):
        line = lines[position]
        stripped = line.strip()
        if not stripped or ANCHOR_LINE.fullmatch(stripped):
            continue
        if stripped == "$$" or HEADING.match(stripped) or FENCE.match(line) or stripped.startswith("|"):
            return False
        caption = CAPTION.search(stripped)
        if caption:
            if caption.group(1) not in {"Формула", "Рівняння"}:
                return False
            for next_line in lines[position + 1:]:
                next_stripped = next_line.strip()
                if not next_stripped or ANCHOR_LINE.fullmatch(next_stripped):
                    continue
                next_fence = FENCE.match(next_line)
                return not (
                    next_stripped == "$$"
                    or next_fence and next_fence.group(2).strip().lower() == "math"
                )
            return True
    return False


def scan_document(path: Path, text: str) -> tuple[list[Finding], list[CodeBlock]]:
    findings: list[Finding] = []
    blocks: list[CodeBlock] = []
    headings: list[tuple[int, str, int]] = []
    captions: list[tuple[int, int, int]] = []
    anchors: dict[str, int] = {}
    details: list[Details] = []
    pending_caption: str | None = None
    fence: tuple[str, int, str, Details | None, bool, str | None] | None = None
    body: list[str] = []
    previous_nonblank = ""
    in_math = False
    math_start = 0
    math_has_caption = False
    table_group_has_caption = False
    chapter_match = CHAPTER.match(path.name)
    chapter = int(chapter_match.group(1)) if chapter_match else None
    lines = text.splitlines()

    for index, line in enumerate(lines, 1):
        if fence:
            marker, start, language, context, skip, caption_text = fence
            close = re.fullmatch(rf" {{0,3}}{re.escape(marker[0])}{{{len(marker)},}}\s*", line)
            if close:
                blocks.append(CodeBlock(path, start, language, "\n".join(body), len(body),
                                        context, skip, caption_text))
                if language == "math" and chapter is not None and not math_has_caption:
                    if not following_math_caption(lines, index):
                        findings.append(Finding(path, start, "MATH_CAPTION",
                                                "Для блокової формули потрібен нумерований підпис"))
                fence, body = None, []
            else:
                body.append(line)
            continue

        opening = FENCE.match(line)
        if opening:
            marker, info = opening.groups()
            language = info.strip().split(maxsplit=1)[0].lower() if info.strip() else ""
            skip = bool(re.search(r"<!--\s*skip-validate\s*-->", previous_nonblank))
            if language == "math" and chapter is not None:
                caption = CAPTION.search(pending_caption or previous_nonblank)
                math_has_caption = bool(caption and caption.group(1) in {"Формула", "Рівняння"})
            fence = (marker, index, language, details[-1] if details else None, skip, pending_caption)
            pending_caption = None
            continue

        if line.strip() == "$$":
            if not in_math:
                math_start = index
                caption = CAPTION.search(pending_caption or previous_nonblank)
                math_has_caption = bool(caption and caption.group(1) in {"Формула", "Рівняння"})
            elif chapter is not None and not math_has_caption:
                if not following_math_caption(lines, index):
                    findings.append(Finding(path, math_start, "MATH_CAPTION",
                                            "Для блокової формули потрібен нумерований підпис"))
            in_math = not in_math
            if not in_math:
                pending_caption = None
            previous_nonblank = "$$"
            continue
        if in_math:
            continue
        if (chapter is not None and "|" in line and index < len(lines) and "|" in lines[index]
                and TABLE_SEPARATOR.match(lines[index])):
            caption = CAPTION.search(pending_caption or previous_nonblank)
            table_group_has_caption = bool(
                caption and caption.group(1) == "Таблиця"
                or table_group_has_caption and "|" in previous_nonblank
            )
            after = index + 1
            while after < len(lines) and "|" in lines[after]:
                after += 1
            table_group_has_caption |= following_caption(lines, after, {"Таблиця"})
            if not table_group_has_caption:
                findings.append(Finding(path, index, "TABLE_CAPTION",
                                        "Для таблиці потрібен нумерований підпис"))

        heading = HEADING.match(line)
        if heading:
            headings.append((len(heading.group(1)), heading.group(2), index))
        markup = INLINE_CODE.sub(lambda match: " " * len(match.group(0)), line)
        for match in ID.finditer(markup):
            anchor = match.group(1)
            if anchor in anchors:
                findings.append(Finding(path, index, "ANCHOR_DUPLICATE",
                                        f"Якір {anchor!r} уже визначено в рядку {anchors[anchor]}"))
            else:
                anchors[anchor] = index
        caption = CAPTION.search(markup)
        if caption:
            pending_caption = line
            captions.append((int(caption.group(2)), int(caption.group(3)), index))

        tags = list(HTML_TAG.finditer(markup))
        for tag in tags:
            value = tag.group(0)
            if value.lower().startswith("<details"):
                opened = bool(re.search(r"\sopen(?:\s|=|>)", value, re.I))
                details.append(Details(index, pending_caption, opened))
                pending_caption = None
            elif value.lower() == "</details>":
                if details:
                    details.pop()
                else:
                    findings.append(Finding(path, index, "DETAILS_UNMATCHED",
                                            "Закривальний </details> не має відкривального"))
            elif details:
                details[-1].summary = tag.group(1).strip()

        stripped = line.strip()
        if stripped and not caption and not tags and not ANCHOR_LINE.fullmatch(stripped):
            pending_caption = None
        if stripped:
            previous_nonblank = stripped

    if fence:
        findings.append(Finding(path, fence[1], "FENCE_UNCLOSED", "Кодовий блок не закрито"))
    for context in details:
        findings.append(Finding(path, context.line, "DETAILS_UNCLOSED", "Блок <details> не закрито"))

    if chapter is None:
        return findings, blocks

    numbered_sections = 0
    subsection = 0
    item = 0
    active_section: int | None = None
    active_subsection: int | None = None
    expected_toc: list[tuple[str, int]] = []
    title_count = 0
    for depth, title, line in headings:
        if depth == 1:
            title_count += 1
            if not re.match(rf"^Глава {chapter}\.\s+\S", title):
                findings.append(Finding(path, line, "HEADING_TITLE", f"Очікується «Глава {chapter}. …»"))
        elif depth == 2:
            if title not in {"Анотація", "Зміст глави"}:
                expected_toc.append((title, line))
            match = re.match(r"^(\d+)\.\s+\S", title)
            if match:
                actual = int(match.group(1))
                if actual != numbered_sections + 1:
                    findings.append(Finding(path, line, "HEADING_SECTION",
                                            f"Очікується розділ {numbered_sections + 1}, отримано {actual}"))
                numbered_sections, active_section = actual, actual
            elif title in META_HEADINGS or title.startswith("Підсумки Тому "):
                active_section = None
            else:
                findings.append(Finding(path, line, "HEADING_SECTION", "Розділ не має номера"))
                active_section = None
            subsection, item, active_subsection = 0, 0, None
        elif depth == 3:
            if active_section is not None:
                expected_toc.append((title, line))
                match = re.match(r"^(\d+)\.(\d+)\.\s+\S", title)
                expected = (active_section, subsection + 1)
                if not match or tuple(map(int, match.groups())) != expected:
                    findings.append(Finding(path, line, "HEADING_SUBSECTION",
                                            f"Очікується підрозділ {expected[0]}.{expected[1]}"))
                if match:
                    subsection = int(match.group(2))
                    active_subsection = subsection
                item = 0
        elif depth == 4 and active_section is not None:
            match = re.match(r"^(\d+)\.(\d+)\.(\d+)\.\s+\S", title)
            expected = (active_section, active_subsection, item + 1)
            if not match or tuple(map(int, match.groups())) != expected:
                findings.append(Finding(path, line, "HEADING_ITEM",
                                        f"Очікується пункт {expected[0]}.{expected[1]}.{expected[2]}"))
            if match:
                item = int(match.group(3))
        elif depth > 4:
            findings.append(Finding(path, line, "HEADING_DEPTH", "Дозволено не більше чотирьох рівнів"))
    if title_count != 1:
        findings.append(Finding(path, 1, "HEADING_TITLE", f"Очікується один заголовок глави, знайдено {title_count}"))

    start = text.find("<!-- chapter-toc:start -->")
    end = text.find("<!-- chapter-toc:end -->")
    if start < 0 or end < start:
        findings.append(Finding(path, 1, "TOC_MISSING", "Не знайдено меж змісту глави"))
    else:
        titles = re.findall(r"(?m)^\s*[-*]\s+\[([^\]]+)\]\([^)]+\)", text[start:end])
        for title, line in expected_toc:
            if title not in titles:
                findings.append(Finding(path, line, "TOC_ENTRY", f"У змісті немає заголовка «{title}»"))
        for title in titles:
            if title not in {expected for expected, _ in expected_toc}:
                findings.append(Finding(path, text[:start].count("\n") + 1, "TOC_STALE",
                                        f"Зайвий або застарілий пункт змісту «{title}»"))

    for expected, (caption_chapter, number, line) in enumerate(captions, 1):
        if caption_chapter != chapter or number != expected:
            findings.append(Finding(path, line, "OBJECT_NUMBER",
                                    f"Очікується об'єкт {chapter}.{expected}, отримано {caption_chapter}.{number}"))
    if not captions:
        findings.append(Finding(path, 1, "OBJECT_MISSING", "У главі немає пронумерованих навчальних об'єктів"))

    diagrams = sum(block.language == "mermaid" for block in blocks)
    if diagrams < 2:
        findings.append(Finding(path, 1, "DIAGRAM_COUNT",
                                "Потрібні принаймні дві діаграми Mermaid; їхнє призначення перевіряється окремо"))
    for block in blocks:
        if block.language == "mermaid":
            caption = CAPTION.search(block.caption or "")
            if not caption or caption.group(1) not in {"Схема", "Діаграма", "Графік", "Рисунок"}:
                findings.append(Finding(path, block.line, "DIAGRAM_CAPTION",
                                        "Перед діаграмою потрібен нумерований підпис"))
            for label in ("accTitle", "accDescr"):
                if not re.search(rf"(?m)^\s*{label}:\s*\S", block.source):
                    findings.append(Finding(path, block.line, "DIAGRAM_ACCESSIBILITY",
                                            f"У діаграмі немає {label}"))
        if block.language in VISIBLE_BLOCKS:
            continue
        if not block.language:
            findings.append(Finding(path, block.line, "LISTING_LANGUAGE", "Для лістингу не задано мову"))
        context = block.details
        if context is None:
            findings.append(Finding(path, block.line, "LISTING_COLLAPSE", "Лістинг не згорнуто в <details>"))
            continue
        if context.opened:
            findings.append(Finding(path, context.line, "LISTING_OPEN", "Лістинг має бути згорнутим"))
        summary = LISTING_SUMMARY.match(context.summary or "")
        caption = LISTING_CAPTION.match(context.caption or "")
        if not summary:
            findings.append(Finding(path, context.line, "LISTING_SUMMARY",
                                    "Підпис має містити номер, кількість рядків, мову та версію"))
        if not caption:
            findings.append(Finding(path, context.line, "LISTING_CAPTION",
                                    "Перед лістингом потрібен нумерований підпис із мовою та версією"))
        if summary:
            caption_chapter, number, count, word, version = summary.groups()
            actual_count = block.line_count
            if int(count) != actual_count or word != plural_for(actual_count):
                findings.append(Finding(path, block.line, "LISTING_LINES",
                                        f"У підписі {count} {word}, фактично {actual_count} {plural_for(actual_count)}"))
            name = LANGUAGE_NAMES.get(block.language, block.language).casefold()
            if not name or name not in version.casefold() or not re.search(r"\d", version):
                findings.append(Finding(path, context.line, "LISTING_VERSION",
                                        f"У підписі не вказано мову {block.language!r} та її версію"))
            if int(caption_chapter) != chapter:
                findings.append(Finding(path, context.line, "LISTING_NUMBER", "Номер глави в підписі не збігається"))
            if caption and (caption.group(1), caption.group(2), caption.group(3)) != (
                caption_chapter, number, version
            ):
                findings.append(Finding(path, context.line, "LISTING_NUMBER",
                                        "Номер або версія в підписі й summary не збігаються"))
        if block.line_count > 30:
            findings.append(Finding(path, block.line, "LISTING_LONG",
                                    "Лістинг довший за рекомендовані 30 рядків", warning=True))
    return findings, blocks


PROLOG_READER = (
    "catch((repeat,read_term(user_input,T,[syntax_errors(error)]),"
    "(T == end_of_file -> ! ; fail)),E,(print_message(error,E),halt(1)))"
)


def go_candidates(source: str) -> list[str]:
    source = re.sub(r"\{\s*\.\.\.\s*\}", "{}", source)
    package = "package example\n"
    imports = re.match(
        r'(?s)^\s*((?:import\s+(?:\([^)]*\)|(?:(?:[._]|\w+)\s+)?"[^"\n]+")\s*)+)(.*)$',
        source,
    )
    declarations, statements = imports.groups() if imports else ("", source)
    return [
        source,
        package + source,
        package + "func example() {\n" + source + "\n}\n",
        package + "var example = (\n" + source + "\n)\n",
        package + "func example() {\nfor {\n" + source + "\n}\n}\n",
        package + declarations + "\nfunc example() {\n" + statements + "\n}\n",
    ]


def run_syntax(blocks: list[CodeBlock]) -> tuple[list[Finding], Counter[str], Counter[str]]:
    findings: list[Finding] = []
    checked: Counter[str] = Counter()
    skipped: Counter[str] = Counter()
    missing: set[str] = set()
    for block in blocks:
        language = block.language
        if language in VISIBLE_BLOCKS or not language:
            continue
        if block.skip or "skip-validate" in block.source:
            skipped[f"{language}: директива"] += 1
            continue
        if language in {"python", "python3", "py"}:
            try:
                ast.parse(block.source, filename=display_path(block.path))
            except SyntaxError as error:
                findings.append(Finding(block.path, block.line + (error.lineno or 1),
                                        "PYTHON_SYNTAX", error.msg))
            else:
                checked["Python"] += 1
            continue
        if language == "json":
            try:
                json.loads(block.source)
            except json.JSONDecodeError as error:
                findings.append(Finding(block.path, block.line + error.lineno,
                                        "JSON_SYNTAX", error.msg))
            else:
                checked["JSON"] += 1
            continue
        if language == "prolog":
            tool, candidates = "swipl", [block.source]
            command = [tool, "-q", "-f", "none", "-g", PROLOG_READER, "-t", "halt"]
        elif language in {"go", "golang"}:
            tool, candidates = "gofmt", go_candidates(block.source)
            command = [tool, "-e"]
        elif language in {"bash", "sh"}:
            tool, candidates = language, [block.source]
            command = [tool, "-n"]
        else:
            skipped[f"{language}: немає синтаксичного аналізатора"] += 1
            continue
        if not shutil.which(tool):
            if tool not in missing:
                findings.append(Finding(block.path, block.line, "TOOL_MISSING",
                                        f"Не знайдено {tool}; синтаксис {language} не перевірено"))
                missing.add(tool)
            skipped[f"{language}: немає {tool}"] += 1
            continue
        errors = []
        for candidate in candidates:
            try:
                result = subprocess.run(command, input=candidate, text=True, capture_output=True, timeout=30)
            except (OSError, subprocess.TimeoutExpired) as error:
                findings.append(Finding(block.path, block.line, "TOOL_ERROR", f"{tool}: {error}"))
                break
            if result.returncode == 0:
                checked[language] += 1
                break
            errors.append(result.stderr.strip())
        else:
            message = next((error for error in errors if error), "синтаксичний аналіз завершився помилкою")
            findings.append(Finding(block.path, block.line, "CODE_SYNTAX", f"{language}: {message}"))
    return findings, checked, skipped


def print_findings(title: str, findings: list[Finding]) -> bool:
    errors = sum(not finding.warning for finding in findings)
    warnings = len(findings) - errors
    print(f"\n{title}: {errors} помилок, {warnings} попереджень")
    for finding in findings:
        level = "ПОПЕРЕДЖЕННЯ" if finding.warning else "ПОМИЛКА"
        print(f"{display_path(finding.path)}:{finding.line}: {level} [{finding.code}] {finding.message}")
    return errors == 0


def run_inspector(name: str, arguments: list[str]) -> bool:
    print(f"\n--- {name} ---", flush=True)
    try:
        result = subprocess.run(arguments, cwd=ROOT, text=True, capture_output=True, timeout=300)
    except (OSError, subprocess.TimeoutExpired) as error:
        print(f"ПОМИЛКА запуску {name}: {error}")
        return False
    if result.stdout:
        print(result.stdout, end="" if result.stdout.endswith("\n") else "\n")
    if result.stderr:
        print(result.stderr, file=sys.stderr, end="" if result.stderr.endswith("\n") else "\n")
    if result.returncode:
        print(f"Перевірка завершилася з кодом {result.returncode}")
    reported_error = bool(re.search(r"(?m)^\s*\[ERROR\]", result.stderr))
    return result.returncode == 0 and not reported_error


def run_render(files: list[Path], report: Path) -> bool:
    print("\n--- MathJax, Mermaid і локальні посилання ---", flush=True)
    command = [
        "node", str(SCRIPT_ROOT / "inspect_render_syntax.mjs"),
        "--workspace-only", f"--report={report}", *(str(path) for path in files),
    ]
    try:
        result = subprocess.run(command, cwd=ROOT, text=True, capture_output=True, timeout=600)
    except (OSError, subprocess.TimeoutExpired) as error:
        print(f"ПОМИЛКА запуску перевірки рендерингу: {error}")
        return False
    if not report.is_file():
        print(result.stdout, end="")
        print(result.stderr, file=sys.stderr, end="")
        print("Звіт рендерингу не створено; перевірте Node.js і залежності з scripts/render-validation.")
        return False
    try:
        data = json.loads(report.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as error:
        print(f"ПОМИЛКА читання звіту рендерингу: {error}")
        return False
    for kind in ("math", "mermaid", "links"):
        counts = data["summary"][kind]
        print(f"{kind}: {counts['passed']}/{counts['total']} пройшли, "
              f"{counts['failed']} помилок, {counts.get('skipped', 0)} пропущено")
    failed = [item for key in ("objects", "links") for item in data[key] if item["status"] == "FAIL"]
    for item in failed:
        print(f"{item['file']}:{item['line']}: {item.get('type', 'link')}: {item.get('error', '')}")
    restrictions = data["summary"]["github_observed_macro_restrictions"]
    if restrictions["failed"]:
        for item in data["objects"]:
            if item["type"] == "math" and re.search(r"\\operatorname(?![A-Za-z])", item["source"]):
                print(f"{item['file']}:{item['line']}: GitHub може відхилити \\operatorname")
    if result.returncode and not failed and not restrictions["failed"]:
        print(result.stdout, end="")
        print(result.stderr, file=sys.stderr, end="")
    print("Зовнішні URL не перевіряються; посилання на заморожені переклади пропускаються.")
    return result.returncode == 0 and not failed and not restrictions["failed"]


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--volume", type=int, choices=(1, 2, 3), help="перевірити лише вказаний том")
    args = parser.parse_args(argv)
    try:
        scope, files = discover_documents(args.volume)
        documents = [(path, path.read_text(encoding="utf-8")) for path in files]
    except (OSError, UnicodeError, ValueError) as error:
        print(f"ПОМИЛКА пошуку або читання документів: {error}", file=sys.stderr)
        return 2
    print(f"Область: {display_path(scope)}; документів Markdown: {len(files)}")
    print("Переклади й каталоги .tmp не перевіряються. Лістинги аналізуються без виконання коду.")
    results: list[bool] = []
    findings: list[Finding] = []
    blocks: list[CodeBlock] = []
    for path, text in documents:
        document_findings, document_blocks = scan_document(path, text)
        findings.extend(document_findings)
        blocks.extend(document_blocks)
        for line, code, message in inspect_structural_hierarchy.inspect_file(path):
            findings.append(Finding(path, line, code, message))
    results.append(print_findings("Структура, зміст, нумерація й лістинги", findings))
    syntax_findings, checked, skipped = run_syntax(blocks)
    results.append(print_findings("Синтаксис програмних лістингів", syntax_findings))
    print(f"Перевірено блоків: {dict(checked)}; пропущено: {dict(skipped)}")
    print("Версії в підписах перевіряються як текст; виконання прикладів не підтверджується.")
    results.append(run_inspector(
        "Математичне форматування", [sys.executable, str(SCRIPT_ROOT / "inspect_math_formatting.py"),
                                    *(str(path) for path in files)]
    ))
    results.append(run_inspector(
        "Термінологія", [sys.executable, str(SCRIPT_ROOT / "inspect_buzzwords.py"),
                         *(str(path) for path in files)]
    ))
    results.append(run_inspector(
        "Витік внутрішніх артефактів",
        [sys.executable, str(SCRIPT_ROOT / "inspect_project_leakage.py"), *(str(path) for path in files)]
    ))
    try:
        TEMP_ROOT.mkdir(parents=True, exist_ok=True)
        with tempfile.TemporaryDirectory(prefix="book-check-", dir=TEMP_ROOT) as directory:
            results.append(run_render(files, Path(directory) / "render.json"))
    except OSError as error:
        print(f"ПОМИЛКА створення тимчасового каталогу в {TEMP_ROOT}: {error}", file=sys.stderr)
        results.append(False)
    print(f"\nПідсумок: {'УСПІХ' if all(results) else 'Є ПОМИЛКИ'} "
          f"({sum(results)}/{len(results)} груп перевірок пройдено).")
    return 0 if all(results) else 1


if __name__ == "__main__":
    sys.exit(main())
