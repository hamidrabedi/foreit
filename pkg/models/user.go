package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/forgego/forge/pkg/auth"
	"github.com/forgego/forge/pkg/config"
	"github.com/forgego/forge/pkg/db"
	"github.com/forgego/forge/pkg/query"
	"github.com/forgego/forge/pkg/schema"
)

// User represents a user in the system
type User struct {
	schema.BaseSchema
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username" validate:"required,max=150"`
	Email        string    `json:"email" db:"email" validate:"required,email,max=254"`
	Password     string    `json:"-" db:"password" validate:"required,max=128"` // Never serialize password
	FirstName    string    `json:"first_name" db:"first_name" validate:"max=150"`
	LastName     string    `json:"last_name" db:"last_name" validate:"max=150"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	IsStaff      bool      `json:"is_staff" db:"is_staff"`
	IsSuperuser  bool      `json:"is_superuser" db:"is_superuser"`
	DateJoined   time.Time `json:"date_joined" db:"date_joined"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// GetID returns the user's ID (implements ModelWithID)
func (u *User) GetID() int64 {
	return u.ID
}

// SetID sets the user's ID (implements ModelWithID)
func (u *User) SetID(id int64) {
	u.ID = id
}

// UserManagerImpl provides methods for user management
type UserManagerImpl struct {
	manager *query.Manager[User]
	db      *db.DB
}

// NewUserManager creates a new user manager
func NewUserManager(database *db.DB) *UserManagerImpl {
	manager := query.NewManager[User]("users")
	manager.SetDB(database)
	return &UserManagerImpl{
		manager: manager,
		db:      database,
	}
}

// Get retrieves a user by ID
func (um *UserManagerImpl) Get(ctx context.Context, id int64) (*User, error) {
	return um.manager.Get(ctx, id)
}

// CreateUser creates a new user with hashed password
func (um *UserManagerImpl) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	// Validate inputs
	if username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	// Check if username already exists
	existingUser, _ := um.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with username %s already exists", username)
	}

	// Check if email already exists
	existingUser, _ = um.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		IsActive:   true,
		IsStaff:    false,
		IsSuperuser: false,
		DateJoined: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Use reflection to set fields for insert
	if err := um.createUserWithReflection(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// CreateSuperuser creates a new superuser
func (um *UserManagerImpl) CreateSuperuser(ctx context.Context, username, email, password string) (*User, error) {
	user, err := um.CreateUser(ctx, username, email, password)
	if err != nil {
		return nil, err
	}

	// Set superuser flags
	user.IsStaff = true
	user.IsSuperuser = true

	// Update in database
	if err := um.manager.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update superuser flags: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (um *UserManagerImpl) GetByEmail(ctx context.Context, email string) (*User, error) {
	users, err := um.manager.Filter(query.NewFieldQueryExpr("email", query.OpEquals, email)).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

// GetByUsername retrieves a user by username
func (um *UserManagerImpl) GetByUsername(ctx context.Context, username string) (*User, error) {
	users, err := um.manager.Filter(query.NewFieldQueryExpr("username", query.OpEquals, username)).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("user not found")
	}
	return users[0], nil
}

// Authenticate authenticates a user with username/email and password
func (um *UserManagerImpl) Authenticate(ctx context.Context, usernameOrEmail, password string) (*User, error) {
	// Try to find user by username or email
	var user *User
	var err error

	// Check if it looks like an email
	if strings.Contains(usernameOrEmail, "@") {
		user, err = um.GetByEmail(ctx, usernameOrEmail)
	} else {
		user, err = um.GetByUsername(ctx, usernameOrEmail)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Check password
	if !auth.CheckPassword(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last_login
	now := time.Now()
	user.LastLogin = &now
	if err := um.manager.Update(ctx, user); err != nil {
		// Log error but don't fail authentication
		fmt.Printf("Warning: failed to update last_login: %v\n", err)
	}

	return user, nil
}

// SetPassword sets a user's password (hashed)
func (um *UserManagerImpl) SetPassword(ctx context.Context, user *User, password string) error {
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	return um.manager.Update(ctx, user)
}

// CheckPassword checks if a password matches
func (um *UserManagerImpl) CheckPassword(user *User, password string) bool {
	return auth.CheckPassword(password, user.Password)
}

// createUserWithReflection creates a user using direct SQL since the Manager might not handle all fields
func (um *UserManagerImpl) createUserWithReflection(ctx context.Context, user *User) error {
	// Build INSERT statement
	fields := []string{"username", "email", "password", "first_name", "last_name", "is_active", "is_staff", "is_superuser", "date_joined", "created_at", "updated_at"}
	placeholders := []string{}
	values := []interface{}{}

	userValue := reflect.ValueOf(user).Elem()
	for _, field := range fields {
		fieldValue := userValue.FieldByName(getFieldName(field))
		if fieldValue.IsValid() {
			placeholders = append(placeholders, "?")
			values = append(values, fieldValue.Interface())
		}
	}

	// Adjust for PostgreSQL vs SQLite
	sqlQuery := fmt.Sprintf(
		"INSERT INTO users (%s) VALUES (%s) RETURNING id",
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	// Detect SQLite by checking config
	cfg := config.NewConfig()
	driver := cfg.GetDriver()
	isSQLite := driver == "sqlite" || driver == "sqlite3"
	
	if isSQLite {
		// SQLite syntax
		sqlQuery = fmt.Sprintf(
			"INSERT INTO users (%s) VALUES (%s)",
			strings.Join(fields, ", "),
			strings.Join(placeholders, ", "),
		)
		var result sql.Result
		var err error
		if ctx != nil {
			result, err = um.db.ExecContext(ctx, sqlQuery, values...)
		} else {
			result, err = um.db.Exec(sqlQuery, values...)
		}
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}
		user.ID = id
		return nil
	}

	// PostgreSQL
	var id int64
	err := um.db.QueryRowContext(ctx, sqlQuery, values...).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	user.ID = id
	return nil
}

// getFieldName converts snake_case to PascalCase
func getFieldName(snake string) string {
	parts := strings.Split(snake, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

