ARG ALPINE_VERSION=3
ARG GO_VERSION=1.21

FROM alpine:${ALPINE_VERSION} AS certificator
RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
    update-ca-certificates

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG UNCONDITIONAL_API_BUILD_COMMIT_VERSION
ARG UNCONDITIONAL_API_BUILD_RELEASE_VERSION

ENV CGO_ENABLED=0
ENV UNCONDITIONAL_API_BUILD_COMMIT_VERSION=${UNCONDITIONAL_API_BUILD_COMMIT_VERSION}
ENV UNCONDITIONAL_API_BUILD_RELEASE_VERSION=${UNCONDITIONAL_API_BUILD_RELEASE_VERSION}

RUN go build --ldflags "-X 'main.releaseVersion=${UNCONDITIONAL_API_BUILD_RELEASE_VERSION}' -X 'main.gitCommit=${UNCONDITIONAL_API_BUILD_COMMIT_VERSION}'" --tags=release -o main .

FROM builder as data
WORKDIR /data
COPY --from=builder /app/main /app/main

ARG UNCONDITIONAL_API_SOURCE_REPO
ARG UNCONDITIONAL_API_SOURCE_CLIENT_KEY
ARG UNCONDITIONAL_API_FEED_REPO_INDEX
ARG UNCONDITIONAL_API_FEED_REPO_HOST
ARG UNCONDITIONAL_API_FEED_REPO_KEY
ARG UNCONDITIONAL_API_LOG_ENV

ENV UNCONDITIONAL_API_SOURCE_REPO=${UNCONDITIONAL_API_SOURCE_REPO}
ENV UNCONDITIONAL_API_SOURCE_CLIENT_KEY=${UNCONDITIONAL_API_SOURCE_CLIENT_KEY}
ENV UNCONDITIONAL_API_FEED_REPO_INDEX=${UNCONDITIONAL_API_FEED_REPO_INDEX}
ENV UNCONDITIONAL_API_FEED_REPO_HOST=${UNCONDITIONAL_API_FEED_REPO_HOST}
ENV UNCONDITIONAL_API_FEED_REPO_KEY=${UNCONDITIONAL_API_FEED_REPO_KEY}
ENV UNCONDITIONAL_API_LOG_ENV=${UNCONDITIONAL_API_LOG_ENV}

RUN --mount=type=secret,id=UNCONDITIONAL_API_SOURCE_CLIENT_KEY \
    --mount=type=secret,id=UNCONDITIONAL_API_FEED_REPO_HOST \
    --mount=type=secret,id=UNCONDITIONAL_API_FEED_REPO_INDEX \
    --mount=type=secret,id=UNCONDITIONAL_API_FEED_REPO_KEY \
    UNCONDITIONAL_API_SOURCE_CLIENT_KEY="$(cat /run/secrets/UNCONDITIONAL_API_SOURCE_CLIENT_KEY)" \
    UNCONDITIONAL_API_FEED_REPO_HOST="$(cat /run/secrets/UNCONDITIONAL_API_FEED_REPO_HOST)" \
    UNCONDITIONAL_API_FEED_REPO_INDEX="$(cat /run/secrets/UNCONDITIONAL_API_FEED_REPO_INDEX)" \
    UNCONDITIONAL_API_FEED_REPO_KEY="$(cat /run/secrets/UNCONDITIONAL_API_FEED_REPO_KEY)" \
    /app/main index create --name feeds

FROM scratch as release
COPY --from=certificator /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/main
COPY --from=data /data /data

ARG UNCONDITIONAL_API_ADDRESS
ARG UNCONDITIONAL_API_ALLOWED_ORIGINS
ARG UNCONDITIONAL_API_PORT
ARG UNCONDITIONAL_API_SOURCE_REPO
ARG UNCONDITIONAL_API_SOURCE_CLIENT_KEY
ARG UNCONDITIONAL_API_LOG_ENV
ARG UNCONDITIONAL_API_FEED_REPO_INDEX
ARG UNCONDITIONAL_API_FEED_REPO_HOST
ARG UNCONDITIONAL_API_FEED_REPO_KEY

ENV UNCONDITIONAL_API_ADDRESS=${UNCONDITIONAL_API_ADDRESS}
ENV UNCONDITIONAL_API_ALLOWED_ORIGINS=${UNCONDITIONAL_API_ALLOWED_ORIGINS}
ENV UNCONDITIONAL_API_PORT=${UNCONDITIONAL_API_PORT}
ENV UNCONDITIONAL_API_SOURCE_REPO=${UNCONDITIONAL_API_SOURCE_REPO}
ENV UNCONDITIONAL_API_SOURCE_CLIENT_KEY=${UNCONDITIONAL_API_SOURCE_CLIENT_KEY}
ENV UNCONDITIONAL_API_LOG_ENV=${UNCONDITIONAL_API_LOG_ENV}
ENV UNCONDITIONAL_API_FEED_REPO_INDEX=${UNCONDITIONAL_API_FEED_REPO_INDEX}
ENV UNCONDITIONAL_API_FEED_REPO_HOST=${UNCONDITIONAL_API_FEED_REPO_HOST}
ENV UNCONDITIONAL_API_FEED_REPO_KEY=${UNCONDITIONAL_API_FEED_REPO_KEY}

ENTRYPOINT ["./app/main","serve", "--address", "0.0.0.0", "--port","8080"]
