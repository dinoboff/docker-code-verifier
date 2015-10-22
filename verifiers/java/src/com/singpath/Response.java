package com.singpath;

import net.minidev.json.JSONArray;
import net.minidev.json.JSONObject;
import net.minidev.json.JSONStyle;

import java.io.ByteArrayOutputStream;
import java.nio.charset.Charset;


public class Response {

    private String errors;
    private Boolean solved;
    private JSONArray results;
    private ByteArrayOutputStream out;

    public Response() {
        this(true);
    }

    public Response(ByteArrayOutputStream out) {
        this(true);
        this.setOutputStream(out);
    }

    public Response(boolean solved) {
        this.setSolved(solved);
    }

    public Response(String errors) {
        this.setErrors(errors);
    }

    public void addResult(String call) {
        if (this.results == null) {
            this.results = new JSONArray();
        }

        JSONObject result = new JSONObject();
        result.put("call", call);
        result.put("correct", true);
        this.results.add(result);
    }

    public void addResult(String call, String msg) {
        if (this.results == null) {
            this.results = new JSONArray();
        }

        JSONObject result = new JSONObject();
        result.put("call", call);
        result.put("correct", false);
        result.put("error", msg);
        this.results.add(result);
        this.setSolved(false);
    }

    public void addResult(String call, Object expected, Object received) {
        if (this.results == null) {
            this.results = new JSONArray();
        }

        JSONObject result = new JSONObject();
        result.put("call", call);
        result.put("correct", false);
        result.put("expected", expected);
        result.put("received", received);
        this.results.add(result);
        this.setSolved(false);
    }

    public void setSolved(Boolean solved) {
        this.solved = solved;
    }

    public void setErrors(String errors) {
        this.solved = false;
        this.errors = errors;
    }

    public void setOutputStream(ByteArrayOutputStream out) {
        this.out = out;
    }

    public String getOutput() {
        if (this.out == null) {
            return "";
        }

        return new String(this.out.toByteArray(), Charset.forName("UTF-8"));
    }

    @Override
    public String toString() {
        JSONObject json = new JSONObject();

        json.put("solved", this.solved);
        if (this.errors != null) {
            json.put("errors", this.errors);
        }

        if (this.results != null && this.results.size() > 0) {
            json.put("results", this.results);
        }

        String out = this.getOutput();
        if (out.length() > 0) {
            json.put("printed", out);
        }

        return json.toString(JSONStyle.NO_COMPRESS);
    }
}
