# Remiges SMTP Client Usage Guide

This document provides comprehensive documentation for using the Remiges SMTP client application.

## Table of Contents
1. [Installation](#installation)
2. [Command Line Options](#command-line-options)
3. [Docker Environment](#docker-environment)
4. [DKIM Setup](#dkim-setup)
5. [Configuration](#configuration)
6. [Examples](#examples)

## Installation

### Using Bazel
```sh
bazel build //cmd/smtpclient
```

### Using Docker
```sh
docker build -t remiges-smtp .
```

## Command Line Options

The application provides several commands with their respective options:

### Global Flags
- `--debug, -d`: Run the application in debug mode
- `--config, -c`: Specify config file path (default: $HOME/config.yaml)

### Commands

1. **server** - Run the SMTP client server
```sh
smtpclient server
```
- Starts an HTTP server on port 8000 for administration
- Processes mail queue continuously

2. **sendmail** - Send individual emails
```sh
smtpclient sendmail [flags]
```
Flags:
- `--from, -f`: Sender email address
- `--to, -t`: Destination email address
- `--msg, -m`: Test message content
- `--path, -p`: Path to the directory containing df and qf files
- `--redis-addr, -r`: Redis server address

3. **gendkim** - Generate DKIM keys and configuration
```sh
smtpclient gendkim [flags]
```
Flags:
- `--algorithm`: Key type (default: "rsa")
- `--bit-size`: Key size (default: 2048)
- `--dkim-domain`: Domain for DKIM
- `--hash`: Hash algorithm (default: "sha256")
- `--out-path`: Output path for keys (default: "./config")
- `--selector`: DKIM selector (default: "key001")

4. **lookupmx** - Look up MX records
```sh
smtpclient lookupmx [flags]
```
Flags:
- `--lookup-domain, -l`: Domain to lookup MX entries

5. **readfile** - Read mail files
```sh
smtpclient readfile [flags]
```
Flags:
- `--path, -p`: Path to the directory containing df and qf files

## Docker Environment

### Building the Docker Image
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o smtpclient ./cmd/smtpclient

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/smtpclient .
COPY config.yaml .

EXPOSE 8000
ENTRYPOINT ["/app/smtpclient"]
CMD ["server"]
```

### Running with Docker
```sh
# Build the image
docker build -t remiges-smtp .

# Run the server
docker run -d \
  -p 8000:8000 \
  -v /path/to/config:/app/config \
  -v /path/to/mail:/app/mail \
  --name remiges-smtp \
  remiges-smtp

# Run specific commands
docker run --rm remiges-smtp gendkim --dkim-domain example.com
```

### Docker Compose Setup
```yaml
version: '3.8'
services:
  smtp:
    image: remiges-smtp
    ports:
      - "8000:8000"
    volumes:
      - ./config:/app/config
      - ./mail:/app/mail
    environment:
      - REDIS_ADDR=redis:6379
    depends_on:
      - redis

  redis:
    image: redis:alpine
    ports:
      - "6379:6379"
```

## DKIM Setup

### 1. Generate DKIM Keys
```sh
smtpclient gendkim \
  --dkim-domain example.com \
  --selector key001 \
  --algorithm rsa \
  --bit-size 2048 \
  --hash sha256 \
  --out-path ./config
```

### 2. Configure DNS Records
After generating the keys, add the provided TXT record to your DNS configuration:

```sh
key001._domainkey.example.com IN TXT "v=DKIM1; k=rsa; p=<public-key>"
```

## Configuration

### Sample Configuration
```yaml
debug: false
read_file:
  in_path: /path/to/mail
  concurrency: 4
  poll_interval: 5s
  redis_addr: localhost:6379

mail_processors:
  - type: unixdos
    index: 0
  - type: body
    index: 1
  - type: bodyHeaders
    index: 2
  - type: mergeHeaders
    index: 11
  - type: dkim
    index: 12
    args:
      domain-str: example.com
      dkim:
        selectors:
          key001:
            algorithm: rsa
            hash: sha256
            private-key-file: ./config/key001.pem
  - type: mergeHeaders
    index: 13
  - type: mergeBody
    index: 99

outputs:
  - type: file
    args:
      path: /path/to/output
```

## Examples

### 1. Send a Test Email
```sh
smtpclient sendmail \
  --from sender@example.com \
  --to recipient@example.com \
  --msg "Test message"
```

### 2. Process Mail Queue
```sh
smtpclient server
```

### 3. Look up MX Records
```sh
smtpclient lookupmx --lookup-domain example.com.
```

### 4. Read Mail Files
```sh
smtpclient readfile --path /path/to/mail
```
