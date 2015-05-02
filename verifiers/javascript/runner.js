'use strict';

var vm = require('vm');
var assert = require('assert');

/**
 * Run tests in a sandbox.
 *
 * usage:
 *
 *   testSolution('var foo = 1;', 'test('foo should be set', function(){ assert.equal(1, foo)})')
 *   // would return `{solved: true, results:[{test: 'foo should be set', correct: true}]}`
 *
 *
 * @param  String solution User provided JS solution
 * @param  String tests    User (admin) provided JS tests
 * @return Object          holding `solved` (bool), `results` and `errors`
 *
 */
var testSolution = function(solution, tests) {
  var ctx = {};

  try {
    runSolution(solution, ctx);
    initTests(tests, ctx);
    return runTests(ctx);
  } catch (e) {
    return Promise.resolve({
      solved: false,
      errors: e.toString()
    });
  }
};

/**
 * Run the user provided solution inside a sandbox.
 *
 * `solution` is the the user solution (string) to evaluate. `ctx` will be
 * the global object the solution will have access to.
 *
 * Will throw if the solution is undefined or if the solution defines
 * `assert`, `test` or `__tests__` properties on the global object.
 *
 */
function runSolution(solution, ctx) {
  if (solution == null) {
    throw new Error('Solutions are missings');
  }

  if (ctx == null) {
    ctx = {};
  }

  try {
    vm.runInNewContext(solution, ctx);
  } catch (e) {
    throw new Error('Failed to run solutions: ' + e);
  }

  if (ctx.test || ctx.__tests__ || ctx.assert) {
    throw new Error('"assert", test" and "__tests__" cannot be defined in a solution');
  }
}

/**
 * Run the user provided tests inside a sandbox.
 *
 * `tests` is the the user tests (string) to evaluate. `ctx` will be
 * the global object the solution will have access to.
 *
 * It will hold the properties defined by the user solution, `assert` and
 * `test`.
 *
 */
function initTests(tests, ctx) {
  if (tests == null) {
    throw new Error('Tests are missings.');
  }

  if (ctx == null) {
    throw new Error('Context is missing.');
  }

  ctx.assert = assert;
  ctx.__tests__ = [];
  ctx.test = function(title, cb) {
    ctx.__tests__.push({
      test: title,
      cb: cb
    });
  };

  try {
    vm.runInNewContext(tests, ctx);
  } catch (e) {
    throw new Error('Failed to initiate tests');
  }
}

/**
 * Run all tests defined by the user provided tests asynchronously.
 *
 * Returns a promise resolving to a response object.
 *
 */
function runTests(ctx) {
  return Promise.all(ctx.__tests__.map(function(test) {
    var onSuccess = function() {
      return {test: test.test, correct: true};
    };
    var onError = function(e) {
      return {test: test.test, correct: false, error: e.toString()};
    };

    try {
      return Promise.resolve(test.cb()).then(onSuccess).catch(onError);
    } catch (e) {
      return onError(e);
    }
  })).then(function(results) {
    return {
      results: results,
      solved: results.every(function(r) {
        return r.correct === true;
      })
    };
  });
}

module.exports = {
  testSolution: testSolution,
  runSolution: runSolution,
  initTests: initTests,
  runTests: runTests
};
