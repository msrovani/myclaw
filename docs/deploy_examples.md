# XXXCLAW Deployment Examples

This document provides examples for deploying XXXCLAW in various environments.

## Docker Compose

The simplest way to run XXXCLAW alongside a local Ollama instance and configure persistence is via Docker Compose.

```yaml
version: '3.8'

services:
  xxxclaw:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - xxxclaw_data:/app/data
    environment:
      - XXXCLAW_ENV=prod
      - XXXCLAW_OLLAMA_URL=http://host.docker.internal:11434
      - XXXCLAW_LOG_LEVEL=info
    restart: unless-stopped

volumes:
  xxxclaw_data:
```

*Note: Ensure `host.docker.internal` resolves correctly in your Docker bridging network.*

## Systemd Service (Linux)

For bare-metal or VM deployments, managing XXXCLAW via systemd ensures it starts on boot and restarts on failure.

1. Build the binary and move it to `/usr/local/bin/`
2. Create `/etc/systemd/system/xxxclaw.service`:

```ini
[Unit]
Description=XXXCLAW Autonomous Agent Server
After=network.target

[Service]
Type=simple
User=xxxclaw
Group=xxxclaw
WorkingDirectory=/var/lib/xxxclaw
ExecStart=/usr/local/bin/xxxclaw
Environment="XXXCLAW_HTTP_ADDR=:8080"
Environment="XXXCLAW_ENV=prod"
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## Nginx Reverse Proxy

To expose XXXCLAW to the public securely, use a reverse proxy like Nginx to handle SSL termination.

```nginx
server {
    listen 80;
    server_name claw.yourdomain.com;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```
