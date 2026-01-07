package accounts

import (
	"context"
	"time"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// UserProfile stores exchange-specific user settings.
type UserProfile struct {
	schema.BaseSchema

	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Email       string    `json:"email" db:"email"`
	Tier        string    `json:"tier" db:"tier"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (UserProfile) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.String("display_name").WithMaxLength(100),
		schema.String("email").WithRequired().WithMaxLength(255),
		schema.String("tier").WithRequired().WithMaxLength(50).
			WithChoices(schema.WithChoices("basic", "Basic", "pro", "Pro", "vip", "VIP")...),
		schema.Bool("is_active").WithDefault(true),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (UserProfile) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "user_profiles",
		VerboseName:       "User Profile",
		VerboseNamePlural: "User Profiles",
		Indexes: []schema.Index{
			{Name: "idx_user_profiles_user", Fields: []string{"user_id"}, Unique: true},
		},
	}
}

func (UserProfile) Relations() []schema.Relation { return []schema.Relation{} }

func (UserProfile) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			return nil
		},
	}
}

// APIKey represents a user API key for trading.
type APIKey struct {
	schema.BaseSchema

	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	Key       string    `json:"key" db:"key"`
	Secret    string    `json:"secret" db:"secret"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (APIKey) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.String("name").WithRequired().WithMaxLength(100),
		schema.String("key").WithRequired().WithMaxLength(128).WithUnique(),
		schema.String("secret").WithRequired().WithMaxLength(128),
		schema.Bool("is_active").WithDefault(true),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (APIKey) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "api_keys",
		VerboseName:       "API Key",
		VerboseNamePlural: "API Keys",
		Indexes: []schema.Index{
			{Name: "idx_api_keys_user", Fields: []string{"user_id"}},
		},
	}
}

func (APIKey) Relations() []schema.Relation { return []schema.Relation{} }

func (APIKey) Hooks() *schema.ModelHooks { return nil }

func RegisterModels() {
	_ = registry.RegisterModel(&UserProfile{})
	_ = registry.RegisterModel(&APIKey{})
}
