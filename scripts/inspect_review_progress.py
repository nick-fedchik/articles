#!/usr/bin/env python3
"""Report manually recorded review progress for ExpertSystem chapters 01-10.

Section hashes cover the original UTF-8 bytes from the H2 heading through the
byte before the next H1/H2 or end of file, including original line endings.
Requirement hashes cover original requirement block bytes, excluding trailing
line breaks. Dependencies are explicit IDs and requirements in the chapter
whose source cites the numbered H2 section or a subsection via §N notation.
Unlinked and uncited semantic relationships are not discovered; manual review
is needed. The fingerprint command prints a template only; the reviewer must
supply status, notes and timestamp.
Stale entries are counted separately regardless of their recorded status.
"""

import argparse
import hashlib
import json
import re
import sys
from datetime import datetime
from pathlib import Path

from inspect_requirements_trace import FIELD, LEVEL_DIRS, REQ_HEAD

ROOT = Path(__file__).resolve().parent.parent / "ExpertSystem"
CHAPTERS = tuple(f"{number:02d}" for number in range(1, 11))
LEDGER = "ch01-10-review-progress.json"
SHA256 = re.compile(r"[0-9a-f]{64}\Z")
REQUIREMENT_ID = re.compile(r"REQ-L[123]-\d\d-\d+\Z")
SECTION_NUMBER = re.compile(r"^(\d+(?:\.\d+)*)(?:\.(?=\s|$)|(?=\s|$))")
SOURCE_NUMBER = re.compile(r"§\s*(\d+(?:\.\d+)*)(?!\d|\.\d)")
FENCE = re.compile(r"^ {0,3}(`{3,}|~{3,})(.*)$")
HEADING = re.compile(r"^ {0,3}(#{1,2})(?:[ \t]+(.*)|[ \t]*)$")
RECORD_FIELDS = {"chapter", "section", "status", "reviewer", "notes", "timestamp", "section_sha256",
                 "requirement_sha256", "requirement_ids", "decision_references"}


def digest(data):
    return hashlib.sha256(data).hexdigest()


def one_file(directory, pattern):
    paths = sorted(directory.glob(pattern))
    if len(paths) != 1:
        raise ValueError(f"Expected exactly one {directory / pattern}, found {len(paths)}")
    return paths[0]


def sections(path):
    lines = path.read_bytes().splitlines(keepends=True)
    marks = []
    fence = None
    for index, raw in enumerate(lines):
        line = raw.decode("utf-8").rstrip("\r\n")
        match = FENCE.match(line)
        if fence:
            if match and match.group(1)[0] == fence[0] and len(match.group(1)) >= len(fence) and not match.group(2).strip():
                fence = None
            continue
        if match and (match.group(1)[0] != "`" or "`" not in match.group(2)):
            fence = match.group(1)
            continue
        heading = HEADING.match(line)
        if heading:
            title = re.sub(r"(?:^|[ \t]+)#+[ \t]*$", "", heading.group(2) or "").strip()
            marks.append((index, len(heading.group(1)), title))
    result = {}
    for position, (start, level, title) in enumerate(marks):
        if level != 2:
            continue
        if title in result:
            raise ValueError(f"{path.name}: duplicate H2 heading: {title!r}")
        end = marks[position + 1][0] if position + 1 < len(marks) else len(lines)
        result[title] = digest(b"".join(lines[start:end]))
    return result


def requirements(root):
    requirement_root = root / "requirements"
    by_chapter = {chapter: {"count": 0, "with_check": 0} for chapter in CHAPTERS}
    items = {}
    for level, dirname in LEVEL_DIRS.items():
        for path in sorted((requirement_root / dirname).glob("ch*.md")):
            lines = path.read_bytes().decode("utf-8").splitlines(keepends=True)
            starts = []
            for index, line in enumerate(lines):
                match = REQ_HEAD.fullmatch(line.rstrip("\r\n"))
                if match:
                    starts.append((index, match))
            for index, (start, match) in enumerate(starts):
                identifier = match.group(1)
                if identifier in items:
                    raise ValueError(f"Duplicate requirement ID: {identifier}")
                chapter = match.group(3)
                end = starts[index + 1][0] if index + 1 < len(starts) else len(lines)
                block = "".join(lines[start:end]).rstrip("\r\n")
                fields = dict(FIELD.findall(block.replace("\r\n", "\n")))
                items[identifier] = {"chapter": chapter, "hash": digest(block.encode("utf-8")),
                                     "source": fields.get("Джерело", "")}
                if chapter in by_chapter:
                    by_chapter[chapter]["count"] += 1
                    by_chapter[chapter]["with_check"] += int("Перевірка" in fields)
        for chapter in CHAPTERS:
            one_file(requirement_root / dirname, f"ch{chapter}-{level}-*.md")
    return by_chapter, items


