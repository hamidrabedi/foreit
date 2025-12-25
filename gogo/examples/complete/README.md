# Complete Example

A complete Gogo application demonstrating all modules working together.

## Features

- All 13 modules enabled
- Authentication setup
- Resource registration
- Console (admin) setup
- Workers and scheduler
- Graceful shutdown

## Running

```bash
export DATABASE_URL="postgres://user:pass@localhost/dbname"
export PORT=8080
export SECRET_KEY="your-secret-key"
export DEBUG=false
export REDIS_ADDR="localhost:6379"
export REDIS_PASSWORD=""
export REDIS_DB=0
export WORKERS_CONCURRENCY=10

go run main.go
```

## Structure

This example shows:
- Application initialization
- Module configuration
- Resource registration
- Route setup
- Graceful shutdown

## Next Steps

1. Define Ent schemas
2. Create resource handlers
3. Define auth policies
4. Add translations
5. Customize console

