ARG GO_VERSION=1.19.1

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

ARG IT_ADDRESS
ARG IT_PORT
ENV IT_ADDRESS=${IT_ADDRESS}
ENV IT_PORT=${IT_PORT}

ENTRYPOINT ["./app/main","serve","--index","/data/index"]