import json
import sys
import tempfile
import unittest
from contextlib import redirect_stdout
from io import StringIO
from pathlib import Path
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))
import inspect_review_progress as progress


class ReviewProgressTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        self.req_root = self.root / "requirements"
        for chapter in progress.CHAPTERS:
            (self.root / f"ch{chapter}-chapter.md").write_bytes(b"# Chapter\r\n")
            for level, dirname in progress.LEVEL_DIRS.items():
                directory = self.req_root / dirname
                directory.mkdir(parents=True, exist_ok=True)
                (directory / f"ch{chapter}-{level}-fixture.md").write_text("# Requirements\n", encoding="utf-8")
        self.write_ledger([])

    def write_ledger(self, records):
        (self.req_root / progress.LEDGER).write_text(json.dumps({"version": 1, "records": records}), encoding="utf-8")

    def chapter(self, text):
        (self.root / "ch01-chapter.md").write_bytes(text.encode("utf-8"))

    def record(self, status="reviewed", section="Section", requirement_ids=None):
        _, sections, items = progress.calculate(self.root)
        requirement_ids = requirement_ids or []
        return {"chapter": "01", "section": section, "status": status, "reviewer": "Reviewer", "notes": "Manual review",
                "timestamp": "2026-10-27T12:00:00+00:00", "section_sha256": sections["01"][section],
                "requirement_sha256": progress.dependency_hashes("01", section, requirement_ids, items),
                "requirement_ids": requirement_ids, "decision_references": []}

    def test_fences_boundaries_and_exact_bytes(self):
        self.chapter("# Chapter\r\n## Section\r\n### Subsection\r\n```md\r\n## Fake\r\n~~~\r\n```\r\n~~~~\r\n## Also fake\r\n~~~\r\n~~~~\r\n## Next\r\n# End\r\n## Final\r\n")
        found = progress.sections(self.root / "ch01-chapter.md")
        self.assertEqual(list(found), ["Section", "Next", "Final"])
        self.assertEqual(found["Next"], progress.digest(b"## Next\r\n"))

    def test_duplicate_heading_fails(self):
        self.chapter("## Same\n## Same\n")
        with self.assertRaisesRegex(ValueError, "duplicate H2"):
            progress.calculate(self.root)

    def test_empty_headings_are_counted_and_end_sections(self):
        self.chapter("## Section\nBody\n#\nIgnored\n##\nEmpty\n")
        found = progress.sections(self.root / "ch01-chapter.md")
        self.assertEqual(list(found), ["Section", ""])
        self.assertEqual(found["Section"], progress.digest(b"## Section\nBody\n"))

    def test_current_stale_and_blocked(self):
        self.chapter("## Section\nBody\n## Blocked\n")
        current = self.record()
        blocked = dict(current, section="Blocked", status="blocked")
        blocked["section_sha256"] = progress.sections(self.root / "ch01-chapter.md")["Blocked"]
        self.write_ledger([current, blocked])
        counts = progress.calculate(self.root)[0]["chapters"]["01"]
        self.assertEqual((counts["reviewed"], counts["blocked"], counts["stale"]), (1, 1, 0))
        self.chapter("## Section\nChanged\n## Blocked\n")
        report = progress.calculate(self.root)[0]
        counts = report["chapters"]["01"]
        self.assertEqual((counts["reviewed"], counts["blocked"], counts["stale"]), (0, 1, 1))
        self.assertEqual(report["stale_records"][0]["changed"], ["section"])

    def test_partial_and_linked_requirement(self):
        self.chapter("## Section\n")
        path = self.req_root / "L3-component" / "ch01-L3-fixture.md"
        path.write_text("## REQ-L3-01-01\n", encoding="utf-8")
        record = self.record("partial", requirement_ids=["REQ-L3-01-01"])
        record["decision_references"] = ["decision-1"]
        self.write_ledger([record])
        counts = progress.calculate(self.root)[0]["chapters"]["01"]
        self.assertEqual((counts["partial"], counts["pending"], counts["reviewed"]), (1, 0, 0))

    def test_empty_and_structural_signal(self):
        report = progress.calculate(self.root)[0]
        self.assertEqual(report["total"]["reviewed_percent"], 0.0)
        path = self.req_root / "L1-stakeholder" / "ch01-L1-fixture.md"
        path.write_text("## REQ-L1-01-01\n- **Перевірка:** Test\n\n## REQ-L1-01-02\n- **Метод верифікації:** Test\n", encoding="utf-8")
        counts = progress.calculate(self.root)[0]["chapters"]["01"]
        self.assertEqual((counts["requirements"], counts["requirements_with_check"]), (2, 1))

    def test_dependency_changes_and_section_number_boundaries(self):
        self.chapter("## 5. Section\nBody\n## 50. Other\n")
        path = self.req_root / "L1-stakeholder" / "ch01-L1-fixture.md"
        cross = self.req_root / "L2-system" / "ch02-L2-fixture.md"
        initial = ("## REQ-L1-01-01\n- **Джерело:** §4, §5.1 and §8\n- **Вимога:** cited\n\n"
                   "## REQ-L1-01-02\n- **Джерело:** §50\n- **Вимога:** unrelated\n")
        path.write_text(initial, encoding="utf-8")
        cross.write_text("## REQ-L2-02-01\n- **Вимога:** linked\n", encoding="utf-8")
        record = self.record(section="5. Section", requirement_ids=["REQ-L2-02-01"])
        self.assertEqual(set(record["requirement_sha256"]), {"REQ-L1-01-01", "REQ-L2-02-01"})
        self.write_ledger([record])
        path.write_text(initial.replace("unrelated", "changed") + "\n## REQ-L1-01-03\n- **Джерело:** §50.1\n", encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["chapters"]["01"]["reviewed"], 1)
        path.write_text(initial.replace("cited", "edited"), encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-01"])
        path.write_text(initial + "\n## REQ-L1-01-03\n- **Джерело:** §5\n", encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-03"])
        path.write_text(initial + "\n## REQ-L1-01-03\n- **Джерело:** §5.\n", encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-03"])
        path.write_text("## REQ-L1-01-02" + initial.split("## REQ-L1-01-02")[1], encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-01"])
        path.write_text(initial.replace("§5.1", "§6.1"), encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-01"])
        path.write_text(initial, encoding="utf-8")
        cross.write_text("## REQ-L2-02-01\n- **Вимога:** changed\n", encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L2-02-01"])

    def test_non_numbered_and_fingerprint_cli(self):
        self.chapter("## Section\n")
        path = self.req_root / "L1-stakeholder" / "ch01-L1-fixture.md"
        path.write_text("## REQ-L1-01-01\n- **Джерело:** §5\n", encoding="utf-8")
        self.assertEqual(self.record()["requirement_sha256"], {})
        output = StringIO()
        with patch.object(progress, "ROOT", self.root), redirect_stdout(output):
            self.assertEqual(progress.main(["--fingerprint", "--chapter", "01", "--section", "Section",
                                            "--requirement-id", "REQ-L1-01-01"]), 0)
        template = json.loads(output.getvalue())
        self.assertEqual(template["requirement_ids"], ["REQ-L1-01-01"])
        self.assertEqual(set(template["requirement_sha256"]), {"REQ-L1-01-01"})
        self.write_ledger([self.record(requirement_ids=["REQ-L1-01-01"])])
        path.write_text("## REQ-L1-01-01\n- **Джерело:** §5\n- **Вимога:** changed\n", encoding="utf-8")
        self.assertEqual(progress.calculate(self.root)[0]["stale_records"][0]["changed"], ["REQ-L1-01-01"])
        path.write_text("# Requirement removed\n", encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "nonexistent requirement IDs"):
            progress.calculate(self.root)

    def test_migrated_record_preserves_review_metadata(self):
        self.chapter("## 5. Section\n")
        path = self.req_root / "L1-stakeholder" / "ch01-L1-fixture.md"
        path.write_text("## REQ-L1-01-01\n- **Джерело:** §5.2\n", encoding="utf-8")
        record = self.record(section="5. Section", requirement_ids=["REQ-L1-01-01"])
        record.update(notes="Existing review and unresolved questions", decision_references=["decision-1"])
        legacy = dict(record, requirement_sha256={level: "0" * 64 for level in progress.LEVEL_DIRS})
        _, items = progress.requirements(self.root)
        migrated = dict(legacy, requirement_sha256=progress.dependency_hashes(
            legacy["chapter"], legacy["section"], legacy["requirement_ids"], items))
        self.assertEqual(migrated, record)
        self.write_ledger([migrated])
        report = progress.calculate(self.root)[0]
        self.assertEqual((report["total"]["reviewed"], report["total"]["stale"]), (1, 0))
        self.assertEqual(progress.read_ledger(self.req_root / progress.LEDGER)[0], record)

    def test_existing_ledger_migration(self):
        records = progress.read_ledger(progress.ROOT / "requirements" / progress.LEDGER)
        record = next(record for record in records if record["chapter"] == "09" and record["section"] ==
                      "5. Архітектура промислової платформи управління графом знань")
        self.assertEqual(record["status"], "reviewed")
        self.assertEqual(record["timestamp"], "2026-10-27T11:44:21+02:00")
        self.assertIn("Зафіксовано відкриті питання", record["notes"])
        self.assertIn("REQ-L2-07-06", record["requirement_ids"])
        self.assertIn("REQ-L2-07-06", record["requirement_sha256"])
        self.assertEqual(record["decision_references"], ["А-06", "А-13", "А-15"])

    def test_invalid_or_duplicate_records(self):
        self.chapter("## Section\n")
        record = self.record()
        for invalid in (dict(record, status="done"), dict(record, section="Missing"),
                        dict(record, requirement_ids=["REQ-L1-01-99"]), dict(record, timestamp="yesterday"),
                        dict(record, unexpected=True), dict(record, section_sha256="bad"), dict(record, reviewer=""),
                        dict(record, requirement_ids=["x", "x"])):
            self.write_ledger([invalid])
            with self.assertRaises(ValueError):
                progress.calculate(self.root)
        self.write_ledger([record, record])
        with self.assertRaisesRegex(ValueError, "Duplicate review record"):
            progress.calculate(self.root)
        (self.req_root / progress.LEDGER).write_text('{"version": 1, "version": 1, "records": []}', encoding="utf-8")
        with self.assertRaisesRegex(ValueError, "Duplicate JSON key"):
            progress.calculate(self.root)


if __name__ == "__main__":
    unittest.main()