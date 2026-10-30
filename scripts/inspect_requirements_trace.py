#!/usr/bin/env python3
"""
Requirements Traceability Inspector for the ExpertSystem monograph.

Checks the requirement tree in ExpertSystem/requirements (L1 stakeholder,
L2 system, L3 component) against the conventions of the tree:

  ERRORS (exit code 1):
    - requirement ID is duplicated or does not match its file (level, chapter);
    - a mandatory field is missing (Вимога, Джерело, Батьківська вимога,
      Метод верифікації, Статус);
    - L1 has a parent, or an L2/L3 parent does not exist one level up;
    - the «Канонічна вимога» field points to a missing requirement, to another
      level, to itself, or to a requirement that is itself a duplicate;
    - the source link points to a missing file;
    - the verification method does not start with Тест, Аналіз, Інспекція
      or Демонстрація.

  WARNINGS (exit code 0):
    - parent requirement belongs to another chapter (allowed, but listed);
    - consecutive header fields without a blank line between them;
    - banned words or an em dash in the requirement text;
    - with --check-sections: a cited «§N «Title»» does not match any heading
      of the cited chapter (prefix match after normalisation).

Usage:
  python scripts/inspect_requirements_trace.py [--check-sections] [--quiet]
"""

import re
import sys
from collections import defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent / "ExpertSystem"
REQ_ROOT = ROOT / "requirements"
LEVEL_DIRS = {"L1": "L1-stakeholder", "L2": "L2-system", "L3": "L3-component"}
PARENT_LEVEL = {"L2": "L1", "L3": "L2"}
MANDATORY = ["Вимога", "Джерело", "Батьківська вимога", "Метод верифікації", "Пріоритет", "Статус"]
METHODS = ("Тест", "Аналіз", "Інспекція", "Демонстрація")
BANNED = [
    (re.compile(r"екосистем|ecosystem", re.I), "banned word 'екосистема'"),
    (re.compile(r"\bконтур", re.I), "generated-text marker 'контур'"),
    (re.compile(r"\bключов(?!\w*\s+слов)", re.I), "generated-text marker 'ключовий'"),
    (re.compile(r"фундаментальн", re.I), "generated-text marker 'фундаментальний'"),
    (re.compile(r"\bунікальн", re.I), "generated-text marker 'унікальний'"),
    (re.compile(r"\bпотужн", re.I), "generated-text marker 'потужний'"),
    (re.compile(r"революційн", re.I), "generated-text marker 'революційний'"),
    (re.compile(r"Таким чином"), "generated-text marker 'Таким чином'"),
    (re.compile(r"\bчесн", re.I), "generated-text marker 'чесно'"),
    (re.compile("\u2014"), "em dash in prose"),
]

REQ_HEAD = re.compile(r"^## (REQ-(L[123])-(\d\d)-(\d+))[ \t]*$", re.M)
FIELD = re.compile(r"^- \*\*([^*]+):\*\* (.*)$", re.M)
REF = re.compile(r"REQ-L[123]-\d\d-\d+")
LINK = re.compile(r"\]\(([^)]+\.md)(?:#[^)]*)?\)")
SECTION = re.compile(r"§(\d+(?:\.\d+)*)\s*«([^»]+)»")


def source_references(source):
    matches = list(LINK.finditer(source))
    for index, match in enumerate(matches):
        end = matches[index + 1].start() if index + 1 < len(matches) else len(source)
        yield match.group(1), SECTION.findall(source[match.end():end])


def norm(text):
    text = re.sub(r"<[^>]+>|[*`$\\{}^_]", "", text)
    text = text.replace("’", "'").replace("–", "-").replace("—", "-")
    return re.sub(r"\s+", " ", text).strip().lower()


def stems(text):
    return {w[:5] for w in re.findall(r"[\w'-]+", text) if len(w) > 2}


def title_matches(cited, heading):
    """Cited title matches when it is a prefix of the heading (or vice versa)
    or when at least 60 % of its word stems occur in the heading."""
    if heading.startswith(cited) or cited.startswith(heading):
        return True
    words = stems(cited)
    return bool(words) and len(words & stems(heading)) / len(words) >= 0.6


def numbered_heading(line):
    """Return (number, title) for a section heading whose level matches its
    numbering depth (## 4., ### 4.1., #### 4.1.1.); local lists such as
    '#### 1. ...' inside a section are not sections."""
    m = re.match(r"^(#{2,5}) (\d+(?:\.\d+)*)\.?\s+(.*)$", line)
    if m and m.group(2).count(".") + 1 >= len(m.group(1)) - 1:
        return m.group(2), m.group(3)
    return None


def chapter_headings(path, cache={}):
    if path not in cache:
        heads = defaultdict(list)
        in_code = False
        for line in path.read_text(encoding="utf-8").splitlines():
            if line.lstrip().startswith("```"):
                in_code = not in_code
                continue
            h = None if in_code else numbered_heading(line)
            if h:
                heads[h[0]].append(norm(h[1]))
        cache[path] = heads
    return cache[path]


