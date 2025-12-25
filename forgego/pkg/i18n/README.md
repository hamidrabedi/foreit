# i18n Module

Internationalization support for Fiber applications.

## Usage

### Load Translations

```go
// Load from directory (e.g., locales/en.json, locales/fr.json)
i18n.Load("locales", "en")
```

### Translation Files

`locales/en.json`:
```json
{
  "welcome": "Welcome",
  "hello": "Hello, {0}",
  "items_count": "{0} items"
}
```

`locales/fr.json`:
```json
{
  "welcome": "Bienvenue",
  "hello": "Bonjour, {0}",
  "items_count": "{0} éléments"
}
```

### Use in Handlers

```go
app.Use(i18n.Middleware())

func handler(c *fiber.Ctx) error {
    message := i18n.T(c, "welcome")
    greeting := i18n.T(c, "hello", "John")
    return c.JSON(fiber.Map{"message": message, "greeting": greeting})
}
```

## Features

- JSON-based translations
- Locale detection (Accept-Language header, query param, session)
- Fallback to default locale
- String formatting with arguments
- Middleware integration

