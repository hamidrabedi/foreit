package models

import "github.com/forgego/forge/schema"

// Address represents billing/shipping addresses.
type Address struct {
	schema.BaseSchema
}

func (Address) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("customer_id").WithRequired(),
		schema.String("type").WithRequired().WithMaxLength(20),
		schema.String("first_name").WithRequired().WithMaxLength(100),
		schema.String("last_name").WithRequired().WithMaxLength(100),
		schema.String("address_line1").WithRequired().WithMaxLength(255),
		schema.String("address_line2").WithMaxLength(255).WithOptional(),
		schema.String("city").WithRequired().WithMaxLength(100),
		schema.String("state").WithRequired().WithMaxLength(100),
		schema.String("postal_code").WithRequired().WithMaxLength(20),
		schema.String("country").WithRequired().WithMaxLength(2),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (Address) Meta() schema.Meta {
	return schema.Meta{
		TableName: "addresses",
	}
}

func (Address) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("addresses"),
	}
}

func (Address) Hooks() *schema.ModelHooks {
	return nil
}
