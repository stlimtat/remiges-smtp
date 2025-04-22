# Quick Start Guide

## Prerequisites
Before you begin, ensure you have the following installed:
- Docker (version 20.10.0 or later)
- Docker Compose (version 2.0.0 or later)
- Go (version 1.21 or later) - only if you want to build from source
- Bazel (version 6.0.0 or later) - only if you want to use Bazel

## Installation Options

### Option 1: Using Docker (Recommended for Beginners)
1. Clone the repository:
   ```bash
   git clone https://github.com/stlimtat/remiges-smtp.git
   cd remiges-smtp
   ```

2. Start the services:
   ```bash
   docker compose up -d
   ```

### Option 2: Building from Source
1. Clone the repository:
   ```bash
   git clone https://github.com/stlimtat/remiges-smtp.git
   cd remiges-smtp
   ```

2. Build using Bazel:
   ```bash
   bazel build //cmd/smtpclient
   ```

## Your First Email

Let's send a test email to verify your setup:

1. Create a basic configuration file:
   ```bash
   mkdir -p config
   cat > config/config.yaml << 'EOF'
   debug: true
   read_file:
     in_path: ./testdata
     concurrency: 1
     poll_interval: 5s
     redis_addr: localhost:6379
   EOF
   ```

2. Send a test email:
   ```bash
   # Using Docker
   docker compose exec smtpclientd smtpclient sendmail \
     --from test@example.com \
     --to recipient@example.com \
     --msg "Hello from Remiges SMTP!"

   # Or using the binary directly
   bazel run //cmd/smtpclient -- sendmail \
     --from test@example.com \
     --to recipient@example.com \
     --msg "Hello from Remiges SMTP!"
   ```

## Next Steps
- Learn about [configuration options](./USAGE.md#configuration)
- Set up [DKIM signing](./USAGE.md#dkim-setup)
- Explore [advanced features](./TUTORIAL.md)

## Troubleshooting
If you encounter issues:
1. Check the logs:
   ```bash
   docker compose logs smtpclientd
   ```
2. Verify Redis is running:
   ```bash
   docker compose ps
   ```
3. Check the configuration file permissions and location

## Need Help?
- Check the [FAQ](./FAQ.md)
- Open an [issue](https://github.com/stlimtat/remiges-smtp/issues)
- Join our [community chat](https://github.com/stlimtat/remiges-smtp/discussions)