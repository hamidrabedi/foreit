# Modular Architecture Example

This example demonstrates how to use Gogo's modular architecture.

## Features Demonstrated

- Settings loading
- ORM client setup
- Pipeline middleware
- Sessions
- Authentication
- API endpoints
- Console (admin)
- Background workers
- Job scheduling

## Running

```bash
export DATABASE_URL="postgres://user:pass@localhost/dbname"
export PORT=8080
export SECRET_KEY="your-secret-key"

go run main.go
```

## Modules Used

- `pkg/settings` - Configuration
- `pkg/orm` - Database
- `pkg/pipeline` - Middleware
- `pkg/sessions` - Sessions
- `pkg/auth` - Authentication
- `pkg/endpoints` - API
- `pkg/console` - Admin
- `pkg/workers` - Background jobs
- `pkg/routing` - URL routing

