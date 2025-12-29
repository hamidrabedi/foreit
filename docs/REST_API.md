# REST API Guide

forge provides a comprehensive REST API system similar to Django REST Framework (DRF). This allows you to build React, Vue, or any other frontend that consumes JSON APIs.

## Overview

The REST API is built on the `ViewSet` pattern, providing:

- Full CRUD operations (Create, Read, Update, Delete)
- Pagination
- Filtering and ordering
- Serialization
- Error handling

## Basic Usage

### 1. Define a Serializer

```go
package models

import (
    "github.com/forgego/forge/pkg/api"
)

// UserSerializer serializes User model
type UserSerializer struct {
    *api.BaseSerializer
}

func NewUserSerializer() api.Serializer {
    return &UserSerializer{
        BaseSerializer: api.NewBaseSerializer(),
    }
}

func (s *UserSerializer) Fields() []string {
    return []string{"id", "username", "email", "is_active"}
}
```

### 2. Create a ViewSet

```go
package api

import (
    "your-module/models"
    "github.com/forgego/forge/pkg/api"
)

func RegisterUserViewSet(router *api.Router) {
    viewset := api.NewBaseViewSet(
        func() api.Serializer {
            return models.NewUserSerializer()
        },
        models.User.Objects.Filter(), // QuerySet
        &models.User{},                // Model instance
    )

    router.Register("users", viewset)
}
```

### 3. Register Routes

```go
package main

import (
    "github.com/forgego/forge/pkg/api"
    httplib "github.com/forgego/forge/pkg/http"
)

func main() {
    // ... setup code ...

    // Create API router
    apiRouter := api.NewRouter("/api/v1")

    // Register viewsets
    RegisterUserViewSet(apiRouter)
    RegisterPostViewSet(apiRouter)

    // Register on HTTP router
    server.RegisterRoutes(func(router *httplib.Router) {
        apiRouter.RegisterRoutes(router)
    })
}
```

## API Endpoints

Once registered, your ViewSet automatically provides these endpoints:

### List (GET /api/v1/users/)

Returns paginated list of users.

**Query Parameters:**

- `page` - Page number (default: 1)
- `page_size` - Items per page (default: 20)
- `ordering` - Order by field (e.g., `ordering=-created_at`)
- Field filters (e.g., `is_active=true`)

**Response:**

```json
{
  "count": 100,
  "next": "http://localhost:8000/api/v1/users/?page=2",
  "previous": null,
  "results": [
    {
      "id": 1,
      "username": "john",
      "email": "john@example.com",
      "is_active": true
    }
  ]
}
```

### Create (POST /api/v1/users/)

Create a new user.

**Request Body:**

```json
{
  "username": "jane",
  "email": "jane@example.com",
  "is_active": true
}
```

**Response:**

```json
{
  "id": 2,
  "username": "jane",
  "email": "jane@example.com",
  "is_active": true
}
```

### Retrieve (GET /api/v1/users/{id}/)

Get a single user by ID.

**Response:**

```json
{
  "id": 1,
  "username": "john",
  "email": "john@example.com",
  "is_active": true
}
```

### Update (PUT /api/v1/users/{id}/)

Full update of a user (all fields required).

**Request Body:**

```json
{
  "username": "john_updated",
  "email": "john_updated@example.com",
  "is_active": false
}
```

**Response:**

```json
{
  "id": 1,
  "username": "john_updated",
  "email": "john_updated@example.com",
  "is_active": false
}
```

### Partial Update (PATCH /api/v1/users/{id}/)

Partial update (only provided fields).

**Request Body:**

```json
{
  "is_active": false
}
```

**Response:**

```json
{
  "id": 1,
  "username": "john",
  "email": "john@example.com",
  "is_active": false
}
```

### Destroy (DELETE /api/v1/users/{id}/)

Delete a user.

**Response:** 204 No Content

## Filtering

Filter results using query parameters:

```
GET /api/v1/users/?is_active=true
GET /api/v1/users/?username=john
GET /api/v1/users/?is_active=true&ordering=-created_at
```

## Ordering

Order results using the `ordering` parameter:

```
GET /api/v1/users/?ordering=username          # Ascending
GET /api/v1/users/?ordering=-username         # Descending
GET /api/v1/users/?ordering=username,-created_at  # Multiple fields
```

## Pagination

Pagination is automatic. Use `page` and `page_size` parameters:

```
GET /api/v1/users/?page=2&page_size=50
```

## Error Handling

### Validation Errors (400)