def dependency_hashes(chapter, section, explicit_ids, items):
    unknown = set(explicit_ids) - items.keys()
    if unknown:
        raise ValueError(f"{chapter} / {section}: nonexistent requirement IDs: {', '.join(sorted(unknown))}")
    selected = set(explicit_ids)
    heading = SECTION_NUMBER.match(section)
    if heading:
        number = heading.group(1)
        for identifier, item in items.items():
            if item["chapter"] == chapter:
                for cited in SOURCE_NUMBER.findall(item["source"]):
                    if cited == number or cited.startswith(number + "."):
                        selected.add(identifier)
                        break
    return {identifier: items[identifier]["hash"] for identifier in sorted(selected)}


def reject_duplicate_keys(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f"Duplicate JSON key: {key}")
        result[key] = value
    return result


def read_ledger(path):
    data = json.loads(path.read_text(encoding="utf-8"), object_pairs_hook=reject_duplicate_keys)
    if not isinstance(data, dict) or set(data) != {"version", "records"} or type(data["version"]) is not int or data["version"] != 1 or not isinstance(data["records"], list):
        raise ValueError("Ledger must contain version 1 and a records array only")
    return data["records"]


def validate_record(record, all_sections, requirement_ids):
    if not isinstance(record, dict) or set(record) != RECORD_FIELDS:
        raise ValueError(f"Malformed record: expected fields {sorted(RECORD_FIELDS)}")
    chapter, section = record["chapter"], record["section"]
    if not isinstance(chapter, str) or chapter not in all_sections or not isinstance(section, str) or section not in all_sections[chapter]:
        raise ValueError(f"Orphaned or invalid section: {chapter!r} / {section!r}")
    if record["status"] not in ("reviewed", "partial", "blocked"):
        raise ValueError(f"{chapter} / {section}: invalid status")
    if not isinstance(record["reviewer"], str) or not record["reviewer"].strip():
        raise ValueError(f"{chapter} / {section}: reviewer must be nonempty")
    if not isinstance(record["notes"], str) or not record["notes"].strip():
        raise ValueError(f"{chapter} / {section}: notes must be nonempty")
    timestamp = record["timestamp"]
    try:
        parsed = datetime.fromisoformat(timestamp.replace("Z", "+00:00")) if isinstance(timestamp, str) else None
        if parsed is None or parsed.tzinfo is None or parsed.utcoffset() is None or "T" not in timestamp:
            raise ValueError
    except ValueError:
        raise ValueError(f"{chapter} / {section}: timestamp must be ISO 8601 with timezone") from None
    hashes = record["requirement_sha256"]
    if not isinstance(record["section_sha256"], str) or not SHA256.fullmatch(record["section_sha256"]) or not isinstance(hashes, dict) or any(not isinstance(key, str) or not REQUIREMENT_ID.fullmatch(key) or not isinstance(value, str) or not SHA256.fullmatch(value) for key, value in hashes.items()):
        raise ValueError(f"{chapter} / {section}: invalid SHA-256 hashes")
    for field in ("requirement_ids", "decision_references"):
        values = record[field]
        if not isinstance(values, list) or any(not isinstance(value, str) or not value.strip() for value in values) or len(values) != len(set(values)):
            raise ValueError(f"{chapter} / {section}: invalid or duplicate {field}")
    unknown = set(record["requirement_ids"]) - requirement_ids
    if unknown:
        raise ValueError(f"{chapter} / {section}: nonexistent requirement IDs: {', '.join(sorted(unknown))}")
    return chapter, section


