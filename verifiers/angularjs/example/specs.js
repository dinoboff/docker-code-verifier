describe('angularjs ng-controller', function() {
  it('should say "hello Bob"', function() {
    browser.get('./');

    expect(
      element(by.binding('ctrl.name')).getText()
    ).toBe('Hello Alice!');
  });

  it('should say "hello Bob"', function() {
    browser.get('./');

    expect(
      element(by.binding('ctrl.name')).getText()
    ).toBe('Hello Bob!');
  });
});
