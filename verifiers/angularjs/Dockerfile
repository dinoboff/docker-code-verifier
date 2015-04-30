FROM node:0.10.38-slim

COPY . /app
WORKDIR /app
RUN npm install -g protractor \
	&& npm install

EXPOSE 5000
CMD ["node", "/app/server.js", "--host=0.0.0.0", "--port=5000"]
