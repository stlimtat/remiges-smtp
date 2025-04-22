# SMTP Client Tutorial Example

This tutorial demonstrates a complete setup of SMTP client, including:
- Setting up the application config
- Running and developing with docker compose

## Directory Structure

```
/remiges-smtp
├── config                 # Location to be used to store config.yaml - see step 3
├── docker-compose.yml     # Docker Compose configuration for smtpclientd
├── Dockerfile.smtpclientd # Dockerfile that compiles the application
├── output                 # Location to be used to record output
├── README.md              # Default README.md
├── testdata               # Location to be used to hold data
└── doc                    # Document directory
  ├── TUTORIAL.md          # This document
  └── USAGE.md             # Document providing usage information
```

## Getting Started

1. [OPTIONAL] Reset docker environment
   ```bash
   docker system prune --all --force --volumes
   ```
1. Start redis:
   ```bash
   docker compose up -d redis
   ```
1. Check the [config.yaml](./config/config.yaml) in [config](./config) directory
   1.1. Ensure that the configuration for dkim is correct
   1.2. You can set up the use of a socks proxy to verify the data being exchanged via the smtp protocol
   1.3. Location of the config file is confined to the following locations, based on the following code - [root.go](https://github.com/stlimtat/remiges-smtp/blob/main/internal/config/root.go#L26)
     1.3.1. $HOME
     1.3.2. /app/config
1. Run the following to run the application in the command line as a form of integration test:
   ```bash
   bazel run //cmd/smtpclient sendmail
   ```
1. [ALTERNATIVE#1] You can also choose to run the application via raw go commands
   ```bash
   go run ./cmd/smtpclient/main.go
   ```
1. [ALTERNATIVE#2] You can also choose to run the application via docker
   ```bash
   docker compose up -d smtpclientd
   ```

## Cleanup

To stop and remove etcd container:
```bash
docker compose down
```
