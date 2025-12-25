package models

import (
	"context"
	"time"
)

// Complete example showing type-safe usage

// User model
type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
}

func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(int) }
func (u *User) IsNew() bool { return u.ID == 0 }
func (u *User) String() string { return u.Email }

// Post model
type Post struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *Post) GetID() interface{} { return p.ID }
func (p *Post) SetID(id interface{}) { p.ID = id.(int) }
func (p *Post) IsNew() bool { return p.ID == 0 }
func (p *Post) String() string { return p.Title }

// Type-safe field references
var (
	// User fields
	UserID        = NewFieldRef[int]("id")
	UserEmail     = NewFieldRef[string]("email")
	UserName      = NewFieldRef[string]("name")
	UserRole      = NewFieldRef[string]("role")
	UserIsActive  = NewFieldRef[bool]("is_active")
	UserCreatedAt = NewFieldRef[time.Time]("created_at")
	
	// Post fields
	PostID        = NewFieldRef[int]("id")
	PostUserID   = NewFieldRef[int]("user_id")
	PostTitle     = NewFieldRef[string]("title")
	PostBody      = NewFieldRef[string]("body")
	PostCreatedAt = NewFieldRef[time.Time]("created_at")
)

// Complete example function
func CompleteExample(
	ctx context.Context,
	userManager Manager[*User],
	postManager Manager[*Post],
	session Session,
) {
	// 1. Type-safe queries
	users, err := userManager.Filter(ctx).
		Filter(UserEmail.Contains("@example.com")).
		Filter(UserIsActive.Eq(true)).
		Filter(UserRole.Eq("admin")).
		OrderBy("-created_at").
		Limit(10).
		All(ctx)
	
	// users is []*User - type-safe!
	for _, user := range users {
		_ = user.Email
		_ = user.Name
	}
	
	// 2. Type-safe get
	user, err := userManager.Filter(ctx).
		Filter(UserEmail.Eq("john@example.com")).
		Get(ctx)
	
	// user is *User - type-safe!
	_ = user.Email
	
	// 3. Type-safe create
	newUser := &User{
		Email:    "jane@example.com",
		Name:     "Jane",
		Role:     "user",
		IsActive: true,
	}
	err = userManager.Create(ctx, newUser)
	
	// 4. Type-safe update
	user.Name = "John Updated"
	err = userManager.Update(ctx, user)
	
	// 5. Type-safe delete
	err = userManager.Delete(ctx, user)
	
	// 6. Type-safe relationships
	userPosts := NewOneToMany[*User, *Post](
		UserID,
		PostUserID,
		postManager,
		"user_id",
	)
	
	posts, err := userPosts.Load(ctx, user)
	// posts is []*Post - type-safe!
	for _, post := range posts {
		_ = post.Title
	}
	
	// 7. Eager loading
	users, _ = userManager.All(ctx)
	postsByUser, _ := EagerLoad(ctx, users, userPosts)
	// postsByUser is map[interface{}][]*Post
	
	// 8. Transactions (SQLAlchemy-style)
	result, err := Transactional(ctx, session, func(txCtx context.Context) ([]*User, error) {
		// Create multiple users in transaction
		users := []*User{
			{Email: "user1@example.com", Name: "User 1"},
			{Email: "user2@example.com", Name: "User 2"},
		}
		
		for _, u := range users {
			if err := userManager.Create(txCtx, u); err != nil {
				return nil, err
			}
		}
		
		return users, nil
	})
	
	// result is []*User - type-safe!
	_ = result
	
	// 9. Aggregations
	count, _ := userManager.Filter(ctx).Count(ctx)
	_ = count
	
	// 10. Complex Q objects
	q := Q{
		"email__icontains": "@example.com",
		"is_active": true,
		"created_at__gte": time.Now().AddDate(0, -1, 0),
	}
	
	users, _ = userManager.Filter(ctx).
		Filter(q).
		Exclude(Q{"role": "guest"}).
		All(ctx)
}

