package models

import (
	"github.com/forgego/forge/pkg/schema"
)

// Order represents a customer order
type Order struct {
	schema.BaseSchema
}

// Fields returns all field definitions for Order
func (Order) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.UUID("order_number").Required().Unique().VerboseName("Order Number").Build(),
		schema.Int64("customer_id").Required().VerboseName("Customer ID").Build(),
		schema.String("status").Required().MaxLength(50).Choices(
			schema.Choice{Value: "pending", Label: "Pending"},
			schema.Choice{Value: "confirmed", Label: "Confirmed"},
			schema.Choice{Value: "processing", Label: "Processing"},
			schema.Choice{Value: "shipped", Label: "Shipped"},
			schema.Choice{Value: "delivered", Label: "Delivered"},
			schema.Choice{Value: "cancelled", Label: "Cancelled"},
			schema.Choice{Value: "refunded", Label: "Refunded"},
		).Default("pending").VerboseName("Order Status").Build(),
		schema.Decimal("subtotal").MaxDigits(12).DecimalPlaces(2).Required().VerboseName("Subtotal").Build(),
		schema.Decimal("tax_amount").MaxDigits(12).DecimalPlaces(2).Default(0.0).VerboseName("Tax Amount").Build(),
		schema.Decimal("shipping_amount").MaxDigits(12).DecimalPlaces(2).Default(0.0).VerboseName("Shipping Amount").Build(),
		schema.Decimal("discount_amount").MaxDigits(12).DecimalPlaces(2).Default(0.0).VerboseName("Discount Amount").Build(),
		schema.Decimal("total_amount").MaxDigits(12).DecimalPlaces(2).Required().VerboseName("Total Amount").Build(),
		schema.String("currency").MaxLength(3).Default("USD").VerboseName("Currency").Build(),
		schema.Int64("billing_address_id").Required().VerboseName("Billing Address ID").Build(),
		schema.Int64("shipping_address_id").Required().VerboseName("Shipping Address ID").Build(),
		schema.String("payment_status").MaxLength(50).Choices(
			schema.Choice{Value: "pending", Label: "Pending"},
			schema.Choice{Value: "paid", Label: "Paid"},
			schema.Choice{Value: "failed", Label: "Failed"},
			schema.Choice{Value: "refunded", Label: "Refunded"},
		).Default("pending").VerboseName("Payment Status").Build(),
		schema.String("shipping_status").MaxLength(50).Choices(
			schema.Choice{Value: "pending", Label: "Pending"},
			schema.Choice{Value: "processing", Label: "Processing"},
			schema.Choice{Value: "shipped", Label: "Shipped"},
			schema.Choice{Value: "delivered", Label: "Delivered"},
			schema.Choice{Value: "returned", Label: "Returned"},
		).Default("pending").VerboseName("Shipping Status").Build(),
		schema.Text("notes").Optional().VerboseName("Order Notes").Build(),
		schema.Text("customer_notes").Optional().VerboseName("Customer Notes").Build(),
		schema.JSON("metadata").Optional().VerboseName("Additional Metadata").Build(),
		schema.DateTime("placed_at").AutoNowAdd().VerboseName("Placed At").Build(),
		schema.DateTime("confirmed_at").Optional().VerboseName("Confirmed At").Build(),
		schema.DateTime("shipped_at").Optional().VerboseName("Shipped At").Build(),
		schema.DateTime("delivered_at").Optional().VerboseName("Delivered At").Build(),
		schema.DateTime("cancelled_at").Optional().VerboseName("Cancelled At").Build(),
		schema.DateTime("created_at").AutoNowAdd().Build(),
		schema.DateTime("updated_at").AutoNow().Build(),
	}
}

// Meta returns model metadata
func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "orders",
		VerboseName:       "Order",
		VerboseNamePlural: "Orders",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_order_order_number", Fields: []string{"order_number"}, Unique: true},
			{Name: "idx_order_customer_id", Fields: []string{"customer_id"}, Unique: false},
			{Name: "idx_order_status", Fields: []string{"status"}, Unique: false},
			{Name: "idx_order_payment_status", Fields: []string{"payment_status"}, Unique: false},
			{Name: "idx_order_shipping_status", Fields: []string{"shipping_status"}, Unique: false},
			{Name: "idx_order_placed_at", Fields: []string{"placed_at"}, Unique: false},
		},
	}
}

// Relations returns all relationship definitions
func (Order) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").Required().OnDelete(schema.CascadePROTECT).RelatedName("orders").Build(),
		schema.ForeignKey("billing_address_id", "Address").Required().OnDelete(schema.CascadePROTECT).RelatedName("billing_orders").Build(),
		schema.ForeignKey("shipping_address_id", "Address").Required().OnDelete(schema.CascadePROTECT).RelatedName("shipping_orders").Build(),
		schema.OneToMany("items", "OrderItem", "order_id").CascadeOnDelete().Build(),
		schema.OneToMany("payments", "Payment", "order_id").CascadeOnDelete().Build(),
		schema.OneToOne("shipping", "Shipping").Optional().OnDelete(schema.CascadeCASCADE).RelatedName("order").Build(),
	}
}

// Hooks returns model lifecycle hooks
func (Order) Hooks() *schema.ModelHooks {
	return nil
}

