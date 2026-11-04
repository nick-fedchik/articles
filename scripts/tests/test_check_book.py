import contextlib
import io
import json
import shutil
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path
from unittest import mock


ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / "scripts"))
import check_book
import inspect_math_formatting


VALID_CHAPTER = """\
# Глава 1. Навчальний приклад

## Анотація

Опис.

<!-- chapter-toc:start -->
## Зміст глави

- [1. Метод](#1-метод)
  - [1.1. Приклад](#11-приклад)
- [Висновки](#висновки)
<!-- chapter-toc:end -->

## 1. Метод

### 1.1. Приклад

<a id="example"></a>

**Лістинг 1.1. Приклад (Python 3.14.4).**

<details>
<summary>Показати лістинг 1.1 (2 рядки, Python 3.14.4)</summary>

```python
x = 1
print(x)
```

</details>

<a id="diagram"></a>

**Діаграма 1.2. Потік.**

```mermaid
flowchart TB
    accTitle: Потік
    accDescr: Перехід між двома кроками.
    A --> B
```

**Схема 1.3. Стани.**

```mermaid
stateDiagram-v2
    accTitle: Стани
    accDescr: Один навчальний перехід.
    [*] --> Ready
```

## Висновки

Підсумок.
"""


class BookCheckerTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(dir=ROOT / "review/tmp")
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.chapter = self.root / "ch01-example.md"

    def codes(self, text):
        findings, _ = check_book.scan_document(self.chapter, text)
        return {finding.code for finding in findings}

    def test_valid_chapter_passes_and_code_is_parsed(self):
        findings, blocks = check_book.scan_document(self.chapter, VALID_CHAPTER)
        syntax, checked, skipped = check_book.run_syntax(blocks)
        self.assertEqual(findings, [])
        self.assertEqual(syntax, [])
        self.assertEqual(checked["Python"], 1)
        self.assertEqual(skipped, {})

    def test_reports_collapsing_version_count_and_duplicate_anchors(self):
        broken = (
            VALID_CHAPTER.replace("<details>", "<details open>")
            .replace("2 рядки, Python 3.14.4", "3 рядки, Python")
            .replace('<a id="diagram"></a>', '<a id="example"></a>')
        )
        self.assertTrue({
            "LISTING_OPEN", "LISTING_LINES", "LISTING_VERSION",
            "LISTING_NUMBER", "ANCHOR_DUPLICATE",
        } <= self.codes(broken))

    def test_reports_unlabelled_and_unwrapped_program(self):
        broken = VALID_CHAPTER.replace("```python", "```")
        broken = broken.replace("<details>\n", "").replace("</details>\n", "")
        self.assertTrue({"LISTING_LANGUAGE", "LISTING_COLLAPSE"} <= self.codes(broken))

    def test_reports_wrong_headings_toc_and_numbering(self):
        broken = VALID_CHAPTER.replace("### 1.1. Приклад", "### 1.2. Інше")
        broken = broken.replace("Діаграма 1.2.", "Діаграма 1.4.")
        self.assertTrue({
            "HEADING_SUBSECTION", "TOC_ENTRY", "TOC_STALE", "OBJECT_NUMBER",
        } <= self.codes(broken))

    def test_reports_diagram_accessibility_and_unclosed_fence(self):
        broken = VALID_CHAPTER.replace("    accDescr: Перехід між двома кроками.\n", "")
        broken = broken.rsplit("```", 1)[0]
        self.assertIn("FENCE_UNCLOSED", self.codes(broken))
        findings, _ = check_book.scan_document(
            self.chapter, VALID_CHAPTER.replace("    accDescr: Перехід між двома кроками.\n", "")
        )
        self.assertIn("DIAGRAM_ACCESSIBILITY", {finding.code for finding in findings})

    def test_inline_html_examples_are_not_treated_as_open_tags_or_anchors(self):
        text = VALID_CHAPTER.replace(
            "Опис.", 'Синтаксис: `<details>` і `<a id="example"></a>`.',
        )
        self.assertNotIn("DETAILS_UNCLOSED", self.codes(text))
        self.assertNotIn("ANCHOR_DUPLICATE", self.codes(text))

    def test_math_inspector_ignores_literal_macro_but_rejects_it_in_math(self):
        source = self.root / "math.md"
        source.write_text(
            "Текст: `\\operatorname{literal}`.\n\n"
            "$\\operatorname{a}$\n\n"
            "$`\\operatorname{b}`$\n\n"
            "$$\n\\operatorname{c}\n$$\n\n"
            "```math\n\\operatorname{d}\n```\n\n"
            "```python\nr'\\operatorname{literal}'\n```\n",
            encoding="utf-8",
        )
        errors, _ = inspect_math_formatting.inspect_file(str(source))
        lines = [line for line, code, _ in errors if code == "ERR-MATH-11"]
        self.assertEqual(lines, [3, 5, 8, 12])

    def test_line_count_includes_final_blank_code_line(self):
        text = VALID_CHAPTER.replace("print(x)\n```", "print(x)\n\n```")
        self.assertIn("LISTING_LINES", self.codes(text))
        self.assertNotIn("LISTING_LINES", self.codes(text.replace("2 рядки", "3 рядки")))

    def test_existing_dash_caption_with_inline_explanation_is_counted(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**", "**Діаграма 1.2 — Потік.** Пояснення."
        )
        self.assertNotIn("OBJECT_MISSING", self.codes(text))
        self.assertNotIn("OBJECT_NUMBER", self.codes(text))

    def test_listing_caption_accepts_dash_and_optional_final_period(self):
        text = VALID_CHAPTER.replace(
            "**Лістинг 1.1. Приклад (Python 3.14.4).**",
            "**Лістинг 1.1 — Приклад (Python 3.14.4)**",
        )
        self.assertNotIn("LISTING_CAPTION", self.codes(text))

    def test_numbered_equation_caption_may_follow_formula(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "$$\nx=1\n$$\n\n**Рівняння 1.2 — Приклад**\n\n**Діаграма 1.3. Потік.**",
        ).replace("**Схема 1.3. Стани.**", "**Схема 1.4. Стани.**")
        self.assertNotIn("MATH_CAPTION", self.codes(text))
        self.assertNotIn("OBJECT_NUMBER", self.codes(text))

    def test_math_caption_can_follow_symbol_explanation(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "$$\nx=1\n$$\n\nТут x — навчальне значення.\n\n"
            "**Формула 1.2. Значення.**\n\n**Діаграма 1.3. Потік.**",
        ).replace("**Схема 1.3. Стани.**", "**Схема 1.4. Стани.**")
        self.assertNotIn("MATH_CAPTION", self.codes(text))
        self.assertNotIn("OBJECT_NUMBER", self.codes(text))

    def test_inline_math_caption_and_fenced_math_are_valid(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "Оцінку задано так. **Рівняння 1.2 — Межа**\n\n"
            "```math\nx=1\n```\n\n**Діаграма 1.3. Потік.**",
        ).replace("**Схема 1.3. Стани.**", "**Схема 1.4. Стани.**")
        self.assertNotIn("MATH_CAPTION", self.codes(text))
        self.assertNotIn("OBJECT_NUMBER", self.codes(text))
        unlabelled = text.replace("**Рівняння 1.2 — Межа**", "Ненумерована межа.")
        self.assertIn("MATH_CAPTION", self.codes(unlabelled))

    def test_caption_before_next_formula_does_not_label_previous_formula(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "$$\nx=0\n$$\n\nПопередню формулу не підписано.\n\n"
            "**Формула 1.2. Наступне значення.**\n\n$$\nx=1\n$$\n\n"
            "**Діаграма 1.3. Потік.**",
        ).replace("**Схема 1.3. Стани.**", "**Схема 1.4. Стани.**")
        findings, _ = check_book.scan_document(self.chapter, text)
        self.assertEqual(sum(finding.code == "MATH_CAPTION" for finding in findings), 1)

    def test_table_separator_accepts_single_and_multiple_columns(self):
        for separator in ("| --- |", "| :---: |", "| --- | --- |", "--- | ---",
                          "| --- | ---", "--- | --- |"):
            with self.subTest(separator=separator):
                self.assertIsNotNone(check_book.TABLE_SEPARATOR.fullmatch(separator))

    def test_one_caption_can_cover_two_adjacent_tables(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "**Таблиця 1.2 — Частини таблиці**\n\n| x |\n|---|\n| 1 |\n\n"
            "| y |\n|---|\n| 2 |\n\n**Діаграма 1.3. Потік.**",
        ).replace("**Схема 1.3. Стани.**", "**Схема 1.4. Стани.**")
        self.assertNotIn("TABLE_CAPTION", self.codes(text))
        self.assertNotIn("OBJECT_NUMBER", self.codes(text))
        self.assertIn("TABLE_CAPTION", self.codes(text.replace("**Таблиця 1.2 — Частини таблиці**", "")))

    def test_unlabelled_single_column_table_is_reported(self):
        text = VALID_CHAPTER.replace(
            "**Діаграма 1.2. Потік.**",
            "| Назва |\n| --- |\n| Значення |\n\n**Діаграма 1.2. Потік.**",
        )
        self.assertIn("TABLE_CAPTION", self.codes(text))

    def test_unlabelled_math_tables_and_diagrams_are_reported(self):
        text = VALID_CHAPTER.replace(
            "### 1.1. Приклад\n", "### 1.1. Приклад\n\n$$\nx=1\n$$\n\n| x |\n|---|\n| 1 |\n"
        )
        text = text.replace("**Діаграма 1.2. Потік.**", "Ненумерована схема.")
        self.assertTrue({"MATH_CAPTION", "TABLE_CAPTION", "DIAGRAM_CAPTION"} <= self.codes(text))

    def test_syntax_validation_never_executes_python(self):
        marker = self.root / "must-not-exist"
        source = f"from pathlib import Path\nPath({str(marker)!r}).write_text('executed')"
        block = check_book.CodeBlock(self.chapter, 1, "python", source, 2, None, False)
        findings, checked, skipped = check_book.run_syntax([block])
        self.assertEqual(findings, [])
        self.assertEqual(checked["Python"], 1)
        self.assertFalse(marker.exists())
        self.assertEqual(skipped, {})

    def test_syntax_error_and_explicit_skip_are_reported(self):
        text = VALID_CHAPTER.replace("x = 1", "x = (")
        _, blocks = check_book.scan_document(self.chapter, text)
        findings, _, _ = check_book.run_syntax(blocks)
        self.assertIn("PYTHON_SYNTAX", {finding.code for finding in findings})
        text = text.replace("```python\n", "<!-- skip-validate -->\n```python\n")
        _, blocks = check_book.scan_document(self.chapter, text)
        findings, _, skipped = check_book.run_syntax(blocks)
        self.assertEqual(findings, [])
        self.assertEqual(skipped["python: директива"], 1)

    def test_missing_tool_is_a_failure_not_a_successful_skip(self):
        block = check_book.CodeBlock(self.chapter, 1, "prolog", "p(a).", 1, None, False)
        with mock.patch.object(check_book.shutil, "which", return_value=None):
            findings, checked, skipped = check_book.run_syntax([block, block])
        self.assertEqual([finding.code for finding in findings], ["TOOL_MISSING"])
        self.assertEqual(checked, {})
        self.assertEqual(skipped["prolog: немає swipl"], 2)

    def test_inspector_error_on_stderr_is_not_reported_as_success(self):
        result = subprocess.CompletedProcess(["inspector"], 0, stdout="", stderr="[ERROR] Cannot read file\n")
        with mock.patch.object(check_book.subprocess, "run", return_value=result):
            with contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
                self.assertFalse(check_book.run_inspector("Тест", ["inspector"]))

    @unittest.skipUnless(shutil.which("swipl"), "SWI-Prolog is unavailable")
    def test_prolog_is_parsed_without_executing_directives(self):
        marker = self.root / "must-not-exist"
        source = f":- open('{marker}', write, _).\np(a)."
        valid = check_book.CodeBlock(self.chapter, 1, "prolog", source, 2, None, False)
        invalid = check_book.CodeBlock(self.chapter, 10, "prolog", "p(.", 1, None, False)
        findings, checked, _ = check_book.run_syntax([valid, invalid])
        self.assertEqual(checked["prolog"], 1)
        self.assertEqual([finding.code for finding in findings], ["CODE_SYNTAX"])
        self.assertFalse(marker.exists())

    @unittest.skipUnless(shutil.which("gofmt"), "gofmt is unavailable")
    def test_go_wrappers_accept_statements_and_reject_invalid_syntax(self):
        good = check_book.CodeBlock(self.chapter, 1, "go", "x := 1\n_ = x", 2, None, False)
        bad = check_book.CodeBlock(self.chapter, 4, "go", "if ( {", 1, None, False)
        findings, checked, _ = check_book.run_syntax([good, bad])
        self.assertEqual(len(check_book.go_candidates("func f() { ... }")), 6)
        self.assertEqual(checked["go"], 1)
        self.assertEqual([finding.code for finding in findings], ["CODE_SYNTAX"])

    def test_volume_scope_excludes_other_volumes_and_translations(self):
        book = self.root / "Expert_Systems"
        for directory in ("Volume-1", "Volume-2", "Volume-1/translations"):
            (book / directory).mkdir(parents=True)
            (book / directory / "sample.md").write_text("# Тест\n", encoding="utf-8")
        with mock.patch.object(check_book, "ROOT", self.root), mock.patch.object(check_book, "BOOK", book):
            scope, files = check_book.discover_documents(1)
        self.assertEqual(scope, book / "Volume-1")
        self.assertEqual(files, [book / "Volume-1/sample.md"])

    def test_volume_scope_rejects_symbolic_documents(self):
        book = self.root / "Expert_Systems"
        (book / "Volume-1").mkdir(parents=True)
        document = book / "Volume-1/sample.md"
        document.write_text("# Тест\n", encoding="utf-8")
        (book / "Volume-1/link.md").symlink_to(document)
        with mock.patch.object(check_book, "ROOT", self.root), mock.patch.object(check_book, "BOOK", book):
            with self.assertRaisesRegex(ValueError, "Символічний документ"):
                check_book.discover_documents(1)

    def test_terminology_inspector_accepts_multiple_explicit_files(self):
        good = self.root / "good.md"
        bad = self.root / "bad.md"
        good.write_text("Система безпеки.\n", encoding="utf-8")
        bad.write_text("контур безпеки\n", encoding="utf-8")
        command = [sys.executable, str(ROOT / "scripts/inspect_buzzwords.py")]
        selected = subprocess.run(command + [str(good)], text=True, capture_output=True)
        both = subprocess.run(command + [str(good), str(bad)], text=True, capture_output=True)
        missing = subprocess.run(command + [str(self.root / "missing.md")], text=True, capture_output=True)
        self.assertEqual(selected.returncode, 0, selected.stdout + selected.stderr)
        self.assertEqual(both.returncode, 1, both.stdout + both.stderr)
        self.assertIn("bad.md", both.stdout)
        self.assertEqual(missing.returncode, 1, missing.stdout + missing.stderr)
        self.assertIn("Could not read", missing.stdout)

    @unittest.skipUnless(
        shutil.which("node") and (ROOT / "review/tmp/technical-tools/node_modules").is_dir(),
        "Node.js render-validation dependencies are unavailable",
    )
    def test_volume_command_returns_success_and_failure_for_synthetic_chapters(self):
        book = self.root / "Expert_Systems"
        volume = book / "Volume-1"
        volume.mkdir(parents=True)
        document = volume / "ch01-example.md"
        document.write_text(VALID_CHAPTER, encoding="utf-8")
        with mock.patch.object(check_book, "BOOK", book):
            output = io.StringIO()
            with contextlib.redirect_stdout(output):
                result = check_book.main(["--volume", "1"])
            self.assertEqual(result, 0, output.getvalue())

            document.write_text(VALID_CHAPTER.replace("x = 1", "x = ("), encoding="utf-8")
            output = io.StringIO()
            with contextlib.redirect_stdout(output):
                result = check_book.main(["--volume", "1"])
            self.assertEqual(result, 1, output.getvalue())
            self.assertIn("PYTHON_SYNTAX", output.getvalue())


