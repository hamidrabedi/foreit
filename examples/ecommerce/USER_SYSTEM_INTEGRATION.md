# User System Integration Guide

The user system has been successfully integrated into the ecommerce project!

## What's Been Added

### 1. Database Migrations
- **File**: `migrations/000002_create_users_tables.up.sql`
- **File**: `migrations/000002_create_users_tables.down.sql`
- Creates all necessary tables for the user system:
  - `users` - Main user table
  - `user_sessions` - Session management
  - `permissions` - Permission definitions
  - `groups` - User groups
  - `user_permissions`, `user_groups`, `group_permissions` - Many-to-many tables
  - `email_verification_tokens` - Email verification
  - `password_reset_tokens` - Password reset

### 2. Code Integration
- User system initialized in `main.go`
- Routes automatically registered
- Full authentication and authorization support

## API Endpoints

### Authentication Endpoints
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout (requires authentication)
- `GET /api/auth/me` - Get current user (requires authentication)

### User Management Endpoints
- `GET /api/users/` - List users (requires authentication)
- `POST /api/users/` - Create user (requires superuser)
- `GET /api/users/{id}` - Get user (requires authentication)
- `PUT /api/users/{id}` - Update user (requires authentication)
- `DELETE /api/users/{id}` - Delete user (requires superuser)

## Setup Instructions

### 1. Run Migrations

```bash
cd examples/ecommerce
forge migrate up
```

This will create all user system tables.

### 2. Start the Server

```bash
go run cmd/server/main.go
```

The server will start with user system enabled.

### 3. Test the API

#### Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecurePassword123!"
  }'
```

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecurePassword123!"
  }'
```

Response will include a `session_key` that you can use for authenticated requests.

#### Get Current User
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "X-Session-Key: YOUR_SESSION_KEY"
```

## Integration with Ecommerce Models

The user system is separate from the ecommerce `Customer` model. You can:

1. **Link customers to users**: Add a `user_id` foreign key to the `customers` table
2. **Use user system for authentication**: Use user system for login/auth, then link to customer profile
3. **Use both independently**: Keep them separate if needed

### Example: Link Customer to User

```sql
-- Add user_id to customers table
ALTER TABLE customers ADD COLUMN user_id BIGINT REFERENCES users(id);
CREATE INDEX idx_customers_user_id ON customers(user_id);
```

Then in your code:
```go
// After user registration
user, _ := userSystem.UserService.Register(ctx, &service.RegisterRequest{
    Username: "johndoe",
    Email: "john@example.com",
    Password: "SecurePassword123!",
})

// Create customer linked to user
customer := &models.Customer{
    UserID: user.ID,
    Email: user.Email,
    // ... other fields
}
```

## Configuration

The user system uses default configuration. To customize:

```go
import "github.com/forgego/forge/pkg/users/config"

userConfig := config.DefaultUserConfig()
userConfig.PasswordPolicy.MinLength = 12
userConfig.LockoutConfig.MaxFailedAttempts = 5

userSystem, err := users.SetupUserSystem(database, userConfig)
```

## Security Features

- ✅ Password hashing (bcrypt)
- ✅ Session management
- ✅ Account lockout support
- ✅ Permission system (RBAC)
- ✅ Group management
- ✅ Email verification support
- ✅ Password reset support

## Next Steps

1. **Run migrations** to create tables
2. **Test the API endpoints** using curl or Postman
3. **Link customers to users** if needed
4. **Add email service** for email verification
5. **Customize configuration** as needed

## Documentation

For more details, see:
- `docs/USER_SYSTEM_ROADMAP.md` - Complete feature roadmap
- `docs/USER_SYSTEM_ARCHITECTURE.md` - Architecture details
- `docs/USER_SYSTEM_QUICK_START.md` - Quick start guide
- `pkg/users/README.md` - Package documentation
