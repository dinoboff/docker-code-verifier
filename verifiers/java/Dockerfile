FROM java:8-jdk

RUN mkdir -p /app/build/dist && \
	mkdir -p /app/build/tests && \
	adduser --system --home /playground verifier

COPY ./lib /app/lib
COPY ./src /app/src
ENV CLASSPATH /app/lib/*:/app/build/dist
RUN javac -sourcepath /app/src -d /app/build/dist -g /app/src/com/singpath/server/Main.java

COPY ./tests /app/tests
RUN javac -sourcepath /app/tests -d /app/build/tests -g /app/tests/com/singpath/TestSuite.java

COPY ./log4j2.xml /app/log4j2.xml

USER verifier
WORKDIR /playground

EXPOSE 5000

CMD ["java", "-Xmx70m", "-Dlog4j.configurationFile=/app/log4j2.xml", "com.singpath.server.Main"]