# Use Node.js to build Astro frontend
FROM node:18 AS build

WORKDIR /app

# Install dependencies
COPY package.json package-lock.json ./
RUN npm install

# Copy the rest of the files and build Astro
COPY . .
RUN npm run build

# Use Caddy to serve the static files
FROM caddy:latest
COPY --from=build /app/dist /srv

EXPOSE 80
CMD ["caddy", "file-server", "--root", "/srv"]