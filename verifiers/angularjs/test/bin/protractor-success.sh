#!/bin/bash

conf=$1

if [[ -z "$conf" ]]; then
  >&2 echo "No configuration file given!"
  exit 127
fi

if [[ ! -f "$conf" ]]; then
    >&2 echo "File not found!"
    exit 127
fi

read -d '' results << EOF
[
  {
    "description": "should alias index positions",
    "assertions": [
      {
        "passed": true
      },
      {
        "passed": true
      },
      {
        "passed": true
      },
      {
        "passed": true
      }
    ],
    "duration": 2020
  }
]
EOF

echo "[launcher] Running 1 instances of WebDriver"
sleep 0.5
echo "."
echo ""
echo "Finished in 0.552 seconds"
echo "1 test, 4 assertions, 0 failures"
echo "[launcher] 0 instance(s) of WebDriver still running"
echo "[launcher] chrome #1 passed"

echo $results > `dirname $conf`/results.json