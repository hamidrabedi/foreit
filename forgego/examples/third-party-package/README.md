# Third-Party Package Integration Example

This example demonstrates how third-party packages can integrate with Forge by:

1. Defining their own config struct in their package
2. Users embedding the config in their Settings struct
3. Loading and using the config

## Structure

```
examples/third-party-package/
├── main.go              # Application entry point
├── config.yaml          # Configuration file
├── mypackage/           # Example third-party package
│   ├── config.go        # Package config definition
│   └── package.go       # Package implementation
└── README.md            # This file
```

## How It Works

### 1. Package Defines Its Config

```go
// mypackage/config.go
type MyPackageConfig struct {
    Enable bool   `mapstructure:"enable" default:"true"`
    Path   string `mapstructure:"path" default:"/mypackage"`
    APIKey string `mapstructure:"api_key" validate:"required"`
}
```

### 2. User Embeds in Settings

```go
// main.go
type Settings struct {
    app.AppSettings
    MyPackage mypackage.MyPackageConfig `mapstructure:"mypackage"`
}
```

### 3. Load and Use

```go
settings, _ := settings.Load[Settings]()
myPackage := mypackage.New(&settings.MyPackage)
```

## Benefits

- **Package ownership**: Each package owns its config
- **Type safety**: Compile-time checking
- **Composability**: Users choose which packages to include
- **No coupling**: Packages don't depend on framework settings

## Running

```bash
go run main.go
```

The config is loaded from `config.yaml` or environment variables (e.g., `MYPACKAGE_ENABLE`, `MYPACKAGE_PATH`, `MYPACKAGE_API_KEY`).

