exports.config = {
  seleniumAddress: 'http://selenium:4444/wd/hub',
  specs: ['/www/_protractor/example/specs.js'],
  baseUrl: 'http://static/_protractor/example/index.html',
  capabilities: {
    'browserName': 'phantomjs'
  },
  framework: 'jasmine',
  resultJsonOutputFile: '/www/_protractor/example/results.js',
  jasmineNodeOpts: {
    showColors: false
  },
  allScriptsTimeout: 5000,
  getPageTimeout: 2500
};
