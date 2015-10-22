package com.singpath;

import org.junit.Assert;
import org.junit.ComparisonFailure;
import org.junit.Test;

import static org.junit.Assert.*;

public class FailureTest {

    @Test
    public void testFromComparisonFailure() throws Exception {
        ComparisonFailure e = new ComparisonFailure("", "foo", "bar");
        Failure f = Failure.fromError(e);
        assertTrue(f.isAssertionError());
    }

    @Test
    public void testFromAssertionError() throws Exception {
        AssertionError e = new AssertionError("Some message");
        Failure f = Failure.fromError(e);
        assertTrue(f.isAssertionError());
    }

    @Test
    public void testFromSomeError() throws Exception {
        Exception e = new Exception();
        Failure f = Failure.fromError(e);
        assertFalse(f.isAssertionError());
    }


    @Test
    public void testGetExpected() throws Exception {
        ComparisonFailure e = new ComparisonFailure("", "foo", "bar");
        Failure f = new Failure(e);
        assertEquals("foo", f.getExpected());
    }

    @Test
    public void testGetActual() throws Exception {
        ComparisonFailure e = new ComparisonFailure("", "foo", "bar");
        Failure f = new Failure(e);
        assertEquals("bar", f.getActual());
    }

    @Test
    public void testGetMessage() throws Exception {
        AssertionError e = new AssertionError("Some error");
        Failure f = new Failure(e);
        assertEquals("Some error", f.getMessage());
    }

    @Test
    public void testIsParsed1() throws Exception {
        ComparisonFailure e = new ComparisonFailure("", "foo", "bar");
        Failure f = new Failure(e);
        assertTrue(f.isParsed());
    }

    @Test
    public void testIsParsed2() throws Exception {
        AssertionError e = new AssertionError("^expected:<foo> but was:<bar>$");
        Failure f = new Failure(e);
        assertTrue(f.isParsed());
    }

    @Test
    public void testIsNotParsed() throws Exception {
        AssertionError e = new AssertionError("Some error");
        Failure f = new Failure(e);
        assertFalse(f.isParsed());
    }
}