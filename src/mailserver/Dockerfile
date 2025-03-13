FROM node:18-bullseye

WORKDIR /app
COPY . .

RUN npm install
ENV HOSTNAME=localhost
# Ensure correct port is exposed
EXPOSE 0000
RUN npm rebuild better-sqlite3
# Use `sh -c` to run the command properly
CMD ["sh", "-c", "npm start"]
