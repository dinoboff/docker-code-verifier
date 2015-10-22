package com.singpath;

import net.minidev.json.JSONArray;
import net.minidev.json.JSONObject;
import net.minidev.json.JSONValue;
import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;

import static org.junit.Assert.*;

public class ResponseTest {

    @Test
    public void testAddResultPass() throws Exception {
        Response resp = new Response();
        resp.addResult("assertTrue(true)");

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        JSONArray results = (JSONArray) dict.get("results");

        assertEquals(1, results.size());

        JSONObject firstTest = (JSONObject) results.get(0);
        assertEquals("assertTrue(true)", firstTest.get("call"));
        assertEquals(true, firstTest.get("correct"));
    }

    @Test
    public void testAddResultFailed() throws Exception {
        Response resp = new Response();
        resp.addResult("assertTrue(true)");
        resp.addResult("assertTrue(false)", true, false);

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        JSONArray results = (JSONArray) dict.get("results");

        assertEquals(2, results.size());

        JSONObject firstTest = (JSONObject) results.get(1);
        assertEquals("assertTrue(false)", firstTest.get("call"));
        assertEquals(false, firstTest.get("correct"));
        assertEquals(true, firstTest.get("expected"));
        assertEquals(false, firstTest.get("received"));
    }

    @Test
    public void testAddResultFailedJustMsg() throws Exception {
        Response resp = new Response();
        resp.addResult("assertTrue(true)");
        resp.addResult("assertTrue(false)", true, false);
        resp.addResult("assertSomething(foo)", "some message I can't parse");

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        JSONArray results = (JSONArray) dict.get("results");

        assertEquals(3, results.size());

        JSONObject firstTest = (JSONObject) results.get(2);
        assertEquals("assertSomething(foo)", firstTest.get("call"));
        assertEquals(false, firstTest.get("correct"));
        assertEquals("some message I can't parse", firstTest.get("error"));
    }

    @Test
    public void testSetSolvedFalse() throws Exception {
        Response resp = new Response();
        resp.setSolved(false);

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        boolean solved = (boolean) dict.get("solved");

        assertEquals(false, solved);
    }

    @Test
    public void testSetSolvedTrue() throws Exception {
        Response resp = new Response();
        resp.setSolved(true);

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        boolean solved = (boolean) dict.get("solved");

        assertEquals(true, solved);
    }

    @Test
    public void testSetErrors() throws Exception {
        Response resp = new Response();
        resp.setErrors("You're bad.");

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        boolean solved = (boolean) dict.get("solved");
        String error = (String) dict.get("errors");

        assertEquals(false, solved);
        assertEquals("You're bad.", error);
    }

    @Test
    public void testSetOutputStream() throws Exception {
        Response resp = new Response();

        JSONObject before = (JSONObject) JSONValue.parse(resp.toString());
        String noOutput = (String) before.get("printed");

        assertNull(noOutput);

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        PrintStream s = new PrintStream(out);
        resp.setOutputStream(out);

        s.print("some output");

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        String printed = (String) dict.get("printed");

        assertEquals("some output", printed);
    }
}