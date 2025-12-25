package models

import (
	"context"
	"time"
)

// Example: Type-safe model definition and usage

// User model struct
type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	IsActive  bool      `db:"is_active"`
}

// Implement Model interface
func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(int) }
func (u *User) IsNew() bool { return u.ID == 0 }
func (u *User) String() string { return u.Email }

// Type-safe field references for queries
var (
	UserEmail     = NewFieldRef[string]("email")
	UserName      = NewFieldRef[string]("name")
	UserID        = NewFieldRef[int]("id")
	UserCreatedAt = NewFieldRef[time.Time]("created_at")
	UserIsActive  = NewFieldRef[bool]("is_active")
)

// Type-safe model definition
var UserModel = NewModelBuilder[*User]("User").
	String("email").Required().Unique().
	String("name").Required().
	Int("id").Required().Indexed().
	Time("created_at").Default(time.Now()).
	Bool("is_active").Default(true).
	Build()

// Example usage with type-safe manager
func ExampleUsage(ctx context.Context, userManager Manager[*User]) {
	// Type-safe queries - no type assertions needed!
	users, err := userManager.Filter(ctx).
		Filter(UserEmail.Contains("@example.com")).
		Filter(UserIsActive.Eq(true)).
		OrderBy("-created_at").
		Limit(10).
		All(ctx)
	
	// users is []*User - fully type-safe!
	for _, user := range users {
		_ = user.Email // Direct access, no type assertion
		_ = user.Name
	}
	
	// Get single user - type-safe
	user, err := userManager.Filter(ctx).
		Filter(UserEmail.Eq("john@example.com")).
		Get(ctx)
	
	// user is *User - type-safe!
	_ = user.Email
	
	// Create user - type-safe
	newUser := &User{
		Email:    "jane@example.com",
		Name:     "Jane",
		IsActive: true,
	}
	err = userManager.Create(ctx, newUser)
	
	// Update user - type-safe
	user.Name = "John Updated"
	err = userManager.Update(ctx, user)
	
	// Delete user - type-safe
	err = userManager.Delete(ctx, user)
}

// Type-safe with Q objects
func ExampleWithQ(ctx context.Context, userManager Manager[*User]) {
	// Using Q map directly (still type-safe field refs)
	users, err := userManager.Filter(ctx).
		Filter(Q{
			"email__icontains": "@example.com",
			"is_active": true,
			"created_at__gte": time.Now().AddDate(0, -1, 0),
		}).
		Exclude(Q{"name__startswith": "Test"}).
		All(ctx)
	
	// Still type-safe result!
	_ = users[0].Email
}

