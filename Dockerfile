ARG GO_VERSION=1.19.1

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o cmd .

FROM builder as data
WORKDIR /data
COPY --from=builder /app/cmd /app/cmd
RUN /app/cmd source download --path /data/source.json && 
RUN /app/cmd index create --source /data/source.json --name /data/index

FROM scratch as release
COPY --from=builder /app/cmd /app/cmd
COPY --from=data /data /data

ARG IT_ADDRESS
ARG IT_PORT
ENV IT_ADDRESS=${IT_ADDRESS}
ENV IT_PORT=${IT_PORT}

ENTRYPOINT ["./app/cmd","serve","--index","/data/index"]