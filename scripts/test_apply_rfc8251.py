#!/usr/bin/env python3
"""Check post-generation reproducibility, idempotence and fail-closed writes."""
import importlib.util
import pathlib
import subprocess
import sys
import tempfile
import unittest

root = pathlib.Path(__file__).resolve().parents[1]
script = root / 'scripts/apply_rfc8251.py'
spec = importlib.util.spec_from_file_location('postprocess', script)
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)


class PostGenerationTests(unittest.TestCase):
    def setUp(self):
        self.patched = (root / 'lib.go').read_text()
        self.original = self.patched
        for old, new in [(module.old_copy, module.new_copy), (module.old_condition, module.new_condition), (module.old_end, module.new_end)]:
            self.assertEqual(self.original.count(new), 1)
            self.original = self.original.replace(new, old, 1)

    def invoke(self, contents, check=False):
        with tempfile.TemporaryDirectory() as directory:
            path = pathlib.Path(directory) / 'lib.go'
            path.write_text(contents)
            cmd = [sys.executable, str(script), '--file', str(path)]
            if check:
                cmd.append('--check')
            result = subprocess.run(cmd, capture_output=True, text=True)
            return result.returncode, path.read_text()

    def test_reapplication_matches_maintained_file(self):
        self.assertEqual(self.invoke(self.original), (0, self.patched))

    def test_idempotent(self):
        self.assertEqual(self.invoke(self.patched), (0, self.patched))

    def test_check_does_not_change_unpatched_file(self):
        code, contents = self.invoke(self.original, check=True)
        self.assertNotEqual(code, 0)
        self.assertEqual(contents, self.original)
        self.assertEqual(self.invoke(self.patched, check=True), (0, self.patched))

    def test_partial_or_unknown_output_is_not_written(self):
        for text in [self.original.replace(module.old_condition, module.new_condition, 1), self.original.replace(module.old_end, '// changed generator output', 1)]:
            code, contents = self.invoke(text)
            self.assertNotEqual(code, 0)
            self.assertEqual(contents, text)


if __name__ == '__main__':
    unittest.main()
