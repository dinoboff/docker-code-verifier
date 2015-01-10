#!/bin/bash

conf=$1

if [[ -z "$conf" ]]; then
  >&2 echo "No configuration file given."
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
        "passed": false,
        "errorMsg": "Expected 'list[ 1 ][ 1 ] = d;' to be 'list[ 1 ][ 1 ] = e;'.",
        "stackTrace": "Error: Failed expectation\n    at [object Object].<anonymous> (/www/test.js:8:39)\n    at /usr/local/lib/node_modules/protractor/node_modules/jasminewd/index.js:94:14\n    at [object Object].webdriver.promise.ControlFlow.runInNewFrame_ (/usr/local/lib/node_modules/protractor/node_modules/selenium-webdriver/lib/webdriver/promise.js:1654:20)\n    at [object Object].webdriver.promise.ControlFlow.runEventLoop_ (/usr/local/lib/node_modules/protractor/node_modules/selenium-webdriver/lib/webdriver/promise.js:1518:8)\n    at [object Object].wrapper [as _onTimeout] (timers.js:261:14)\n    at Timer.listOnTimeout [as ontimeout] (timers.js:112:15)"
      }
    ],
    "duration": 1982
  }
]
EOF

echo "[launcher] Running 1 instances of WebDriver"
sleep 0.5
echo "F"
echo ""
echo "Failures:"
echo ""
echo "  1) test should alias index positions"
echo "   Message:"
echo "     Expected 'list[ 1 ][ 1 ] = d;' to be 'list[ 1 ][ 1 ] = e;'."
echo "   Stacktrace:"
echo "     Error: Failed expectation"
echo "    at [object Object].<anonymous> (/www/foo/test.js:8:39)"
echo ""
echo "Finished in 0.588 seconds"
echo "1 test, 4 assertions, 1 failure"
echo ""
echo "[launcher] 0 instance(s) of WebDriver still running"
echo "[launcher] chrome #1 failed 1 test(s)"
echo "[launcher] overall: 1 failed spec(s)"
echo "[launcher] Process exited with error code 1"

echo $results > `dirname $conf`/results.json

exit 1