#!/usr/bin/env python3
#
# CLI to run user solution and tests
#
import json
import sys

from codeverifier import TestRunner


def main():
    args = sys.argv[-2:]
    if len(args) == 1:
        solution, tests = args[0], ""
    elif len(args) == 2:
        solution, tests = args
    else:
        exit(1)

    runner = TestRunner(solution, tests)
    runner.run()
    json.dump(runner.to_dict(), sys.stdout, indent=2)


if __name__ == '__main__':
    main()
