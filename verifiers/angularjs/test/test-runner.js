/* global describe, it, before, after */
var q = require('q');
var runner = require('../runner.js');
var rimraf = require('rimraf');
var fs = require('fs');
var path = require('path');
var readFile = q.nfbind(fs.readFile);
var writeFile = q.nfbind(fs.writeFile);
require('chai').should();


function exists(path) {
  'use strict';

  var f = q.defer();

  fs.exists(path, function(found) {
    if (found) {
      f.resolve(true);
    } else {
      f.reject(new Error(path + ' was not found.'));
    }
  });

  return f.promise;
}


if (!String.prototype.startsWith) {
  Object.defineProperty(String.prototype, 'startsWith', {
    enumerable: false,
    configurable: false,
    writable: false,
    value: function(searchString, position) {
      'use strict';

      position = position || 0;
      return this.lastIndexOf(searchString, position) === position;
    }
  });
}


describe('runner', function() {
  'use strict';

  var staticDir = path.join(__dirname, 'static');

  before(function(done) {
    fs.mkdir(staticDir, function() {
      done();
    });
  });

  after(function(done) {
    rimraf(staticDir, done);
  });


  describe('fTempDir', function() {

    it('should create a temporary directory', function() {
      return runner.fTempDir(staticDir).then(function(path) {
        path.startsWith(staticDir).should.equal(true);
        return exists(path);
      });
    });

    it('should resolve to an error when the static dir is missing', function() {
      var badDir = path.join(__dirname, 'bad');

      return runner.fTempDir(badDir).then(function() {
        throw new Error('Unexpected');
      }, function(e) {
        e.code.should.equal('ENOENT');
      });
    });
  });


  describe('setUpTest', function() {

    it('should create a temporary dir in static directory', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(paths) {
        paths.dir.startsWith(staticDir).should.equal(true);
        return exists(paths.dir);
      });
    });

    it('should create solution file', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(paths) {
        paths.html.startsWith(paths.dir);
        return exists(paths.html).then(function() {
          return readFile(paths.html);
        });
      }).then(function(content) {
        content.toString().should.equal('solution');
      });
    });

    it('should create test file', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(paths) {
        paths.specs.startsWith(paths.dir);
        return exists(paths.specs).then(function() {
          return readFile(paths.specs);
        });
      }).then(function(content) {
        content.toString().should.equal('tests');
      });
    });

    it('should create config file', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static/_protractor'
      }).then(function(paths) {
        paths.config.startsWith(paths.dir);
        return exists(paths.config).then(function() {
          var config = require(paths.config).config;

          config.seleniumAddress.should.equal('http://selenium:4444/wd/hub');
          config.specs.length.should.equal(1);
          config.specs[0].should.equal(paths.specs);
          config.baseUrl.startsWith('http://static/_protractor');
          config.resultJsonOutputFile.should.equal(paths.results);
        });
      });
    });

  });


  describe('cleanUpTests', function() {

    it('should remove temporary files', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(paths) {
        return writeFile(paths.results, '{}').then(function() {
          return paths;
        });
      }).then(function(paths) {
        return runner.cleanUpTests(paths).then(function() {
          return exists(paths.dir).then(function() {
            throw new Error('The temporary directory should be gone');
          }, function() {
            // we expect it gone.
          });
        });
      });
    });

    it('should remove temporary files when results.json is missing', function() {
      return runner.setUpTest('solution', 'tests', {
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(paths) {
        return runner.cleanUpTests(paths).then(function() {
          return exists(paths.dir).then(function() {
            throw new Error('The temporary directory should be gone');
          }, function() {
            // we expect it gone.
          });
        });
      });
    });

  });


  describe('testSolution', function() {

    it('should run protractor and report success', function() {
      return runner.testSolution('solution', 'tests', {
        protractorPath: path.join(__dirname, 'bin/protractor-success.sh'),
        uid: process.getuid(),
        timeoutDelay: 2000,
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(results){
        results.solved.should.equal(true);
        results.results.length.should.equal(1);
      });
    });

    it('should run protractor and report failure', function() {
      return runner.testSolution('solution', 'tests', {
        protractorPath: path.join(__dirname, 'bin/protractor-failure.sh'),
        uid: process.getuid(),
        timeoutDelay: 2000,
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(results){
        results.solved.should.equal(false);
        results.results.length.should.equal(1);
      });
    });

    it('should run protractor and timeout', function() {
      return runner.testSolution('solution', 'tests', {
        protractorPath: path.join(__dirname, 'bin/protractor-success.sh'),
        uid: process.getuid(),
        timeoutDelay: 100,
        staticDir: staticDir,
        seleniumUrl: 'http://selenium:4444/wd/hub',
        staticPrefixRegex: new RegExp('^' + staticDir),
        staticUrl: 'http://static'
      }).then(function(results){
        results.solved.should.equal(false);
        results.errors.should.equal('timeout');
      });
    });

  });

});
