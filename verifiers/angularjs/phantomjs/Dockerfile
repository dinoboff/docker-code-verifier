FROM debian:wheezy

USER root

ENV SCREEN_WIDTH 1360
ENV SCREEN_HEIGHT 1020
ENV SCREEN_DEPTH 24
ENV DISPLAY :99.0
ENV PHANTOMJS_VERSION 1.9.8-linux-x86_64
RUN apt-get update \
	&& apt-get -qqy install \
		ca-certificates \
		bzip2 \
		chrpath \
		libfontconfig1 \
		libfreetype6 \
		netcat \
		wget \
		xvfb \
	&& rm -rf /var/lib/apt/lists/* \
	&& wget -q https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-${PHANTOMJS_VERSION}.tar.bz2 -O /tmp/phantomjs.tar.bz2 \
	&& (cd /opt/ ; tar xvjf /tmp/phantomjs.tar.bz2 ; ) \
	&& ln -s /opt/phantomjs-${PHANTOMJS_VERSION}/bin/phantomjs /usr/bin/phantomjs \
	&& /usr/bin/phantomjs -v

COPY entry_point.sh /opt/bin/entry_point.sh
RUN chmod +x /opt/bin/entry_point.sh

RUN useradd --system verifier
USER verifier

EXPOSE 8080
CMD ["/opt/bin/entry_point.sh"]
