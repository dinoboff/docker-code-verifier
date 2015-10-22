package com.singpath;

import org.junit.Test;

import static org.junit.Assert.*;


public class RequestTest {

    @Test
    public void testGetSolution() throws Exception {
        Request req = new Request("solution", "test");
        assertEquals("solution", req.getSolution());
    }

    @Test
    public void testGetTests() throws Exception {
        Request req = new Request("solution", "test1\ntest2");
        String[] tests = req.getTests();
        assertEquals(2, tests.length);
        assertEquals("test1", tests[0]);
        assertEquals("test2", tests[1]);
    }

    @Test
    public void testIsValid() throws Exception {
        Request req = new Request("", "");
        assertFalse(req.isValid());
    }
}