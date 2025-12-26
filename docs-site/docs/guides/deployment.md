---
sidebar_position: 7
---

# Deployment Guide

This guide covers deploying forge applications to production.

## Preparation

### Build for Production

Build your application:

```bash
go build -o myapp ./main.go
```

Or with optimizations:

```bash
go build -ldflags="-s -w" -o myapp ./main.go
```

### Environment Variables

Use environment variables for configuration:

```go
import (
    "os"
    "github.com/forgego/forge/pkg/config"
)

func main() {
    cfg := config.NewConfig()
    
    // Override with environment variables
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        cfg.Database.Host = dbHost
    }
    
    // ... rest of setup
}
```

### Configuration File

Create `config/production.yaml`:

```yaml
database:
  driver: postgres
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  dbname: ${DB_NAME}
  sslmode: require

server:
  host: 0.0.0.0
  port: "8000"

admin:
  enabled: true
  path: "/admin"

logging:
  level: info
  format: json
```

## Database Setup

### Run Migrations

Always run migrations before starting the application:

```bash
# Set environment variables
export DB_HOST=your-db-host
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=your-db-name

# Run migrations
./myapp migrate
```

### Database Backup

Set up regular backups:

```bash
# Backup script
#!/bin/bash
pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql
```

## Deployment Options

### Docker

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o myapp ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
COPY --from=builder /app/config ./config
CMD ["./myapp"]
```

Build and run:

```bash
docker build -t myapp .
docker run -p 8000:8000 \
  -e DB_HOST=db \
  -e DB_USER=user \
  -e DB_PASSWORD=pass \
  -e DB_NAME=mydb \
  myapp
```

### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:15
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: mydb
    volumes:
      - postgres_data:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      DB_HOST: db
      DB_USER: user
      DB_PASSWORD: pass
      DB_NAME: mydb
    depends_on:
      - db
    command: ./myapp migrate && ./myapp runserver

volumes:
  postgres_data:
```

Run:

```bash
docker-compose up -d
```

### Systemd Service

Create `/etc/systemd/system/myapp.service`:

```ini
[Unit]
Description=My Forge Application
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/myapp
ExecStart=/opt/myapp/myapp runserver
Restart=always
RestartSec=10
Environment="DB_HOST=localhost"
Environment="DB_USER=myapp"
Environment="DB_PASSWORD=secret"
Environment="DB_NAME=myapp_db"

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable myapp
sudo systemctl start myapp
sudo systemctl status myapp
```

### Nginx Reverse Proxy

Create `/etc/nginx/sites-available/myapp`:

```nginx
server {
    listen 80;
    server_name example.com;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable:

```bash
sudo ln -s /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL with Let's Encrypt

Install certbot:

```bash
sudo apt-get install certbot python3-certbot-nginx
```

Get certificate:

```bash
sudo certbot --nginx -d example.com
```

## Monitoring

### Logging

Configure structured logging:

```go
import "github.com/forgego/forge/pkg/logging"

logger, err := logging.NewLogger(false) // production mode
if err != nil {
    log.Fatal(err)
}
defer logger.Sync()

logger.Info("Application started",
    zap.String("version", "1.0.0"),
    zap.String("environment", "production"),
)
```

### Health Checks

Add health check endpoint:

```go
router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    // Check database connection
    if err := db.Ping(); err != nil {
        http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
    })
})
```

### Metrics

Add metrics endpoint:

```go
import "github.com/prometheus/client_golang/prometheus"

router.Get("/metrics", promhttp.Handler())
```

## Performance

### Database Connection Pooling

Configure connection pool:

```go
database.SetMaxOpenConns(25)
database.SetMaxIdleConns(5)
database.SetConnMaxLifetime(5 * time.Minute)
```

### Caching

Add caching for frequently accessed data:

```go
import "github.com/forgego/forge/pkg/cache"

cache := cache.NewRedisCache(redisClient)

// Cache query results
key := fmt.Sprintf("user:%d", userID)
if cached, err := cache.Get(key); err == nil {
    return cached.(*User), nil
}

user, err := User.Objects.Get(ctx, userID)
if err == nil {
    cache.Set(key, user, 5*time.Minute)
}
return user, err
```

## Security

### Environment Variables

Never commit secrets. Use environment variables or secrets management.

### HTTPS

Always use HTTPS in production. Set up SSL certificates.

### Security Headers

Add security headers (see [Security Guide](security)).

### Firewall

Configure firewall to only allow necessary ports:

```bash
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS
sudo ufw enable
```

## Backup and Recovery

### Database Backups

Set up automated backups:

```bash
# Daily backup script
0 2 * * * pg_dump -h $DB_HOST -U $DB_USER $DB_NAME | gzip > /backups/db_$(date +\%Y\%m\%d).sql.gz
```

### Application Backups

Backup application files:

```bash
tar -czf app_backup_$(date +%Y%m%d).tar.gz /opt/myapp
```

## Troubleshooting

### Check Logs

```bash
# Systemd
sudo journalctl -u myapp -f

# Docker
docker logs -f myapp

# Application logs
tail -f /var/log/myapp/app.log
```

### Database Issues

```bash
# Check connection
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1"

# Check migrations
./myapp migrate version
```

### Performance Issues

- Check database query performance
- Monitor connection pool usage
- Review application logs
- Use profiling tools

## Next Steps

- [Security Guide](security) - Secure your deployment
- [Performance Guide](../advanced/performance) - Optimize performance

