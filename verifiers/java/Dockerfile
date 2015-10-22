FROM java:8-jre

RUN mkdir -p /app && \
	adduser --system --home /playground verifier

COPY ./lib /app/lib
COPY ./log4j2.xml /app/log4j2.xml
COPY ./build/dist/verifier.jar /app/verifier.jar

USER verifier
WORKDIR /playground

EXPOSE 5000

CMD ["java", "-Xmx70m", "-Dlog4j.configurationFile=/app/log4j2.xml", "-cp", "/app/lib/*:/app/verifier.jar", "com.singpath.Server"]