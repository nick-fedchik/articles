import runpy
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
MATH = runpy.run_path(str(ROOT / "scripts/inspect_math_formatting.py"))
STRUCTURE = runpy.run_path(str(ROOT / "scripts/inspect_structural_hierarchy.py"))
TRACE = runpy.run_path(str(ROOT / "scripts/inspect_requirements_trace.py"))


class TechnicalInspectorTests(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(dir=ROOT / "review/tmp")
        self.addCleanup(self.directory.cleanup)
        self.path = Path(self.directory.name) / "sample.md"

    def math_codes(self, text):
        self.path.write_text(text)
        return {code for _, code, _ in MATH["inspect_file"](str(self.path))[0]}

    def structure_codes(self, text):
        self.path.write_text(text)
        return {code for _, code, _ in STRUCTURE["inspect_file"](self.path)}

    def test_literal_block_delimiters_are_code(self):
        self.assertEqual(self.math_codes("Формулу записано як `$$...$$`.\n"), set())

    def test_inverted_inline_math_is_still_rejected(self):
        self.assertIn("ERR-MATH-01", self.math_codes("Формула `$x_1$`.\n"))

    def test_protected_inline_math_is_accepted(self):
        self.assertEqual(self.math_codes("Формула $`x_1`$.\n"), set())

    def test_corrupted_macro_is_still_rejected(self):
        self.assertIn("ERR-MATH-08", self.math_codes("Формула $" + chr(9) + "op$.\n"))

    def test_catalogue_title_is_not_authorial_declaration(self):
        text = "| Назва | Опис |\n|---|---|\n| **Закон перевірки М. Федчика** | Посилання |\n"
        self.assertEqual(self.structure_codes(text), set())

    def test_real_declaration_still_needs_note(self):
        self.assertIn("ERR-STRUCT-02", self.structure_codes("**Закон перевірки М. Федчика**: опис.\n"))
        self.assertEqual(self.structure_codes("> [!NOTE]\n> **Закон перевірки М. Федчика**\n> 1. Опис.\n"), set())

    def test_source_reference_keeps_section_after_anchor(self):
        source = '[Джерело](../../source.md#part-11), §11 «Джерело».'
        self.assertEqual(list(TRACE["source_references"](source)),
                 [("../../source.md", [("11", "Джерело")])])

    def test_sections_belong_to_their_own_source(self):
        source = '[Глава](chapter.md), §3 «Метод»; [Діалог](dialogue.md#part), §7 «Межі».'
        self.assertEqual(list(TRACE["source_references"](source)),
                         [("chapter.md", [("3", "Метод")]), ("dialogue.md", [("7", "Межі")])])

    def test_repeated_source_links_keep_distinct_sections(self):
        source = '[Глава](chapter.md#one), §1 «Початок»; [Глава](chapter.md#two), §2 «Продовження».'
        self.assertEqual(list(TRACE["source_references"](source)),
                         [("chapter.md", [("1", "Початок")]), ("chapter.md", [("2", "Продовження")])])


if __name__ == "__main__":
    unittest.main()