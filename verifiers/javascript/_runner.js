'use strict';

var minimist = require('minimist');
var runner = require('./runner');

function runTests(solution, tests) {
  var ctx = {};

  try {
    runner.runSolution(solution, ctx);
    runner.initTests(tests, ctx);
    return runner.runTests(ctx);
  } catch (e) {
    return Promise.resolve({
      solved: false,
      errors: e.toString()
    });
  }
}

function main(argv) {
  /* eslint no-console: 0 */
  runTests(argv.solution, argv.tests).then(function(resp) {
    console.log(JSON.stringify(resp));
  });
}

main(minimist(process.argv));
