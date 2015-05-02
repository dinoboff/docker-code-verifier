FROM node:0.12.2-slim

COPY . /app/
WORKDIR /app
RUN npm install && adduser --system verifer
USER verifer


EXPOSE 5000
CMD ["node", "/app/server.js", "--host=0.0.0.0", "--port=5000"]
