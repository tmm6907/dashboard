# Use a lightweight Go base image
FROM golang:1.23.4 AS build

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server .
# Create a minimal final image
FROM debian:bookworm-slim AS final
RUN apt-get update && apt-get install -y \
    ca-certificates \
    sqlite3
WORKDIR /root/
COPY --from=build /app/server .
COPY build.sql .
EXPOSE 8080

CMD ["./server"]

