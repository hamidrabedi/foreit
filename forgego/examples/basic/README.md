# Basic Example

This example demonstrates **automatic initialization** - the framework automatically sets up enabled features based on your settings.

## Key Principles

1. **Automatic initialization** - Framework features are initialized automatically based on settings
2. **Error handling** - If initialization fails, `app.New()` returns an error
3. **Configuration-driven** - Enable/disable features via settings
4. **No manual setup** - Users don't need to manually initialize framework features

## Structure

```
examples/basic/
├── main.go          # Simple application - framework handles initialization
├── config.yaml      # Configuration file
└── README.md        # This file
```

## Running

1. Make sure you have a database running (if using database features)
2. Update `config.yaml` with your settings
3. Run the application:

```bash
go run main.go
```

## Automatic Initialization

The framework automatically initializes enabled features:

```go
// Load settings
settings, _ := settings.Load[Settings]()

// Create app - framework automatically:
// - Sets up middleware
// - Connects to database (if database.url is set)
// - Sets up API router (if api.path is set)
// - Initializes i18n (if i18n.enable=true)
// - Sets up static files (if static.enable=true)
// - Initializes workers (if workers.enable=true)
// - Sets up admin (if admin.enable=true and database is connected)
app, err := app.New(&settings.AppSettings)
if err != nil {
    // Handle initialization errors
    log.Fatal(err)
}
```

## What Gets Initialized

Based on your `config.yaml`:

- **Middleware**: Always initialized (logging, recovery, CORS, etc.)
- **Database**: If `database.url` is set
- **API Router**: If `api.path` is set
- **i18n**: If `i18n.enable=true` and `i18n.locales_path` is set
- **Static Files**: If `static.enable=true` and `static.root` is set
- **Workers**: If `workers.enable=true` and Redis config is valid
- **Admin**: If `admin.enable=true` and database is connected

## Error Handling

If any required initialization fails, `app.New()` returns an error:

```go
app, err := app.New(settings)
if err != nil {
    // Error messages are descriptive:
    // - "database connection failed: ..."
    // - "i18n setup failed: ..."
    // - "workers setup failed: ..."
    // - "admin setup failed: database connection required..."
    log.Fatal(err)
}
```

## Configuration

Settings can be provided via:

1. **Config file** (`config.yaml` in the same directory)
2. **Environment variables** (e.g., `DATABASE_URL`, `SERVER_PORT`, `MYAPP_API_KEY`)

Environment variables take precedence over config file values.

## Benefits

- **Simple**: Just configure and create the app
- **Safe**: Errors are caught at startup
- **Clear**: Descriptive error messages
- **Flexible**: Enable only what you need via settings
