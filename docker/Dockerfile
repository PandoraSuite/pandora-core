FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -ldflags "-X 'internal.version.Version=0.0.3-beta'" \
    -o pandora-core ./cmd/http/main.go

FROM alpine:3.19

RUN apk add --no-cache postgresql postgresql-contrib su-exec bash

RUN addgroup -S pandora && adduser -S -G pandora pandora

RUN mkdir -p /etc/pandora /var/lib/pandora /opt/pandora \
        /run/postgresql /var/lib/postgresql/data && \
    chown -R pandora:pandora /etc/pandora /var/lib/pandora \
        /opt/pandora /run/postgresql /var/lib/postgresql

COPY --from=builder /app/pandora-core /usr/local/bin/pandora-core

WORKDIR /var/lib/pandora

COPY ./docker/run.sh ./db/init.sql /opt/pandora/

RUN chmod +x /opt/pandora/run.sh

ENV PANDORA_DB_DNS="postgresql://pandora_user:supersecure@localhost:5432/pandora"

USER pandora

# VOLUME /var/lib/postgresql/data

EXPOSE 80 5432

ENTRYPOINT ["/opt/pandora/run.sh"]
