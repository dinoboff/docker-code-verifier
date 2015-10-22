package com.singpath;

import org.junit.ComparisonFailure;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Failure {
    public static final String assertErrorMsgPatternStr = "^expected:<(.*)> but was:<(.*)>$";
    public static final Pattern assertErrorMsgPattern = Pattern.compile(assertErrorMsgPatternStr);

    private String expected;
    private String actual;
    private String message;

    public Failure(String msg) {
        super();
        this.message = msg;
    }

    public Failure(ComparisonFailure e) {
        this(e.getMessage());
        this.expected = e.getExpected();
        this.actual = e.getActual();
    }

    public Failure(AssertionError e) {
        this(e.getMessage());

        Matcher matcher = assertErrorMsgPattern.matcher(this.message);
        if (matcher.find()) {
            this.expected = matcher.group(1);
            this.actual = matcher.group(2);
            return;
        }

        //if the regular expression fails, use the old method
        String failS = this.message.replace("expected:<", "");
        failS = failS.replace("> but was:<", ":SEP:");
        failS = failS.replace(">", "");
        String[] ss = failS.split(":SEP:", 2);
        if (ss.length == 2) {
            this.expected = ss[0];
            this.actual = ss[1];
        }
    }

    public static Failure fromError(Throwable e) {
        Class<? extends Throwable> k = e.getClass();

        if (k.equals(ComparisonFailure.class)) {
            return new Failure((ComparisonFailure) e);
        }

        if (k.equals(AssertionError.class)) {
            return new Failure((AssertionError) e);
        }

        return new Failure(e.getMessage());
    }

    public String getExpected() {
        return this.expected;
    }

    public String getActual() {
        return this.actual;
    }

    public String getMessage() {
        return this.message;
    }

    public boolean isAssertionError() {
        return this.message != null;
    }

    public boolean isParsed() {
        return this.actual != null;
    }
}
