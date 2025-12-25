package admin

import (
	"context"
	"github.com/gogo/pkg/models"
)

// Example: Type-safe admin definition

// User model (same as in models package)
type User struct {
	ID        int
	Email     string
	Name      string
	IsActive  bool
}

func (u *User) GetID() interface{} { return u.ID }
func (u *User) SetID(id interface{}) { u.ID = id.(int) }
func (u *User) IsNew() bool { return u.ID == 0 }
func (u *User) String() string { return u.Email }

// Type-safe field references for admin
var (
	UserEmail    = NewFieldRef[*User]("email")
	UserName     = NewFieldRef[*User]("name")
	UserID       = NewFieldRef[*User]("id")
	UserIsActive = NewFieldRef[*User]("is_active")
)

// Type-safe admin definition - no strings!
var UserAdmin = NewBaseModelAdmin[*User]("User").
	SetListDisplay(UserID, UserEmail, UserName, UserIsActive).
	SetListDisplayLinks(UserEmail).
	SetListEditable(UserIsActive).
	SetSearchFields(UserEmail, UserName).
	SetListFilter(
		FilterSpec[*User]{
			Field: UserIsActive,
			Type:  FilterTypeBoolean,
		},
		FilterSpec[*User]{
			Field: UserEmail,
			Type:  FilterTypeChoice,
		},
	).
	SetFields(UserEmail, UserName, UserIsActive).
	SetReadonlyFields(UserID)

// Custom admin with overrides (still type-safe!)
type CustomUserAdmin struct {
	*BaseModelAdmin[*User]
}

func NewCustomUserAdmin() *CustomUserAdmin {
	return &CustomUserAdmin{
		BaseModelAdmin: NewBaseModelAdmin[*User]("User").
			SetListDisplay(UserID, UserEmail, UserName),
	}
}

// Override SaveModel - still type-safe!
func (a *CustomUserAdmin) SaveModel(ctx context.Context, obj *User, form interface{}, change bool) error {
	// obj is *User - type-safe!
	_ = obj.Email
	_ = obj.Name
	
	// Custom logic here
	return a.BaseModelAdmin.SaveModel(ctx, obj, form, change)
}

// Override GetQueryset - returns type-safe QuerySet!
func (a *CustomUserAdmin) GetQueryset(ctx context.Context) models.QuerySet[*User] {
	// Return type-safe queryset
	// users, err := queryset.All(ctx) returns []*User
	return nil
}

