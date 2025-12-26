---
sidebar_position: 2
---

# Plugins

forge has a plugin system that allows you to extend and customize framework behavior.

## Plugin Interface

Plugins implement the `Plugin` interface:

```go
type Plugin interface {
    Name() string
    Initialize(app *Application) error
}
```

## Creating a Plugin

### Basic Plugin

```go
package plugins

import (
    "github.com/forgego/forge/pkg/registry"
)

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Initialize(app *registry.Application) error {
    // Plugin initialization logic
    return nil
}
```

### Registering a Plugin

```go
import (
    "github.com/forgego/forge/pkg/registry"
    "myapp/plugins"
)

func main() {
    app := registry.NewApplication()
    
    // Register plugin
    app.RegisterPlugin(&plugins.MyPlugin{})
    
    // Initialize all plugins
    if err := app.Initialize(); err != nil {
        log.Fatal(err)
    }
}
```

## Plugin Hooks

Plugins can hook into various framework events:

### Model Registration

```go
func (p *MyPlugin) Initialize(app *registry.Application) error {
    app.OnModelRegister(func(model schema.Schema) {
        // Called when a model is registered
        log.Printf("Model registered: %T", model)
    })
    return nil
}
```

### Before Request

```go
func (p *MyPlugin) Initialize(app *registry.Application) error {
    app.OnBeforeRequest(func(w http.ResponseWriter, r *http.Request) {
        // Called before each request
    })
    return nil
}
```

### After Request

```go
func (p *MyPlugin) Initialize(app *registry.Application) error {
    app.OnAfterRequest(func(w http.ResponseWriter, r *http.Request) {
        // Called after each request
    })
    return nil
}
```

## Example Plugins

### Logging Plugin

```go
type LoggingPlugin struct{}

func (p *LoggingPlugin) Name() string {
    return "logging"
}

func (p *LoggingPlugin) Initialize(app *registry.Application) error {
    app.OnBeforeRequest(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s", r.Method, r.URL.Path)
    })
    return nil
}
```

### Caching Plugin

```go
type CachingPlugin struct {
    cache cache.Cache
}

func (p *CachingPlugin) Name() string {
    return "caching"
}

func (p *CachingPlugin) Initialize(app *registry.Application) error {
    p.cache = cache.NewRedisCache(redisClient)
    
    app.OnModelRegister(func(model schema.Schema) {
        // Add caching to models
    })
    
    return nil
}
```

## Plugin Configuration

Plugins can have configuration:

```go
type MyPlugin struct {
    config Config
}

type Config struct {
    Enabled bool
    Option  string
}

func (p *MyPlugin) Initialize(app *registry.Application) error {
    // Use configuration
    if !p.config.Enabled {
        return nil
    }
    
    // Plugin logic
    return nil
}
```

## Best Practices

1. **Keep Plugins Focused** - Each plugin should do one thing well
2. **Use Hooks** - Use framework hooks instead of modifying core code
3. **Handle Errors** - Always return errors from Initialize
4. **Document Plugins** - Document what your plugin does
5. **Test Plugins** - Write tests for your plugins

## See Also

- [Code Generation](code-generation) - Code generation system
- [Custom Fields](custom-fields) - Create custom field types

