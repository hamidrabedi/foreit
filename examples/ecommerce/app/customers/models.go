package customers

import (
	"context"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// CustomerGroup represents customer segmentation for pricing/promotions
type CustomerGroup struct {
	schema.BaseSchema
}

func (CustomerGroup) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.String("name").Required().MaxLength(200).Unique().
			HelpText("Group name (e.g., 'VIP', 'Wholesale')").Build(),
		schema.String("code").Required().MaxLength(50).Unique().Build(),
		schema.Text("description").Optional().Build(),
		
		// Pricing
		schema.Float64("discount_percentage").Default(0.0).
			HelpText("Default discount percentage for this group").Build(),
		
		// Status
		schema.Bool("is_active").Default(true).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (CustomerGroup) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "customer_groups",
		VerboseName:      "Customer Group",
		VerboseNamePlural: "Customer Groups",
		OrderBy:          []string{"name"},
	}
}

func (CustomerGroup) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (CustomerGroup) Hooks() *schema.ModelHooks {
	return nil
}

// Customer represents a customer account
type Customer struct {
	schema.BaseSchema
}

func (Customer) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		
		// Authentication (if using built-in auth)
		schema.String("email").Required().MaxLength(255).Unique().
			HelpText("Customer email address").Build(),
		schema.String("password_hash").Required().MaxLength(255).
			HelpText("Hashed password").Build(),
		
		// Personal information
		schema.String("first_name").Required().MaxLength(100).Build(),
		schema.String("last_name").Required().MaxLength(100).Build(),
		schema.String("phone").MaxLength(20).Optional().Build(),
		schema.Date("date_of_birth").Optional().Build(),
		schema.String("gender").MaxLength(20).Optional().
			HelpText("Gender: male, female, other, prefer_not_to_say").Build(),
		
		// Business (for B2B)
		schema.String("company_name").MaxLength(200).Optional().Build(),
		schema.String("tax_id").MaxLength(50).Optional().
			HelpText("Tax ID / VAT number").Build(),
		
		// Customer group
		schema.Int64("customer_group_id").Optional().Build(),
		
		// Status
		schema.Bool("is_active").Default(true).
			HelpText("Is account active").Build(),
		schema.Bool("is_verified").Default(false).
			HelpText("Is email verified").Build(),
		schema.Bool("accepts_marketing").Default(false).
			HelpText("Accepts marketing emails").Build(),
		
		// Account security
		schema.String("verification_token").MaxLength(255).Optional().Build(),
		schema.String("reset_password_token").MaxLength(255).Optional().Build(),
		schema.Time("reset_password_expires_at").Optional().Build(),
		schema.Time("last_login_at").Optional().Build(),
		schema.String("last_login_ip").MaxLength(45).Optional().Build(),
		
		// Statistics
		schema.Int32("total_orders").Default(0).
			HelpText("Total number of orders").Build(),
		schema.Float64("total_spent").Default(0.0).
			HelpText("Total amount spent").Build(),
		schema.Float64("average_order_value").Default(0.0).Build(),
		
		// Preferences
		schema.String("preferred_language").MaxLength(10).Default("en").Build(),
		schema.String("preferred_currency").MaxLength(3).Default("USD").Build(),
		
		// Notes
		schema.Text("notes").Optional().
			HelpText("Admin notes about customer").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Customer) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "customers",
		VerboseName:      "Customer",
		VerboseNamePlural: "Customers",
		OrderBy:          []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_customer_email", Fields: []string{"email"}, Unique: true},
			{Name: "idx_customer_group", Fields: []string{"customer_group_id"}},
			{Name: "idx_customer_active", Fields: []string{"is_active"}},
			{Name: "idx_customer_verified", Fields: []string{"is_verified"}},
		},
	}
}

func (Customer) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_group_id", "CustomerGroup").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("customers").Build(),
	}
}

func (Customer) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Hash password before storing
			// Generate verification token
			return nil
		},
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Validate email format
			// Validate phone format
			return nil
		},
	}
}

// Address represents shipping/billing addresses
type Address struct {
	schema.BaseSchema
}