def section_bodies(path, cache={}):
    """Normalised text of each numbered section, including its subsections."""
    if path not in cache:
        lines = path.read_text(encoding="utf-8").splitlines()
        marks = []
        in_code = False
        for i, line in enumerate(lines):
            if line.lstrip().startswith("```"):
                in_code = not in_code
                continue
            h = None if in_code else numbered_heading(line)
            if h:
                marks.append((i, h[0]))
        bodies = defaultdict(str)
        for k, (start, num) in enumerate(marks):
            end = len(lines)
            for nxt, other in marks[k + 1:]:
                if not other.startswith(num + "."):
                    end = nxt
                    break
            bodies[num] += " " + norm(" ".join(lines[start:end]))
        cache[path] = bodies
    return cache[path]


def load():
    reqs = {}
    problems = []
    header_warnings = []
    for level, dirname in LEVEL_DIRS.items():
        for path in sorted((REQ_ROOT / dirname).glob("ch*.md")):
            chapter = re.match(r"ch(\d\d)-", path.name).group(1)
            text = path.read_text(encoding="utf-8").replace("\r\n", "\n")
            head = text.split("\n---\n", 1)[0].split("\n")
            for a, b in zip(head, head[1:]):
                if a.startswith("**") and b.startswith("**"):
                    header_warnings.append(f"{path.name}: header fields without blank line: {a[:40]!r}")
                    break
            matches = list(REQ_HEAD.finditer(text))
            for i, m in enumerate(matches):
                rid = m.group(1)
                end = matches[i + 1].start() if i + 1 < len(matches) else len(text)
                fields = dict(FIELD.findall(text[m.end():end]))
                if m.group(2) != level or m.group(3) != chapter:
                    problems.append(f"{rid}: ID does not match file {dirname}/{path.name}")
                if rid in reqs:
                    problems.append(f"{rid}: duplicated ID ({reqs[rid]['file']} and {path.name})")
                reqs[rid] = {"level": level, "chapter": chapter, "file": f"{dirname}/{path.name}",
                             "path": path, "fields": fields}
    return reqs, problems, header_warnings


def main():
    check_sections = "--check-sections" in sys.argv
    quiet = "--quiet" in sys.argv
    reqs, errors, warnings = load()

    for rid, r in reqs.items():
        f = r["fields"]
        for name in MANDATORY:
            if name not in f:
                errors.append(f"{rid}: missing field «{name}»")

        parent = f.get("Батьківська вимога", "")
        if r["level"] == "L1":
            if REF.search(parent):
                errors.append(f"{rid}: L1 requirement must not have a parent ({parent})")
        else:
            refs = REF.findall(parent)
            if not refs:
                errors.append(f"{rid}: parent missing or malformed ({parent!r})")
            for p in refs:
                if p not in reqs:
                    errors.append(f"{rid}: parent {p} does not exist")
                elif reqs[p]["level"] != PARENT_LEVEL[r["level"]]:
                    errors.append(f"{rid}: parent {p} is not one level up")
                elif reqs[p]["chapter"] != r["chapter"]:
                    warnings.append(f"{rid}: cross-chapter parent {p}")

        canon = f.get("Канонічна вимога")
        if canon is not None:
            refs = REF.findall(canon)
            if len(refs) != 1:
                errors.append(f"{rid}: «Канонічна вимога» must name exactly one requirement")
            for c in refs:
                if c == rid:
                    errors.append(f"{rid}: canonical requirement points to itself")
                elif c not in reqs:
                    errors.append(f"{rid}: canonical requirement {c} does not exist")
                elif reqs[c]["level"] != r["level"]:
                    errors.append(f"{rid}: canonical requirement {c} is on another level")
                elif "Канонічна вимога" in reqs[c]["fields"]:
                    errors.append(f"{rid}: canonical requirement {c} is itself marked as a duplicate")

        method = f.get("Метод верифікації", "")
        if method and not method.startswith(METHODS):
            errors.append(f"{rid}: unknown verification method {method!r}")

        source = f.get("Джерело", "")
        for link, sections in source_references(source):
            target = (r["path"].parent / link).resolve()
            if not target.exists():
                errors.append(f"{rid}: source link {link} does not exist")
            elif check_sections:
                heads = chapter_headings(target)
                for num, title in sections:
                    got = heads.get(num)
                    want = norm(title)
                    if not got:
                        warnings.append(f"{rid}: §{num} not found in {target.name}")
                    elif not any(title_matches(want, g) for g in got) and want not in section_bodies(target).get(num, ""):
                        warnings.append(f"{rid}: §{num} «{title}» differs from heading «{got[0][:80]}»")

        text = f.get("Вимога", "")
        for pattern, message in BANNED:
            if pattern.search(text):
                warnings.append(f"{rid}: {message}")

    per_level = defaultdict(int)
    for r in reqs.values():
        per_level[r["level"]] += 1
    duplicates = sum(1 for r in reqs.values() if "Канонічна вимога" in r["fields"])

    print("=== Requirements traceability inspection ===")
    print("Requirements: " + ", ".join(f"{k} {per_level[k]}" for k in sorted(per_level))
          + f"; marked duplicates: {duplicates}")
    if not quiet:
        for w in warnings:
            print(f"  [WARN] {w}")
    for e in errors:
        print(f"  [ERR]  {e}")
    print(f"Errors: {len(errors)}; warnings: {len(warnings)}")
    return 1 if errors else 0


if __name__ == "__main__":
    sys.exit(main())
