FROM node:18 AS build

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

CMD ["npm", "run", "preview", "--", "--port", "3030"]
