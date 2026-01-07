package models

import "github.com/forgego/forge/schema"

// CustomerProfile stores additional customer details.
type CustomerProfile struct {
	schema.BaseSchema
}

func (CustomerProfile) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("customer_id").WithRequired().WithUnique(),
		schema.String("first_name").WithMaxLength(100).WithOptional(),
		schema.String("last_name").WithMaxLength(100).WithOptional(),
		schema.String("phone").WithMaxLength(20).WithOptional(),
		schema.Date("date_of_birth").WithOptional(),
		schema.DateTime("created_at").WithAutoNowAdd(),
		schema.DateTime("updated_at").WithAutoNow(),
	}
}

func (CustomerProfile) Meta() schema.Meta {
	return schema.Meta{
		TableName: "customer_profiles",
	}
}

func (CustomerProfile) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("profile"),
	}
}

func (CustomerProfile) Hooks() *schema.ModelHooks {
	return nil
}
