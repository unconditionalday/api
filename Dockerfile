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

ENV CGO_ENABLED=0
RUN go build -o main .

FROM builder as data
WORKDIR /data
COPY --from=builder /app/main /app/main

ARG UNCONDITIONAL_API_SOURCE_REPO
ARG UNCONDITIONAL_API_LOG_ENV
ENV UNCONDITIONAL_API_SOURCE_REPO=${UNCONDITIONAL_API_SOURCE_REPO}
ENV UNCONDITIONAL_API_LOG_ENV=${UNCONDITIONAL_API_LOG_ENV}

RUN --mount=type=secret,id=UNCONDITIONAL_API_SOURCE_CLIENT_KEY \
    UNCONDITIONAL_API_SOURCE_CLIENT_KEY="$(cat /run/secrets/UNCONDITIONAL_API_SOURCE_CLIENT_KEY)" \
    /app/main source download --path /data/source.json 

RUN /app/main index create --source /data/source.json --name /data/index

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

ENV UNCONDITIONAL_API_ADDRESS=${UNCONDITIONAL_API_ADDRESS}
ENV UNCONDITIONAL_API_ALLOWED_ORIGINS=${UNCONDITIONAL_API_ALLOWED_ORIGINS}
ENV UNCONDITIONAL_API_PORT=${UNCONDITIONAL_API_PORT}
ENV UNCONDITIONAL_API_SOURCE_REPO=${UNCONDITIONAL_API_SOURCE_REPO}
ENV UNCONDITIONAL_API_SOURCE_CLIENT_KEY=${UNCONDITIONAL_API_SOURCE_CLIENT_KEY}
ENV UNCONDITIONAL_API_LOG_ENV=${UNCONDITIONAL_API_LOG_ENV}

ENTRYPOINT ["./app/main","serve", "--address", "0.0.0.0", "--port","8080", "--index","/data/index"]
