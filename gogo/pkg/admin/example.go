package admin

import (
	"context"
	"time"
	"github.com/gogo/pkg/models"
	"gorm.io/gorm"
)

// Example: Admin with type-safe queries

// User model with GORM tags
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"not null" json:"name"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Implement Model interface
func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(uint) }
func (u *User) IsNew() bool { return u.ID == 0 }
func (u *User) String() string { return u.Email }
func (u *User) GetCreatedAt() *time.Time { return &u.CreatedAt }
func (u *User) GetUpdatedAt() *time.Time { return &u.UpdatedAt }

// Register model and create field references
func init() {
	models.RegisterModel[User](map[string]string{
		"ID":        "id",
		"Email":     "email",
		"Name":      "name",
		"IsActive":  "is_active",
		"CreatedAt": "created_at",
		"UpdatedAt": "updated_at",
	})
}

// Type-safe field references (from models package)
var (
	UserEmail    = models.NewStringFieldRef[User]("Email")
	UserName     = models.NewStringFieldRef[User]("Name")
	UserID       = models.NewFieldRef[uint, User]("ID")
	UserIsActive = models.NewFieldRef[bool, User]("IsActive")
)

// Admin field references (for display/search)
var (
	AdminUserEmail    = NewFieldRef[*User]("email")
	AdminUserName     = NewFieldRef[*User]("name")
	AdminUserID       = NewFieldRef[*User]("id")
	AdminUserIsActive = NewFieldRef[*User]("is_active")
)

// Create admin instance
func NewUserAdmin(db *gorm.DB) *Admin[*User] {
	admin := NewAdmin[*User](db, "User")
	
	// Configure admin using admin FieldRef (for display)
	admin.SetListDisplay(AdminUserID, AdminUserEmail, AdminUserName, AdminUserIsActive).
		SetListDisplayLinks(AdminUserEmail).
		SetListEditable(AdminUserIsActive).
		SetSearchFields(AdminUserEmail, AdminUserName).
		SetListFilter(
			FilterSpec[*User]{
				Field: AdminUserIsActive,
				Type:  FilterTypeBoolean,
			},
			FilterSpec[*User]{
				Field: AdminUserEmail,
				Type:  FilterTypeChoice,
			},
		).
		SetFields(AdminUserEmail, AdminUserName, AdminUserIsActive).
		SetReadonlyFields(AdminUserID)
	
	return admin
}

// Example: Custom admin with type-safe queryset filtering
type CustomUserAdmin struct {
	*Admin[*User]
}

func NewCustomUserAdmin(db *gorm.DB) *CustomUserAdmin {
	return &CustomUserAdmin{
		Admin: NewUserAdmin(db),
	}
}

// Override GetQueryset with type-safe filtering
func (a *CustomUserAdmin) GetQueryset(ctx context.Context) *models.QuerySetImpl[*User] {
	queryset := a.manager.Filter(ctx)
	
	// Use type-safe field references for filtering
	activeCondition := models.NewCondition[*User](UserIsActive.ApplyEq(true))
	
	// Apply filter - returns *QuerySetImpl[*User]
	return queryset.Filter(activeCondition)
}

// Example: Using admin in a handler
func ExampleAdminUsage(db *gorm.DB, ctx context.Context) {
	admin := NewUserAdmin(db)
	
	// Get queryset (type-safe!)
	queryset := admin.GetQueryset(ctx)
	
	// Apply search using models field references
	searchTerm := "example"
	queryset = admin.ApplySearch(ctx, queryset, searchTerm)
	
	// Get results - returns []*User
	users, err := queryset.All(ctx)
	if err != nil {
		panic(err)
	}
	
	// Type-safe access
	for _, user := range users {
		_ = user.Email
		_ = user.Name
		_ = user.IsActive
	}
	
	_ = users
}

