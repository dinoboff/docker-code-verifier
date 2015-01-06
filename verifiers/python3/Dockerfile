FROM python:3.4-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		make \
	&& rm -rf /var/lib/apt/lists/* \
	&& pip install flask \
	&& adduser --system --home /playground verifier

COPY . /app
RUN python3 -m compileall /app/

USER verifier
WORKDIR /playground

EXPOSE 5000
CMD ["python3", "/app/server.py", "--host=0.0.0.0", "--port=5000"]
