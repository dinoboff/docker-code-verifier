#!/usr/bin/env python3
import argparse
import base64

from codeverifier import TestRunner


parser = argparse.ArgumentParser(
    description='Run some python code and test it with some doctest tests.'
)
parser.add_argument(
    '-e',
    '--encoded',
    action='store_true',
    help="the input will be base64 encoded"
)
parser.add_argument('--tests', help="The doctest test")
parser.add_argument('solution', help="the user solution")


def parse():
    args = parser.parse_args()
    if args.encoded:
        solution = base64.b64decode(args.solution).decode('utf8')
        tests = base64.b64decode(args.tests).decode('utf8')
    else:
        solution = args.solution
        tests = args.tests

    return solution, tests


def main():
    solution, tests = parse()
    runner = TestRunner(solution, tests)
    runner.run()
    print(runner.to_json())

if __name__ == '__main__':
    main()
