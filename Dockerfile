ARG GO_VERSION=1.20.6

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
RUN /app/main source download --path /data/source.json 
RUN /app/main index create --source /data/source.json --name /data/index

FROM scratch as release
COPY --from=builder /app/main /app/main
COPY --from=data /data /data

ARG UNCONDITIONAL_ADDRESS
ARG UNCONDITIONAL_ALLOWED_ORIGINS
ARG UNCONDITIONAL_PORT

ENV UNCONDITIONAL_ADDRESS=${UNCONDITIONAL_ADDRESS}
ENV UNCONDITIONAL_ALLOWED_ORIGINS=${UNCONDITIONAL_ALLOWED_ORIGINS}
ENV UNCONDITIONAL_PORT=${UNCONDITIONAL_PORT}

ENTRYPOINT ["./app/main","serve", "--address", "0.0.0.0", "--port","8080", "--index","/data/index"]
