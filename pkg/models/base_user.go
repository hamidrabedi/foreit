package models

import (
	"context"

	"github.com/forgego/forge/pkg/schema"
	field "github.com/forgego/forge/pkg/schema"
)

// BaseUser is a Django-like base user model that can be extended
type BaseUser struct {
	schema.Schema
}

// Fields returns the field definitions for BaseUser
func (BaseUser) Fields() []schema.Field {
	return []schema.Field{
		field.Int64("id").Primary().AutoIncrement().Build(),
		field.String("username").Unique().Required().MaxLength(150).VerboseName("Username").Build(),
		field.String("email").Unique().Required().MaxLength(254).VerboseName("Email Address").Build(),
		field.String("password").Required().MaxLength(128).VerboseName("Password").Build(),
		field.String("first_name").Optional().MaxLength(150).VerboseName("First Name").Build(),
		field.String("last_name").Optional().MaxLength(150).VerboseName("Last Name").Build(),
		field.Bool("is_active").Default(true).VerboseName("Active").Build(),
		field.Bool("is_staff").Default(false).VerboseName("Staff Status").Build(),
		field.Bool("is_superuser").Default(false).VerboseName("Superuser Status").Build(),
		field.Time("date_joined").AutoNowAdd().VerboseName("Date Joined").Build(),
		field.Time("last_login").Optional().VerboseName("Last Login").Build(),
		field.Time("created_at").AutoNowAdd().VerboseName("Created At").Build(),
		field.Time("updated_at").AutoNow().VerboseName("Updated At").Build(),
	}
}

// Relations returns the relationship definitions for BaseUser
func (BaseUser) Relations() []schema.Relation {
	return []schema.Relation{
		// User has many Posts (example - can be extended)
		// field.OneToMany("posts", "Post", "author_id").CascadeOnDelete(),
	}
}

// Meta returns the model metadata for BaseUser
func (BaseUser) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "users",
		VerboseName:       "User",
		VerboseNamePlural: "Users",
		OrderBy:           []string{"-date_joined"},
		Indexes: []schema.Index{
			{Name: "idx_users_email", Fields: []string{"email"}, Unique: false},
			{Name: "idx_users_username", Fields: []string{"username"}, Unique: false},
			{Name: "idx_users_is_active", Fields: []string{"is_active"}, Unique: false},
		},
	}
}

// Hooks returns the lifecycle hooks for BaseUser
func (BaseUser) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Hash password before creating user
			// This will be implemented when password hashing is integrated
			return nil
		},
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Send welcome email after creating user
			// This will be implemented when email system is integrated
			return nil
		},
		BeforeUpdate: func(ctx context.Context, instance interface{}) error {
			// Update last_login if it's a login operation
			// This will be implemented when authentication is integrated
			return nil
		},
	}
}

// UserManager provides methods for user management
type UserManager interface {
	// CreateUser creates a new user with hashed password
	CreateUser(ctx interface{}, username, email, password string) (interface{}, error)

	// CreateSuperuser creates a new superuser
	CreateSuperuser(ctx interface{}, username, email, password string) (interface{}, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx interface{}, email string) (interface{}, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx interface{}, username string) (interface{}, error)

	// Authenticate authenticates a user with username/email and password
	Authenticate(ctx interface{}, usernameOrEmail, password string) (interface{}, error)

	// SetPassword sets a user's password (hashed)
	SetPassword(ctx, user interface{}, password string) error

	// CheckPassword checks if a password matches
	CheckPassword(ctx, user interface{}, password string) bool
}

// UserPermissions provides permission checking methods
type UserPermissions interface {
	// HasPerm checks if user has a specific permission
	HasPerm(user interface{}, perm string) bool

	// HasPerms checks if user has all specified permissions
	HasPerms(user interface{}, perms []string) bool

	// HasAnyPerm checks if user has any of the specified permissions
	HasAnyPerm(user interface{}, perms []string) bool

	// IsActive checks if user is active
	IsActive(user interface{}) bool

	// IsStaff checks if user is staff
	IsStaff(user interface{}) bool

	// IsSuperuser checks if user is superuser
	IsSuperuser(user interface{}) bool
}
