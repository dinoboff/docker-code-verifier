import json
import unittest

from codeverifier import TestRunner


class TestTestRunner(unittest.TestCase):

    def test_run_one_line(self):
        runner = TestRunner(
            solution='foo = 1',
            tests='>>> foo\n1'
        )
        runner.run()
        self.assertEqual(1, len(runner.results))
        self.assertEqual(
            {
                'call': 'foo',
                'expected': '1',
                'received': '1',
                'correct': True
            },
            runner.results[0]
        )
        self.assertEqual(True, runner.solved)
        self.assertEqual('', runner.printed)

    def test_run_unsolved(self):
        runner = TestRunner(
            solution='foo = 2',
            tests='>>> foo\n1'
        )
        runner.run()
        self.assertEqual(1, len(runner.results))
        self.assertEqual(
            {
                'call': 'foo',
                'expected': '1',
                'received': '2',
                'correct': False
            },
            runner.results[0]
        )
        self.assertEqual(False, runner.solved)

    def test_run_except(self):
        runner = TestRunner(
            solution='foo = bar',
            tests='>>> foo\n1'
        )
        runner.run()
        self.assertIsNone(runner.results)
        self.assertEqual("name 'bar' is not defined", runner.errors)
        self.assertEqual(False, runner.solved)

    def test_to_json_solved(self):
        runner = TestRunner(
            solution='foo = 1',
            tests='>>> foo\n1'
        )
        runner.run()
        data = json.loads(runner.to_json())
        self.assertEqual(
            {'solved', 'results', 'printed'},
            data.keys()
        )
        self.assertTrue(data['solved'])
        self.assertEqual(
            [{
                'call': 'foo',
                'expected': '1',
                'received': '1',
                'correct': True
            }],
            data['results']
        )
        self.assertEqual('', data['printed'])
