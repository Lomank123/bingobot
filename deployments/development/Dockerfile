FROM golang:1.22.4 as builder

WORKDIR /app

COPY . .

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    curl \
    iputils-ping && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean && \
    go mod download
