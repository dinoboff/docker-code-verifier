var spawn = require('child_process').spawn;
var path = require('path');
var q = require('q');
var fs = require('fs');
var writeFile = q.nfbind(fs.writeFile);
var readFile = q.nfbind(fs.readFile);
var deleteFile = q.nfbind(fs.unlink);
var chmod = q.nfbind(fs.chmod);
var rmdir = q.nfbind(fs.rmdir);
var tmp = require('tmp');


function deleteFileIfExist(path) {
  'use strict';

  return deleteFile(path).catch(function(err) {
    if (!err.code || err.code !== 'ENOENT') {
      throw err;
    }
  });
}


function exists(path) {
  'use strict';

  var f = q.defer();

  fs.exists(path, function(found) {
    if (found) {
      f.resolve(path);
    } else {
      f.reject(new Error(path + ' was not found.'));
    }
  });

  return f.promise;
}


// Create a temporary directory inside the static directory.
//
// Returns a promise resolving to the directory path.
//
// Note that the directory need to be removed manually. There's no automatic
// cleanup.
//
var fTempDir = function(staticDir) {
  'use strict';

  var f = q.defer();

  tmp.dir({
    mode: '0777',
    dir: staticDir
  }, function(err, path) {
    if (err) {
      f.reject(err);
    } else {
      f.resolve(path);
    }
  });

  return f.promise.then(function(path){
    return chmod(path, '0777').then(function(){
      return path;
    });
  });
};


// Writes the html document to test, the scenario to run and protractor
// config in files in a temporary directory.
//
// Returns a promise resolving to an object with the path to the
// temporary directory and the testing files.
//
var setUpTest = function(solution, tests, options) {
  'use strict';

  return fTempDir(options.staticDir).then(function(dirPath) {
    var paths = {
      dir: dirPath,
      specs: path.join(dirPath, 'specs.js'),
      html: path.join(dirPath, 'index.html'),
      config: path.join(dirPath, 'config.js'),
      results: path.join(dirPath, 'results.json')
    };

    var config = {
      seleniumAddress: options.seleniumUrl,
      specs: [paths.specs],
      baseUrl: paths.html.replace(options.staticPrefixRegex, options.staticUrl),
      framework: 'jasmine',
      resultJsonOutputFile: paths.results,
      jasmineNodeOpts: {
        showColors: false
      }
    };

    var fSpec = writeFile(paths.specs, tests);
    var fHtml = writeFile(paths.html, solution);
    var fConfig = writeFile(paths.config, 'exports.config = ' + JSON.stringify(config)) + ';';

    return q.all([fSpec, fHtml, fConfig]).then(function() {
      return paths;
    });
  });
};


// Cleanup the temporary folder holding the testing files.
//
var cleanUpTests = function(paths) {
  'use strict';

  return q.all([
    deleteFileIfExist(paths.specs),
    deleteFileIfExist(paths.html),
    deleteFileIfExist(paths.config),
    deleteFileIfExist(paths.results)
  ]).then(function() {
    return rmdir(paths.dir);
  });
};


// Test a html document against protractor scenarios.
//
// For this verifier, both solution and tests should be provided.
//
var testSolution = function(solution, tests, options) {
  'use strict';

  return setUpTest(solution, tests, options).then(function(paths) {
    var f = q.defer();
    var cleaningUp = null;
    var results = {
      solved: false,
      results: null,
    };

    var protractor = spawn(options.protractorPath, [paths.config], {
      uid: options.uid
    });

    var timeout = setTimeout(function() {
      f.resolve({
        solved: false,
        errors: 'timeout'
      });

      protractor.kill('SIGTERM');
      if (cleaningUp === null) {
        cleaningUp = cleanUpTests(paths);
      }
      timeout = null;
    }, options.timeoutDelay);

    var cleanUp = function() {
      if (timeout !== null) {
        clearTimeout(timeout);
        timeout = null;
      }

      if (cleaningUp === null) {
        cleaningUp = cleanUpTests(paths);
      }
    };

    protractor.on('exit', function(code) {
      if (timeout !== null) {
        clearTimeout(timeout);
        timeout = null;
      }

      results.solved = code === 0;

      exists(paths.results).then(function() {
        return readFile(paths.results);
      }, function() {
        throw new Error('Protractor results are missing at ' + paths.results);
      }).then(function(content) {
        try {
          results.results = JSON.parse(content);
        } catch (e) {
          throw new Error('Failed to parse protractor results (' + e + ')');
        }
        f.resolve(results);
      }).catch(function(e){
        f.reject({errors: e.toString()});
      }).finally(cleanUp);

    });

    protractor.on('error', function(err) {
      f.reject(err);
      cleanUp();
    });

    return f.promise;
  });
};


module.exports = {
  testSolution: testSolution,
  setUpTest: setUpTest,
  cleanUpTests: cleanUpTests,
  fTempDir: fTempDir
};
