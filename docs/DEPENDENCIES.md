# Required Dependencies

## Core Dependencies

The following dependencies need to be added to `go.mod`:

```go
require (
    gopkg.in/natefinch/lumberjack.v2 v2.0.0
)
```

## Optional Dependencies (for future features)

### OpenTelemetry Integration
```go
require (
    go.opentelemetry.io/otel v1.21.0
    go.opentelemetry.io/otel/exporters/jaeger v1.21.0
    go.opentelemetry.io/otel/sdk v1.21.0
    go.opentelemetry.io/otel/trace v1.21.0
)
```

### Redis (for idempotency store)
```go
require (
    github.com/redis/go-redis/v9 v9.3.0
)
```

## Installation

To add the required dependency:

```bash
cd forge
go get gopkg.in/natefinch/lumberjack.v2
go mod tidy
```

## Current Status

- ✅ **zap** - Already in dependencies
- ✅ **google/uuid** - Already in dependencies  
- ⚠️ **lumberjack.v2** - Needs to be added for file rotation
- ⏳ **OpenTelemetry** - Optional, for future observability features
- ⏳ **Redis** - Optional, for distributed idempotency store

## Notes

- The file exporter currently works without lumberjack but lacks rotation
- Once lumberjack is added, file rotation will be automatically enabled
- All other features work without additional dependencies