@unittest.skipUnless(
    shutil.which("node") and (ROOT / "review/tmp/technical-tools/node_modules").is_dir(),
    "Node.js render-validation dependencies are unavailable",
)
class RenderScopeTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(dir=ROOT / "review/tmp")
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.selected = self.root / "selected.md"
        self.linked = self.root / "linked.md"
        self.report = self.root / "render.json"
        self.selected.write_text(
            "# Selected\n\n[Інший документ](linked.md#other)\n\n$x$\n\n"
            "```mermaid\nflowchart TB\nA --> B\n```\n", encoding="utf-8"
        )
        self.linked.write_text(
            "# Other\n\n[Застаріле посилання](missing.md)\n\n```python\nx = 1\n```\n",
            encoding="utf-8",
        )

    def render(self):
        return subprocess.run(
            ["node", str(ROOT / "scripts/inspect_render_syntax.mjs"),
             "--workspace-only", f"--report={self.report}", str(self.selected)],
            cwd=ROOT, text=True, capture_output=True, timeout=120,
        )

    def test_linked_document_is_used_for_anchors_but_not_added_to_scope(self):
        result = self.render()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertEqual(data["summary"]["links"]["total"], 1)
        self.assertEqual(len(data["code_blocks"]), 1)
        self.assertEqual(data["scope"], [str(self.selected.relative_to(ROOT))])

    def test_missing_anchor_in_linked_document_is_still_a_failure(self):
        self.linked.write_text("# Different\n", encoding="utf-8")
        result = self.render()
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertEqual(data["summary"]["links"]["failed"], 1)

    def test_outside_workspace_target_is_reported_without_reading_it(self):
        target = ROOT.parent / "not-accessed.md"
        self.selected.write_text(f"# Тест\n\n[Посилання]({target}#section)\n", encoding="utf-8")
        result = self.render()
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertIn("outside the workspace", data["links"][0]["error"])

    def test_symbolic_target_is_reported_without_following_it(self):
        (self.root / "link.md").symlink_to(self.linked)
        self.selected.write_text("# Тест\n\n[Посилання](link.md#other)\n", encoding="utf-8")
        result = self.render()
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertIn("Symbolic local target", data["links"][0]["error"])

    def test_restricted_dot_tmp_target_is_not_accessed(self):
        self.selected.write_text("# Тест\n\n[Посилання](.tmp/hidden.md#section)\n", encoding="utf-8")
        result = self.render()
        self.assertEqual(result.returncode, 1, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertIn("Restricted .tmp", data["links"][0]["error"])

    def test_documents_without_math_or_diagrams_are_valid_input(self):
        self.selected.write_text("# Текст без формул і діаграм\n", encoding="utf-8")
        result = self.render()
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)
        data = json.loads(self.report.read_text(encoding="utf-8"))
        self.assertEqual(data["summary"]["math"]["total"], 0)
        self.assertEqual(data["summary"]["mermaid"]["total"], 0)


if __name__ == "__main__":
    unittest.main()
