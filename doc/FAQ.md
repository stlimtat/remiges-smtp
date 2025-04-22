# Frequently Asked Questions

## General Questions

### What is Remiges SMTP?
Remiges SMTP is an SMTP client that can read files from directories, format them as emails, sign them with DKIM, and send them to specified recipients.

### What are the system requirements?
- Docker (version 20.10.0 or later)
- Docker Compose (version 2.0.0 or later)
- Go (version 1.21 or later) - only if building from source
- Bazel (version 6.0.0 or later) - only if using Bazel

## Installation

### How do I install Remiges SMTP?
The easiest way is using Docker:
```bash
git clone https://github.com/stlimtat/remiges-smtp.git
cd remiges-smtp
docker compose up -d
```

### I'm getting errors during installation. What should I do?
1. Check that Docker and Docker Compose are installed and running:
   ```bash
   docker --version
   docker compose version
   ```
2. Ensure you have sufficient disk space and memory
3. Check the logs for specific error messages:
   ```bash
   docker compose logs
   ```

## Configuration

### Where should I put my configuration file?
The configuration file can be placed in either:
- `$HOME/config.yaml`
- `/app/config/config.yaml` (when using Docker)

### What are the essential configuration options?
At minimum, you need to configure:
```yaml
debug: false
read_file:
  in_path: /path/to/mail
  redis_addr: localhost:6379
```

## DKIM

### How do I set up DKIM?
1. Generate DKIM keys:
   ```bash
   smtpclient gendkim --dkim-domain yourdomain.com
   ```
2. Add the provided TXT record to your DNS
3. Configure the private key path in your config.yaml

### My DKIM signatures are failing. What should I check?
1. Verify the DNS TXT record is correctly set up
2. Check that the private key path in config.yaml is correct
3. Ensure the domain in the DKIM signature matches your DNS record

## Troubleshooting

### How do I check if the service is running?
```bash
docker compose ps
```

### Where can I find the logs?
```bash
docker compose logs smtpclientd
```

### My emails aren't being sent. What should I check?
1. Verify Redis is running:
   ```bash
   docker compose ps redis
   ```
2. Check the SMTP server configuration
3. Look for errors in the logs
4. Verify network connectivity to the SMTP server

## Performance

### How can I improve performance?
1. Increase concurrency in the configuration:
   ```yaml
   read_file:
     concurrency: 4
   ```
2. Adjust the poll interval:
   ```yaml
   read_file:
     poll_interval: 5s
   ```
3. Ensure sufficient resources for Redis

## Security

### How do I secure my configuration?
1. Never commit sensitive configuration files to version control
2. Use environment variables for sensitive data
3. Set appropriate file permissions on configuration files
4. Use secure connections (TLS) for Redis

### How do I handle API keys and passwords?
Store them in environment variables or a secure secrets management system. Never hardcode them in configuration files.