import copy
import json
import tempfile
import unittest
from datetime import datetime, timedelta, timezone
from pathlib import Path

from review import check_execution as checker


class ExecutionTests(unittest.TestCase):
    def setUp(self):
        temporary = checker.ROOT / "review/tmp"
        temporary.mkdir(exist_ok=True)
        fixture = tempfile.TemporaryDirectory(prefix="execution-test-", dir=temporary)
        self.addCleanup(fixture.cleanup)
        self.root = Path(fixture.name)
        (self.root / "ExpertSystem").mkdir()
        (self.root / "review/tmp").mkdir(parents=True)
        for number in range(1, 41):
            (self.root / f"ExpertSystem/ch{number:02}-test.md").write_text(f"# Chapter {number}\n", encoding="utf-8")
        (self.root / "ExpertSystem/README.md").write_text("# Companion\n", encoding="utf-8")
        (self.root / "ExpertSystem/AGENTS.md").write_text("Not a companion", encoding="utf-8")
        self.now = datetime(2026, 10, 29, 12, 0, tzinfo=timezone.utc)
        self.data = checker.new_ledger(self.root, self.now)

    def check(self, data=None, now=None):
        return checker.check(self.root, self.data if data is None else data, now or self.now)

    def completed(self):
        data = copy.deepcopy(self.data)
        output = self.root / "review/accepted.md"
        output.write_text("Verified control artifact\n", encoding="utf-8")
        task = data["tasks"][0]
        task.update(status="completed", outputs=[{"path": "review/accepted.md", "sha256": checker.sha256(output)}],
                    acceptance={"by": "main-reviewer", "at": data["updated_at"],
                                "evidence": "Verified coverage, contents and recorded artifact hashes."})
        return data

    def test_valid_all40_and_report(self):
        self.assertEqual(self.check(), [])
        self.assertEqual(len(self.data["chapters"]), 40)
        self.assertEqual(len(self.data["companions"]), 1)
        self.assertEqual(len(self.data["stages"]), 8)
        self.assertEqual(len(self.data["tasks"]), 5)
        self.assertTrue(all(record["scientific_review"] == "pending" for record in self.data["chapters"]))
        checker.report(self.root, self.data, [])
        text = (self.root / "review/PROGRESS.md").read_text(encoding="utf-8")
        for number in range(1, 41):
            self.assertIn(f"| {number:02}: [Chapter {number}]", text)
        self.assertIn("not scientific truth", text)
        self.assertIn("Planned | Doing | Findings | Next Action", text)

    def test_omitted_and_duplicate_chapter(self):
        for records in (self.data["chapters"][:-1], self.data["chapters"][:-1] + [self.data["chapters"][0]]):
            with self.subTest(records=len(records)):
                data = copy.deepcopy(self.data)
                data["chapters"] = records
                self.assertIn("coverage", " ".join(self.check(data)))

    def test_source_missing_or_duplicate_chapter(self):
        path = self.root / "ExpertSystem/ch40-test.md"
        path.unlink()
        self.assertIn("coverage", " ".join(self.check()))
        path.write_text("# Chapter 40\n", encoding="utf-8")
        (self.root / "ExpertSystem/ch01-duplicate.md").write_text("# Duplicate\n", encoding="utf-8")
        self.assertIn("coverage", " ".join(self.check()))

    def test_duplicate_task_and_json_keys(self):
        self.data["tasks"].append(copy.deepcopy(self.data["tasks"][0]))
        self.assertIn("duplicate task IDs", " ".join(self.check()))
        path = self.root / "review/execution.json"
        for text in ('{"tasks": [], "tasks": []}', '{"acceptance": {"by": "a", "by": "b"}}'):
            with self.subTest(text=text):
                path.write_text(text, encoding="utf-8")
                with self.assertRaisesRegex(ValueError, "duplicate JSON key"):
                    checker.load_ledger(path)

    def test_valid_completed(self):
        self.assertEqual(self.check(self.completed()), [])

    def test_stage_completion_requires_accepted_covering_tasks(self):
        for index in (0, 1, 2):
            with self.subTest(stage=index + 1):
                data = copy.deepcopy(self.data)
                data["stages"][index]["status"] = "completed"
                self.assertIn("stage completion requires accepted tasks", " ".join(self.check(data)))
        data = self.completed()
        data["stages"][0]["status"] = "completed"
        self.assertEqual(self.check(data), [])

    def test_scientific_completion_requires_review_task(self):
        self.data["chapters"][0]["scientific_review"] = "completed"
        self.assertIn("scientific completion requires an accepted S3 task", " ".join(self.check()))

    def test_invalid_completed_artifact(self):
        data = self.completed()
        for path in ("review/missing.md", "review/tmp/draft.md", "ExpertSystem/README.md"):
            with self.subTest(path=path):
                changed = copy.deepcopy(data)
                changed["tasks"][0]["outputs"][0]["path"] = path
                if path == "review/tmp/draft.md":
                    (self.root / path).write_text("Draft\n", encoding="utf-8")
                self.assertTrue(self.check(changed))
        data["tasks"][0]["outputs"] = []
        self.assertIn("no permanent output", " ".join(self.check(data)))

    def test_output_hash_and_acceptance(self):
        data = self.completed()
        data["tasks"][0]["outputs"][0]["sha256"] = "0" * 64
        self.assertIn("output hash mismatch", " ".join(self.check(data)))
        for acceptance in (None, {}, {"by": "main-reviewer", "at": data["updated_at"], "evidence": " "}):
            with self.subTest(acceptance=acceptance):
                changed = self.completed()
                changed["tasks"][0]["acceptance"] = acceptance
                self.assertIn("acceptance evidence", " ".join(self.check(changed)))

    def test_changed_source_and_stale_task_input_after_refresh(self):
        data = self.completed()
        source = self.root / "ExpertSystem/ch01-test.md"
        source.write_text("# Chapter 1\nRevised science.\n", encoding="utf-8")
        self.assertIn("source snapshot hash mismatch", " ".join(self.check(data)))
        data["chapters"], data["companions"] = checker.inventory(self.root)
        self.assertIn("OPS-01: input hash mismatch", " ".join(self.check(data)))

    def test_companion_inventory_drift(self):
        (self.root / "ExpertSystem/part-new.md").write_text("# New part\n", encoding="utf-8")
        self.assertIn("companions paths or inventory changed", " ".join(self.check()))

    def test_overdue_running_not_paused(self):
        later = self.now + timedelta(minutes=10, seconds=1)
        self.assertEqual(len(self.check(now=later)), 5)
        self.assertTrue(all("OVERDUE" in issue and "suspicion" in issue for issue in self.check(now=later)))
        for task in self.data["tasks"]:
            task["status"] = "paused"
        self.assertEqual(self.check(now=later), [])

    def test_next_action_skips_paused_and_honors_resume(self):
        self.data["tasks"][0].update(status="paused", next_action="Deferred visual session")
        self.data["tasks"][1]["next_action"] = "First unpaused chapter"
        checker.report(self.root, self.data, [])
        text = (self.root / "review/PROGRESS.md").read_text(encoding="utf-8")
        self.assertIn("## Next Action\n\nFirst unpaused chapter", text)
        self.data["resume_next_task"] = "MAP-04"
        self.data["tasks"][-1]["next_action"] = "Resume the recorded next chapter"
        checker.report(self.root, self.data, [])
        text = (self.root / "review/PROGRESS.md").read_text(encoding="utf-8")
        self.assertIn("## Next Action\n\nResume the recorded next chapter", text)

    def test_malformed_timestamps_and_paths(self):
        for stamp in ("yesterday", "2026-10-29T12:00:00+02:00", "2026-99-29T12:00:00Z"):
            with self.subTest(stamp=stamp):
                data = copy.deepcopy(self.data)
                data["tasks"][0]["checkpoint_at"] = stamp
                self.assertTrue(self.check(data))
        for path in ("../escape.md", "/etc/passwd", "review//accepted.md", "review/../accepted.md"):
            with self.subTest(path=path):
                data = self.completed()
                data["tasks"][0]["outputs"][0]["path"] = path
                self.assertIn("malformed path", " ".join(self.check(data)))

    def test_symlink_cannot_promote_temporary_output(self):
        data = self.completed()
        temporary = self.root / "review/tmp/draft.md"
        temporary.write_text("Draft\n", encoding="utf-8")
        link = self.root / "review/disguised.md"
        link.symlink_to(temporary)
        data["tasks"][0]["outputs"] = [{"path": "review/disguised.md", "sha256": checker.sha256(temporary)}]
        self.assertIn("permanent output", " ".join(self.check(data)))


if __name__ == "__main__":
    unittest.main()
