package models

import (
	"context"
)

// Example: Type-safe relationships

// User model
type User struct {
	ID    int
	Email string
	Name  string
}

func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(int) }
func (u *User) IsNew() bool { return u.ID == 0 }
func (u *User) String() string { return u.Email }

// Post model
type Post struct {
	ID     int
	UserID int  // Foreign key
	Title  string
	Body   string
}

func (p *Post) GetID() interface{} { return p.ID }
func (p *Post) SetID(id interface{}) { p.ID = id.(int) }
func (p *Post) IsNew() bool { return p.ID == 0 }
func (p *Post) String() string { return p.Title }

// Type-safe field references
var (
	UserID   = NewFieldRef[int]("id")
	UserEmail = NewFieldRef[string]("email")
	
	PostID    = NewFieldRef[int]("id")
	PostUserID = NewFieldRef[int]("user_id")
	PostTitle  = NewFieldRef[string]("title")
)

// Example: One-to-Many (User has many Posts)
func ExampleOneToMany(
	ctx context.Context,
	userManager Manager[*User],
	postManager Manager[*Post],
) {
	// Create type-safe relationship
	userPosts := NewOneToMany[*User, *Post](
		UserID,
		PostUserID,
		postManager,
		"user_id",
	)
	
	// Get a user
	user, _ := userManager.Get(ctx, 1)
	
	// Load user's posts - returns []*Post, type-safe!
	posts, err := userPosts.Load(ctx, user)
	
	// posts is []*Post - no type assertions!
	for _, post := range posts {
		_ = post.Title  // Direct access
		_ = post.Body
	}
	
	// Eager load posts for multiple users
	users, _ := userManager.All(ctx)
	postsByUser, _ := EagerLoad(ctx, users, userPosts)
	
	// postsByUser is map[interface{}][]*Post
	// Access posts for a user
	userPostsList := postsByUser[user.GetID()]
	_ = userPostsList[0].Title  // Type-safe!
}

// Example: Many-to-One (Post belongs to User)
func ExampleManyToOne(
	ctx context.Context,
	userManager Manager[*User],
	postManager Manager[*Post],
) {
	// Create type-safe relationship
	postAuthor := NewManyToOne[*Post, *User](
		PostUserID,
		UserID,
		userManager,
		"user_id",
	)
	
	// Get a post
	post, _ := postManager.Get(ctx, 1)
	
	// Load post's author - returns *User, type-safe!
	author, err := postAuthor.Load(ctx, post)
	
	// author is *User - no type assertions!
	_ = author.Email
	_ = author.Name
}

// Example: One-to-One
func ExampleOneToOne(
	ctx context.Context,
	userManager Manager[*User],
	profileManager Manager[*UserProfile],
) {
	// Create type-safe relationship
	userProfile := NewOneToOne[*User, *UserProfile](
		UserID,
		NewFieldRef[*UserProfile]("id"),
		profileManager,
		"user_id",
		true, // Foreign key is on profile
	)
	
	// Get a user
	user, _ := userManager.Get(ctx, 1)
	
	// Load user's profile - returns *UserProfile, type-safe!
	profile, err := userProfile.Load(ctx, user)
	
	// profile is *UserProfile - no type assertions!
	_ = profile.Bio
}

// UserProfile model (example for one-to-one)
type UserProfile struct {
	ID     int
	UserID int
	Bio    string
}

func (p *UserProfile) GetID() interface{} { return p.ID }
func (p *UserProfile) SetID(id interface{}) { p.ID = id.(int) }
func (p *UserProfile) IsNew() bool { return p.ID == 0 }
func (p *UserProfile) String() string { return p.Bio }

