# Routing Module

URL routing with reverse lookup and route groups.

## Usage

### Basic Routing

```go
router := routing.NewRouter(app)

router.Get("/", homeHandler, routing.Name("home"))
router.Get("/users/:id", showUser, routing.Name("users.show"))
```

### Route Groups

```go
api := router.Group("/api/v1", middleware.Auth())
api.Get("/users", listUsers, routing.Name("api.users.list"))
api.Get("/users/:id", showUser, routing.Name("api.users.show"))
```

### Reverse URL Generation

```go
url, _ := router.URL("users.show", routing.Param("id", "123"))
// Returns: "/users/123"
```

## Features

- Named routes
- Reverse URL generation
- Route groups with middleware
- Route listing
- Type-safe path parameters

