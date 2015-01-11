FROM selenium/node-base:2.44.0

USER root

ENV NODE_PORT 8080
ENV PHANTOMJS_VERSION 1.9.8-linux-x86_64
RUN wget https://bitbucket.org/ariya/phantomjs/downloads/phantomjs-${PHANTOMJS_VERSION}.tar.bz2 -O /tmp/phantomjs.tar.bz2 \
	&& (cd /opt/ ; tar xvjf /tmp/phantomjs.tar.bz2 ; ) \
	&& ln -s /opt/phantomjs-${PHANTOMJS_VERSION}/bin/phantomjs /usr/bin/phantomjs \
	&& /usr/bin/phantomjs -v

COPY entry_point.sh /opt/bin/entry_point.sh
RUN chmod +x /opt/bin/entry_point.sh

USER seluser

CMD ["/opt/bin/entry_point.sh"]