def calculate(root):
    all_sections = {chapter: sections(one_file(root, f"ch{chapter}-*.md")) for chapter in CHAPTERS}
    reqs, items = requirements(root)
    records = read_ledger(root / "requirements" / LEDGER)
    indexed = {}
    stale_records = []
    for record in records:
        key = validate_record(record, all_sections, items.keys())
        if key in indexed:
            raise ValueError(f"Duplicate review record: {key[0]} / {key[1]}")
        indexed[key] = record
    chapters = {}
    for chapter in CHAPTERS:
        counts = {"sections": len(all_sections[chapter]), "reviewed": 0, "partial": 0,
                  "blocked": 0, "stale": 0, "pending": 0,
                  "requirements": reqs[chapter]["count"], "requirements_with_check": reqs[chapter]["with_check"]}
        for title, section_hash in all_sections[chapter].items():
            record = indexed.get((chapter, title))
            if record is None:
                counts["pending"] += 1
            elif record["section_sha256"] != section_hash or record["requirement_sha256"] != dependency_hashes(chapter, title, record["requirement_ids"], items):
                counts["stale"] += 1
                changed = []
                if record["section_sha256"] != section_hash:
                    changed.append("section")
                current = dependency_hashes(chapter, title, record["requirement_ids"], items)
                changed.extend(identifier for identifier in sorted(record["requirement_sha256"].keys() | current.keys())
                               if record["requirement_sha256"].get(identifier) != current.get(identifier))
                stale_records.append({"chapter": chapter, "section": title, "status": record["status"], "changed": changed})
            else:
                counts[record["status"]] += 1
        counts["reviewed_percent"] = 100 * counts["reviewed"] / counts["sections"] if counts["sections"] else 0.0
        chapters[chapter] = counts
    total = {key: sum(counts[key] for counts in chapters.values()) for key in next(iter(chapters.values())) if key != "reviewed_percent"}
    total["reviewed_percent"] = 100 * total["reviewed"] / total["sections"] if total["sections"] else 0.0
    return {"chapters": chapters, "total": total, "stale_records": stale_records,
            "scope": "Unlinked and uncited semantic relationships are not automatically discovered; manual review is needed."}, all_sections, items


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__, epilog="Templates are not valid ledger records until a reviewer supplies a status, reviewer identity, nonempty notes, and an ISO 8601 timestamp with timezone. Lists may be empty; decision references are reviewer-supplied strings.")
    parser.add_argument("--json", action="store_true", help="Print report as JSON")
    parser.add_argument("--fingerprint", action="store_true", help="Print a record template without writing")
    parser.add_argument("--chapter", choices=CHAPTERS)
    parser.add_argument("--section", help="Exact H2 heading text")
    parser.add_argument("--requirement-id", action="append", default=[], help="Explicit requirement ID; repeat for each linked requirement")
    args = parser.parse_args(argv)
    if args.fingerprint != (args.chapter is not None and args.section is not None) or (not args.fingerprint and (args.chapter is not None or args.section is not None)):
        parser.error("--fingerprint requires both --chapter and --section; those options require --fingerprint")
    if args.requirement_id and not args.fingerprint:
        parser.error("--requirement-id requires --fingerprint")
    if len(args.requirement_id) != len(set(args.requirement_id)) or any(not REQUIREMENT_ID.fullmatch(identifier) for identifier in args.requirement_id):
        parser.error("--requirement-id values must be valid and unique")
    try:
        report, all_sections, items = calculate(ROOT)
        if args.fingerprint:
            if args.section not in all_sections[args.chapter]:
                raise ValueError(f"Unknown section: {args.chapter} / {args.section}")
            record = {"chapter": args.chapter, "section": args.section, "status": None,
                      "reviewer": "", "notes": "", "timestamp": None, "section_sha256": all_sections[args.chapter][args.section],
                      "requirement_sha256": dependency_hashes(args.chapter, args.section, args.requirement_id, items), "requirement_ids": args.requirement_id, "decision_references": []}
            print(json.dumps(record, ensure_ascii=False, indent=2))
        elif args.json:
            print(json.dumps(report, ensure_ascii=False, indent=2))
        else:
            print("Manual review: current reviewed / all H2 sections; stale records are separate.")
            print("Requirement counts and explicit Перевірка fields are structural signals, not quality evidence.")
            print(report["scope"])
            for chapter, counts in list(report["chapters"].items()) + [("total", report["total"])]:
                print(f"{chapter}: reviewed {counts['reviewed']}/{counts['sections']} ({counts['reviewed_percent']:.1f}%); "
                      f"partial {counts['partial']}; blocked {counts['blocked']}; stale {counts['stale']}; pending {counts['pending']}; "
                      f"requirements {counts['requirements']}, with Перевірка {counts['requirements_with_check']}")
            for record in report["stale_records"]:
                print(f"  STALE {record['chapter']} / {record['section']} ({record['status']}): {', '.join(record['changed'])}")
    except (OSError, UnicodeError, ValueError, TypeError, json.JSONDecodeError) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())