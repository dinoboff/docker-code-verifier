package com.singpath;

import com.singpath.tools.MemoryClassLoader;
import com.singpath.tools.MemoryJavaFileManager;
import com.singpath.tools.StringJavaSource;

import javax.tools.JavaCompiler;
import javax.tools.JavaFileObject;
import javax.tools.StandardJavaFileManager;
import javax.tools.ToolProvider;
import java.io.Writer;
import java.util.Arrays;

public class Request {

    public final String SOLUTION_CLASS_NAME = "SingPath";
    public final String TEST_CLASS_NAME = "SingPathTest";
    public final String TEST_TEMPLATE = "import org.junit.Test;\n" +
            "import static org.junit.Assert.*;\n" +
            "import junit.framework.*;\n" +
            "import java.io.ByteArrayOutputStream;\n" +
            "import java.io.PrintStream;\n" +
            "import com.singpath.SolutionRunner;\n" +
            "\n" +
            "public class SingPathTest extends SolutionRunner {\n" +
            "\n" +
            "    @Test\n" +
            "    public void testSolution() throws Exception {\n" +
            "        %s\n" +
            "    }\n" +
            "}";

    private String tests;
    private String solution;

    public String getTests() {
        return tests;
    }

    public String getSolution() {
        return solution;
    }

    public Request(String solution, String tests) {
        this.tests = tests;
        this.solution = solution;
    }

    public boolean isValid() {
        return this.tests != null
                && this.tests.length() > 0
                && this.solution != null
                && this.solution.length() > 0;
    }
}
