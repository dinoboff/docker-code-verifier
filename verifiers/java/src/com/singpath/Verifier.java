package com.singpath;

import bsh.EvalError;
import bsh.Interpreter;
import bsh.TargetError;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.concurrent.*;


public class Verifier implements Callable<Response> {

    public static final int TIMEOUT = Integer.getInteger("com.singpath.verifier.timeout", 5000);

    public static final String errTooLong = "Your code took too long to return. Your solution may be stuck in an infinite loop.";
    public static final String errInvalidRequest = "No solution or tests defined";
    public static final String errImportJUnit = "Failed to import JUnit";

    private static final Logger logger = LogManager.getLogger(Verifier.class);

    private Request req;

    public Verifier(Request req) {
        super();
        this.req = req;
    }

    public static Response process(Request req) {
        ExecutorService executor = Executors.newSingleThreadExecutor();

        try {
            Future<Response> f = executor.submit(new Verifier(req));
            return f.get(TIMEOUT, TimeUnit.SECONDS);
        } catch (InterruptedException | TimeoutException | ExecutionException e) {
            return new Response(Verifier.errTooLong);
        } finally {
            executor.shutdown();
        }
    }

    protected static void addFailingResults(Failure failure, String test, Response results) {
        if (test.contains("assertTrue(")) {
            results.addResult(test, true, false);
            return;
        }

        if (test.contains("assertFalse(")) {
            results.addResult(test, false, true);
            return;
        }

        if (failure.isParsed()) {
            results.addResult(test, failure.getExpected(), failure.getActual());
        } else {
            results.addResult(test, failure.getMessage());
        }
    }

    protected static String logError(Throwable e) {
        String error = Verifier.getStackTrace(e);
        logger.error(error);
        return error;
    }

    protected static String logError(Throwable e, Response r) {
        String error = Verifier.logError(e);
        r.setErrors(error);
        return error;
    }

    protected static String getStackTrace(Throwable e) {
        StringWriter sw = new StringWriter();
        PrintWriter pw = new PrintWriter(sw);
        e.printStackTrace(pw);
        return sw.toString();
    }

    @Override
    public Response call() throws Exception {
        if (!this.req.isValid()) {
            return new Response(errInvalidRequest);
        }

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        PrintStream s = new PrintStream(out);
        Interpreter sh = new Interpreter(null, s, s, false, null);
        sh.setStrictJava(true);

        try {
            sh.eval("import static org.junit.Assert.*;");
        } catch (EvalError evalError) {
            Verifier.logError(evalError);
            return new Response(errImportJUnit);
        }

        // Eval the user text.
        Response results = new Response(out);
        try {
            sh.eval(this.req.getSolution());
        } catch (EvalError evalError) {
            Verifier.logError(evalError, results);
            return results;
        }

        // Eval tests one line at a time.
        for (String test : this.req.getTests()) {
            if (test.trim().equals("")) {
                continue;
            }

            try {
                sh.eval(test);
                if (test.contains("assert")) {
                    results.addResult(test);
                }
            } catch (TargetError targetError) {
                Failure f = Failure.fromError(targetError.getTarget());
                if (f.isAssertionError()) {
                    Verifier.addFailingResults(f, test, results);
                    continue;
                }

                logger.error(targetError.getTarget().getClass());
                Verifier.logError(targetError.getTarget(), results);
                return results;
            } catch (EvalError evalError) {
                Verifier.logError(evalError, results);
                return results;
            }
        }

        return results;
    }
}
