FROM python:3.4-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		make \
	&& rm -rf /var/lib/apt/lists/*

COPY . /app
WORKDIR /app

RUN python3 -m compileall .


ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["bash"]
