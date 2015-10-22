package com.singpath;

import net.minidev.json.JSONObject;
import net.minidev.json.JSONValue;
import org.junit.Test;

import static org.junit.Assert.*;

public class VerifierTest {

    @Test
    public void testProcess() throws Exception {
        Request req = new Request("int foo = 1", "assertEquals(1, foo)");
        Response resp = Verifier.process(req);

        JSONObject dict = (JSONObject) JSONValue.parse(resp.toString());
        boolean solved = (boolean) dict.get("solved");

        assertEquals(true, solved);
    }
}