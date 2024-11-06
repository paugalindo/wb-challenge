FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/api

FROM postgres:latest AS final

ENV POSTGRES_USER=wbuser
ENV POSTGRES_PASSWORD=wbpass
ENV POSTGRES_DB=wbdb

COPY --from=builder /app/main /usr/local/bin/main

COPY init.sql /docker-entrypoint-initdb.d/

RUN mkdir -p /path/to/local/data
RUN chmod 700 /path/to/local/data

RUN apt-get update && apt-get install -y \
    wget \
    tar \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/nats-io/nats-server/releases/download/v2.9.0/nats-server-v2.9.0-linux-amd64.tar.gz && \
    tar -xvf nats-server-v2.9.0-linux-amd64.tar.gz && \
    mv nats-server-v2.9.0-linux-amd64/nats-server /usr/local/bin/ && \
    rm -rf nats-server-v2.9.0-linux-amd64.tar.gz nats-server-v2.9.0-linux-amd64

USER postgres

EXPOSE 80 5432 4222

CMD ["sh", "-c", "docker-entrypoint.sh postgres & nats-server & /usr/local/bin/main"]