func (Address) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("customer_id").Required().Build(),
		
		// Address type
		schema.String("address_type").Required().MaxLength(20).
			HelpText("Type: shipping, billing, both").Build(),
		
		// Contact
		schema.String("first_name").Required().MaxLength(100).Build(),
		schema.String("last_name").Required().MaxLength(100).Build(),
		schema.String("company_name").MaxLength(200).Optional().Build(),
		schema.String("phone").MaxLength(20).Optional().Build(),
		
		// Address
		schema.String("address_line1").Required().MaxLength(255).
			HelpText("Street address, P.O. box").Build(),
		schema.String("address_line2").MaxLength(255).Optional().
			HelpText("Apartment, suite, unit, building, floor, etc.").Build(),
		schema.String("city").Required().MaxLength(100).Build(),
		schema.String("state_province").MaxLength(100).Optional().
			HelpText("State, province, region").Build(),
		schema.String("postal_code").Required().MaxLength(20).Build(),
		schema.String("country_code").Required().MaxLength(2).
			HelpText("ISO 3166-1 alpha-2 country code").Build(),
		schema.String("country_name").Required().MaxLength(100).Build(),
		
		// Geolocation (optional)
		schema.Float64("latitude").Optional().Build(),
		schema.Float64("longitude").Optional().Build(),
		
		// Preferences
		schema.Bool("is_default_shipping").Default(false).Build(),
		schema.Bool("is_default_billing").Default(false).Build(),
		
		// Special instructions
		schema.Text("delivery_instructions").Optional().
			HelpText("Delivery notes, gate codes, etc.").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Address) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "addresses",
		VerboseName:      "Address",
		VerboseNamePlural: "Addresses",
		OrderBy:          []string{"-is_default_shipping", "-is_default_billing", "-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_address_customer", Fields: []string{"customer_id"}},
			{Name: "idx_address_type", Fields: []string{"address_type"}},
			{Name: "idx_address_country", Fields: []string{"country_code"}},
		},
	}
}

func (Address) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("addresses").Build(),
	}
}

func (Address) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Ensure only one default shipping address per customer
			// Ensure only one default billing address per customer
			// Validate country code
			// Validate postal code format
			return nil
		},
	}
}

// WishList represents customer wish lists
type WishList struct {
	schema.BaseSchema
}

func (WishList) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("customer_id").Required().Build(),
		
		schema.String("name").Required().MaxLength(200).
			HelpText("Wish list name (e.g., 'Birthday Gifts')").Build(),
		schema.Text("description").Optional().Build(),
		
		// Visibility
		schema.Bool("is_public").Default(false).
			HelpText("Share with others").Build(),
		schema.String("share_token").MaxLength(100).Optional().
			HelpText("Token for sharing wish list").Build(),
		
		// Status
		schema.Bool("is_default").Default(false).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (WishList) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "wish_lists",
		VerboseName:      "Wish List",
		VerboseNamePlural: "Wish Lists",
		OrderBy:          []string{"-is_default", "-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_wishlist_customer", Fields: []string{"customer_id"}},
			{Name: "idx_wishlist_token", Fields: []string{"share_token"}},
		},
	}
}

func (WishList) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("wish_lists").Build(),
	}
}

func (WishList) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeCreate: func(ctx context.Context, instance interface{}) error {
			// Generate share token if public
			return nil
		},
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Ensure only one default wish list per customer
			return nil
		},
	}
}

// WishListItem represents items in a wish list
type WishListItem struct {
	schema.BaseSchema
}

func (WishListItem) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("wish_list_id").Required().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("variant_id").Optional().
			HelpText("Specific product variant").Build(),
		
		schema.Int32("desired_quantity").Default(1).Build(),
		schema.Float64("price_when_added").Optional().
			HelpText("Price when added to wish list").Build(),
		
		schema.Text("notes").Optional().
			HelpText("Personal notes about the item").Build(),
		schema.Int32("priority").Default(0).
			HelpText("Priority: 0=normal, 1=high, etc.").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (WishListItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "wish_list_items",
		VerboseName:      "Wish List Item",
		VerboseNamePlural: "Wish List Items",
		OrderBy:          []string{"-priority", "-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_wishlist_item_list", Fields: []string{"wish_list_id"}},
			{Name: "idx_wishlist_item_product", Fields: []string{"product_id"}},
			{Name: "idx_wishlist_item_variant", Fields: []string{"variant_id"}},
		},
		UniqueTogether: [][]string{
			{"wish_list_id", "product_id", "variant_id"},
		},
	}
}

func (WishListItem) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("wish_list_id", "WishList").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("items").Build(),
		schema.ForeignKey("product_id", "Product").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("wish_list_items").Build(),
		schema.ForeignKey("variant_id", "ProductVariant").
			OnDelete(schema.CascadeCASCADE).
			Optional().
			RelatedName("wish_list_items").Build(),
	}
}

func (WishListItem) Hooks() *schema.ModelHooks {
	return nil
}

// RegisterModels registers customer models with the framework
func RegisterModels() {
	registry.RegisterModel(&CustomerGroup{})
	registry.RegisterModel(&Customer{})
	registry.RegisterModel(&Address{})
	registry.RegisterModel(&WishList{})
	registry.RegisterModel(&WishListItem{})
}
