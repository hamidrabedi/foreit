package admin

import (
	"testing"
	"time"

	query "github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
)

// TestUser is a test model
type TestUser struct {
	TestUserSchema
	ID        int64
	Username  string
	Email     string
	IsActive  bool
	CreatedAt time.Time
}

// TestUserSchema implements schema.Schema
type TestUserSchema struct {
	schema.BaseSchema
}

func (s TestUserSchema) Fields() []schema.Field {
	return []schema.Field{
		{Name: "ID", Type: schema.TypeInt64, PrimaryKey: true},
		{Name: "Username", Type: schema.TypeString, Required: true},
		{Name: "Email", Type: schema.TypeString},
		{Name: "IsActive", Type: schema.TypeBool},
		{Name: "CreatedAt", Type: schema.TypeTime},
	}
}

func (s TestUserSchema) Meta() schema.Meta {
	return schema.Meta{
		VerboseName:       "User",
		VerboseNamePlural: "Users",
	}
}

// TestAdmin_Register tests admin registration
func TestAdmin_Register(t *testing.T) {
	// Create a test manager
	manager, _ := query.NewManager[TestUser]("")

	// Create admin config
	config := &Config[TestUser]{
		VerboseName:       "User",
		VerboseNamePlural: "Users",
		ListPerPage:       25,
	}

	// Create schema
	s := &TestUserSchema{}

	// Register admin
	admin, err := Register[TestUser](s, manager, config)

	// Verify registration
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if admin == nil {
		t.Fatal("Admin should not be nil")
	}

	if admin.ModelName() != "TestUser" {
		t.Errorf("Expected ModelName 'TestUser', got '%s'", admin.ModelName())
	}

	if admin.Config() != config {
		t.Error("Config should match")
	}
}
