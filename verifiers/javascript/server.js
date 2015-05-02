'use strict';
var express = require('express');
var bodyParser = require('body-parser');
var app = express();
var testSolution = require('./runner.js').testSolution;

var argv = require('minimist')(process.argv);
var options = {
  port: argv.port || 5000,
  host: argv.host || '0.0.0.0',
  root: argv.root || '/javascript'
};

var errJsonrequestRequired = 'jsonrequest is a required property.';
var errJsonrequestRequiredBase64 = 'jsonrequest is required property. ' +
  'It should be JSON object, ' +
  'base64 encoded or plain.';
var errNoPayload = 'The payload is missing. It should be a JSON object with' +
  ' "solution" and "tests" attributes.';
var errJsonEncoding = 'The payload is not a valid JSON object.';
var errSolutionRequired = 'The request did not include a user solution.';
var errTestRequired = 'The request did not include a test.';
var codeBadRequest = 400;
var codeInternalError = 500;

// Middlewares

// Body parsers
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
  extended: true
}));

// Allow cross domain requests
app.use(options.root, function setAccessControl(req, res, next) {
  res.set('Access-Control-Allow-Origin', '*');
  next();
});

// Routes

// Handle GET request to the Verifier.
//
// For this verifier, the payload should include both a "solutions" and "tests"
// attribute.
//
app.get(options.root, function verifierGetHandler(req, res) {

  var jsonrequest = req.query.jsonrequest;
  var cb = req.query.vcallback;

  if (!jsonrequest) {
    res.status(codeBadRequest).json({
      errors: errJsonrequestRequired
    });
    return;
  }

  if (jsonrequest[0] !== '{') {
    jsonrequest = new Buffer(jsonrequest, 'base64').toString('utf8');
  }

  try {
    jsonrequest = JSON.parse(jsonrequest);
  } catch (e) {
    res.status(codeBadRequest).json({
      errors: errJsonrequestRequiredBase64
    });
    return;
  }

  if (!jsonrequest.solution) {
    res.status(codeBadRequest).json({
      errors: errSolutionRequired
    });
    return;
  }

  if (!jsonrequest.tests) {
    res.status(codeBadRequest).json({
      errors: errTestRequired
    });
    return;
  }

  testSolution(jsonrequest.solution, jsonrequest.tests, options).then(function(results) {
    if (!cb) {
      res.json(results);
      return;
    }

    var body = cb + '(' + JSON.stringify(results) + ')';
    res.type('application/javascript').send(body);

  }).catch(function(err) {
    if (typeof err === 'object') {
      res.status(codeInternalError).json(err);
    } else {
      res.status(codeInternalError).json({
        solved: false,
        errors: err.toString()
      });
    }
  });

});

// Handle POST request to the Verifier.
//
// Test a html document against protractor scenarios.
//
// It should find the payload in the body, which can be JSON or form data
// encoded. In that case, the payload should be a json object found in the
// "jsonrequest parameter."
//
// For this verifier, the payload should include both a "solutions" and "tests"
// attribute.
//
app.post(options.root, function verifierPostHandler(req, res) {
  var jsonrequest;

  if (!req.body) {
    res.status(codeBadRequest).json({
      errors: errNoPayload
    });
    return;
  }

  var payload = req.body;
  if (payload.jsonrequest) {
    payload = payload.jsonrequest;
    try {
      jsonrequest = JSON.parse(payload);
    } catch (e) {
      res.status(codeBadRequest).json({
        solved: false,
        errors: errJsonEncoding
      });
      return;
    }
  } else {
    jsonrequest = payload;
  }

  if (!jsonrequest.solution) {
    res.status(codeBadRequest).json({
      errors: errSolutionRequired
    });
    return;
  }

  if (!jsonrequest.tests) {
    res.status(codeBadRequest).json({
      errors: errTestRequired
    });
    return;
  }

  testSolution(jsonrequest.solution, jsonrequest.tests).then(function(results) {
    res.json(results);
  }).catch(function(err) {
    if (typeof err === 'object') {
      res.status(codeInternalError).json(err);
    } else {
      res.status(codeInternalError).json({
        solved: false,
        errors: err.toString()
      });
    }
  });
});

app.listen(options.port, options.host, function serverHandler() {
  /* eslint no-console: 0*/
  console.log('Binding server to', options.host + ':' + options.port);
});
