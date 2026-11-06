import runpy
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
INSPECT = runpy.run_path(str(ROOT / "scripts/inspect_math_formatting.py"))["inspect_file"]


class DegreeNotationTests(unittest.TestCase):
    def codes(self, text):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "sample.md"
            path.write_text(text, encoding="utf-8")
            return {code for _, code, _ in INSPECT(str(path))[0]}

    def test_spacing_before_degree_is_rejected(self):
        for expression in (r"85\,^\circ\mathrm C", r"80\;^{\circ}\mathrm{C}",
                           r"2\quad^\circ\mathrm{C}", r"^\circ\mathrm C/\mathrm s",
                           r"(^\circ\mathrm C)^2"):
            for delimiters in ("plain", "protected", "block", "fence"):
                with self.subTest(expression=expression, delimiters=delimiters):
                    text = {
                        "plain": f"${expression}$\n",
                        "protected": f"$`{expression}`$\n",
                        "block": f"$$\n{expression}\n$$\n",
                        "fence": f"```math\n{expression}\n```\n",
                    }[delimiters]
                    self.assertIn("ERR-MATH-12", self.codes(text))

    def test_explicit_base_and_angle_superscripts_are_accepted(self):
        for expression in (r"85\,{}^\circ\mathrm{C}", r"{}^\circ\mathrm{C}/\mathrm{s}",
                           r"90^\circ", r"x^\circ"):
            with self.subTest(expression=expression):
                self.assertEqual(self.codes(f"$`{expression}`$\n"), set())

    def test_program_code_and_inline_code_are_ignored(self):
        text = "```python\n" + r'temperature = r"85\,^\circ\mathrm C"' + "\n```\n"
        text += r"Example code: `85\,^\circ\mathrm C`." + "\n"
        self.assertEqual(self.codes(text), set())

    def test_error_reports_source_line(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / "sample.md"
            path.write_text("Heading\n\n" + r"$`85\,^\circ\mathrm C`$" + "\n",
                            encoding="utf-8")
            findings, modified = INSPECT(str(path))
            self.assertEqual([(line, code) for line, code, _ in findings],
                             [(3, "ERR-MATH-12")])
            self.assertFalse(modified)


if __name__ == "__main__":
    unittest.main()
