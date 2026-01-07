package models

import "github.com/forgego/forge/schema"

// Customer represents a customer account.
type Customer struct {
	schema.BaseSchema
}

func (Customer) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.UUID("uuid").WithRequired(),
		schema.String("email").WithRequired().WithMaxLength(255),
		schema.String("password_hash").WithRequired().WithMaxLength(255),
		schema.String("first_name").WithRequired().WithMaxLength(100),
		schema.String("last_name").WithRequired().WithMaxLength(100),
		schema.Bool("is_active").WithDefault(true),
		schema.JSON("metadata").WithOptional(),
		schema.Decimal("lifetime_value").WithDefault(0.0).WithMaxDigits(12).WithDecimalPlaces(2),
		schema.DateTime("last_login").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Customer) Meta() schema.Meta {
	return schema.Meta{
		TableName: "customers",
		Indexes: []schema.Index{
			{Name: "idx_customer_email", Fields: []string{"email"}, Unique: true},
			{Name: "idx_customer_uuid", Fields: []string{"uuid"}, Unique: true},
		},
	}
}

func (Customer) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Customer) Hooks() *schema.ModelHooks {
	return nil
}
