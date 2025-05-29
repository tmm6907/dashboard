ARG GO_VERSION=1.23.4
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build

WORKDIR /mail

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN apt-get update && apt-get install -y gcc libc6-dev && rm -rf /var/lib/apt/lists/* 

ARG TARGETARCH
ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    GOARCH=$TARGETARCH go build -o /tmp/mail .

# Final stage
FROM debian:bookworm-slim AS final

RUN apt-get update && apt-get install -y \
    postfix \
    ca-certificates \
    tzdata \
    sqlite3 \
    zip && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /mail

COPY postfix/main.cf /etc/postfix/main.cf
COPY start.sh /mail/start.sh
COPY ./build.sql /mail/build.sql
COPY --from=build /tmp/mail /mail/mailserver

RUN chmod 644 /etc/postfix/main.cf
RUN chmod +x /mail/start.sh

EXPOSE 25 8888

ENTRYPOINT ["/mail/start.sh"]
