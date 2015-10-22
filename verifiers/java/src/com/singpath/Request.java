package com.singpath;

public class Request {

    private String tests;
    private String solution;

    public Request(String solution, String tests) {
        this.tests = tests;
        this.solution = solution;
    }

    public String getSolution() {
        return solution;
    }

    public String[] getTests() {
        return tests.split("\n");
    }

    public boolean isValid() {
        return this.tests != null
                && this.tests.length() > 0
                && this.solution != null
                && this.solution.length() > 0;
    }
}
