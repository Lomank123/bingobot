FROM golang:1.22.4 as builder

WORKDIR /app

COPY . .

RUN go mod download && \
    go build -v -o main cmd/telegram/main.go

# We have to run the binary on other image
FROM debian:bookworm-slim
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y \
    curl \
    iputils-ping && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

# To load env using go package
COPY --from=builder /app/.env .
# Copy generated executable file
COPY --from=builder /app/main main

CMD ["/main"]