package com.singpath;

import org.junit.runner.RunWith;
import org.junit.runners.Suite;

@RunWith(Suite.class)
@Suite.SuiteClasses({
        RequestTest.class,
        ResponseTest.class,
        VerifierTest.class
})
public class TestSuite {
}
