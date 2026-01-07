package models

import "github.com/forgego/forge/schema"

// Product demonstrates various field types and constraints
type Product struct {
	schema.BaseSchema
}

func (Product) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(255),
		schema.Text("description"),
		schema.Decimal("price").WithRequired(),
		schema.Int64("stock_quantity").WithDefault(0),
		schema.Bool("is_active").WithDefault(true),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (Product) Meta() schema.Meta {
	return schema.Meta{
		TableName: "products",
		Indexes: []schema.Index{
			schema.Index("name"),
			schema.Index("is_active"),
		},
	}
}

func (Product) Relations() []schema.Relation {
	return []schema.Relation{}
}

// Order demonstrates foreign keys and decimal precision
type Order struct {
	schema.BaseSchema
}

func (Order) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("order_number").WithRequired().WithMaxLength(50).WithUnique(),
		schema.Int64("customer_id").WithRequired(),
		schema.Decimal("total_amount").WithRequired(),
		schema.String("status").WithRequired().WithMaxLength(50).WithDefault("pending"),
		schema.DateTime("created_at"),
		schema.DateTime("updated_at"),
	}
}

func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName: "orders",
		Indexes: []schema.Index{
			schema.Index("customer_id"),
			schema.Index("status"),
		},
	}
}

func (Order) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").WithOnDelete("RESTRICT"),
	}
}

// Customer with unique email and default values
type Customer struct {
	schema.BaseSchema
}

func (Customer) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("email").WithRequired().WithMaxLength(255).WithUnique(),
		schema.String("name").WithRequired().WithMaxLength(255),
		schema.Bool("is_active").WithDefault(true),
		schema.DateTime("created_at"),
	}
}

func (Customer) Meta() schema.Meta {
	return schema.Meta{
		TableName: "customers",
	}
}

func (Customer) Relations() []schema.Relation {
	return []schema.Relation{}
}