```json
{
  "errors": {
    "email": ["This field is required."],
    "username": ["A user with this username already exists."]
  }
}
```

### Not Found (404)

```json
{
  "error": "User not found"
}
```

### Server Error (500)

```json
{
  "error": "Internal server error"
}
```

## Using with React

### Example: Fetch Users

```typescript
// api/users.ts
export interface User {
  id: number;
  username: string;
  email: string;
  is_active: boolean;
}

export interface PaginatedResponse<T> {
  count: number;
  next: string | null;
  previous: string | null;
  results: T[];
}

export async function getUsers(
  page = 1,
  pageSize = 20
): Promise<PaginatedResponse<User>> {
  const response = await fetch(
    `/api/v1/users/?page=${page}&page_size=${pageSize}`
  );
  if (!response.ok) {
    throw new Error("Failed to fetch users");
  }
  return response.json();
}

export async function createUser(user: Partial<User>): Promise<User> {
  const response = await fetch("/api/v1/users/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(user),
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to create user");
  }
  return response.json();
}

export async function updateUser(
  id: number,
  user: Partial<User>
): Promise<User> {
  const response = await fetch(`/api/v1/users/${id}/`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(user),
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to update user");
  }
  return response.json();
}

export async function deleteUser(id: number): Promise<void> {
  const response = await fetch(`/api/v1/users/${id}/`, {
    method: "DELETE",
  });
  if (!response.ok) {
    throw new Error("Failed to delete user");
  }
}
```

### Example: React Component

```tsx
import { useState, useEffect } from "react";
import { getUsers, createUser, User } from "./api/users";

function UserList() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  useEffect(() => {
    loadUsers();
  }, [page]);

  async function loadUsers() {
    setLoading(true);
    try {
      const data = await getUsers(page);
      setUsers(data.results);
    } catch (error) {
      console.error("Failed to load users:", error);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate(user: Partial<User>) {
    try {
      const newUser = await createUser(user);
      setUsers([...users, newUser]);
    } catch (error) {
      console.error("Failed to create user:", error);
    }
  }

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Users</h1>
      <ul>
        {users.map((user) => (
          <li key={user.id}>
            {user.username} - {user.email}
          </li>
        ))}
      </ul>
      <button onClick={() => setPage(page + 1)}>Next Page</button>
    </div>
  );
}
```

## Using with Vue

### Example: Vue Component

```vue
<template>
  <div>
    <h1>Users</h1>
    <ul v-if="!loading">
      <li v-for="user in users" :key="user.id">
        {{ user.username }} - {{ user.email }}
      </li>
    </ul>
    <div v-else>Loading...</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";

const users = ref([]);
const loading = ref(true);
const page = ref(1);

async function loadUsers() {
  loading.value = true;
  try {
    const response = await fetch(`/api/v1/users/?page=${page.value}`);
    const data = await response.json();
    users.value = data.results;
  } catch (error) {
    console.error("Failed to load users:", error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadUsers();
});
</script>
```

## Authentication

For authenticated APIs, add middleware:

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            forgehttp.SendError(w, http.StatusUnauthorized, "Authentication required")
            return
        }
        // Validate token...
        next.ServeHTTP(w, r)
    })
}

// Apply to API routes
apiRouter.RegisterRoutes(router)
router.Use(AuthMiddleware)
```

## CORS

For cross-origin requests, add CORS middleware:

```go
import "github.com/go-chi/cors"

router.Use(cors.Handler(cors.Options{
    AllowedOrigins: []string{"http://localhost:3000"},
    AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
    AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
}))
```

## Best Practices

1. **Use Serializers** - Control which fields are exposed
2. **Validate Input** - Use serializer validation
3. **Handle Errors** - Return proper HTTP status codes
4. **Use Pagination** - Always paginate large datasets
5. **Filter and Order** - Use query parameters for flexibility
6. **Version APIs** - Use `/api/v1/` prefix for versioning

## Comparison: HTMX Admin vs REST API

| Feature         | HTMX Admin      | REST API                |
| --------------- | --------------- | ----------------------- |
| Use Case        | Admin interface | Frontend apps, mobile   |
| Response Format | HTML            | JSON                    |
| Interactivity   | Server-rendered | Client-side             |
| Complexity      | Low             | Medium                  |
| Build Tools     | None            | Optional (for frontend) |
| Best For        | Internal tools  | Public APIs, SPAs       |

## Next Steps

- See [HTMX Patterns](HTMX_PATTERNS.md) for admin interface patterns
- See [API Reference](API_REFERENCE.md) for detailed API documentation
- See examples in `examples/` directory
