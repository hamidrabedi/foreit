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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("name").WithRequired().WithMaxLength(200).WithUnique().
			WithHelpText("Group name (e.g., 'VIP', 'Wholesale')"),
		schema.String("code").WithRequired().WithMaxLength(50).WithUnique(),
		schema.Text("description").WithOptional(),

		// Pricing
		schema.Float64("discount_percentage").WithDefault(0.0).
			WithHelpText("Default discount percentage for this group"),

		// Status
		schema.Bool("is_active").WithDefault(true),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (CustomerGroup) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "customer_groups",
		VerboseName:       "Customer Group",
		VerboseNamePlural: "Customer Groups",
		OrderBy:           []string{"name"},
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),

		// Authentication (if using built-in auth)
		schema.String("email").WithRequired().WithMaxLength(255).WithUnique().
			WithHelpText("Customer email address"),
		schema.String("password_hash").WithRequired().WithMaxLength(255).
			WithHelpText("Hashed password"),

		// Personal information
		schema.String("first_name").WithRequired().WithMaxLength(100),
		schema.String("last_name").WithRequired().WithMaxLength(100),
		schema.String("phone").WithMaxLength(20).WithOptional(),
		schema.Date("date_of_birth").WithOptional(),
		schema.String("gender").WithMaxLength(20).WithOptional().
			WithHelpText("Gender: male, female, other, prefer_not_to_say"),

		// Business (for B2B)
		schema.String("company_name").WithMaxLength(200).WithOptional(),
		schema.String("tax_id").WithMaxLength(50).WithOptional().
			WithHelpText("Tax ID / VAT number"),

		// Customer group
		schema.Int64("customer_group_id").WithOptional(),

		// Status
		schema.Bool("is_active").WithDefault(true).
			WithHelpText("Is account active"),
		schema.Bool("is_verified").WithDefault(false).
			WithHelpText("Is email verified"),
		schema.Bool("accepts_marketing").WithDefault(false).
			WithHelpText("Accepts marketing emails"),

		// Account security
		schema.String("verification_token").WithMaxLength(255).WithOptional(),
		schema.String("reset_password_token").WithMaxLength(255).WithOptional(),
		schema.Time("reset_password_expires_at").WithOptional(),
		schema.Time("last_login_at").WithOptional(),
		schema.String("last_login_ip").WithMaxLength(45).WithOptional(),

		// Statistics
		schema.Int32("total_orders").WithDefault(0).
			WithHelpText("Total number of orders"),
		schema.Float64("total_spent").WithDefault(0.0).
			WithHelpText("Total amount spent"),
		schema.Float64("average_order_value").WithDefault(0.0),

		// Preferences
		schema.String("preferred_language").WithMaxLength(10).WithDefault("en"),
		schema.String("preferred_currency").WithMaxLength(3).WithDefault("USD"),

		// Notes
		schema.Text("notes").WithOptional().
			WithHelpText("Admin notes about customer"),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Customer) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "customers",
		VerboseName:       "Customer",
		VerboseNamePlural: "Customers",
		OrderBy:           []string{"-created_at"},
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
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("customers"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("customer_id").WithRequired(),

		// Address type
		schema.String("address_type").WithRequired().WithMaxLength(20).
			WithHelpText("Type: shipping, billing, both"),

		// Contact
		schema.String("first_name").WithRequired().WithMaxLength(100),
		schema.String("last_name").WithRequired().WithMaxLength(100),
		schema.String("company_name").WithMaxLength(200).WithOptional(),
		schema.String("phone").WithMaxLength(20).WithOptional(),

		// Address
		schema.String("address_line1").WithRequired().WithMaxLength(255).
			WithHelpText("Street address, P.O. box"),
		schema.String("address_line2").WithMaxLength(255).WithOptional().
			WithHelpText("Apartment, suite, unit, building, floor, etc."),
		schema.String("city").WithRequired().WithMaxLength(100),
		schema.String("state_province").WithMaxLength(100).WithOptional().
			WithHelpText("State, province, region"),
		schema.String("postal_code").WithRequired().WithMaxLength(20),
		schema.String("country_code").WithRequired().WithMaxLength(2).
			WithHelpText("ISO 3166-1 alpha-2 country code"),
		schema.String("country_name").WithRequired().WithMaxLength(100),

		// Geolocation (optional)
		schema.Float64("latitude").WithOptional(),
		schema.Float64("longitude").WithOptional(),

		// Preferences
		schema.Bool("is_default_shipping").WithDefault(false),
		schema.Bool("is_default_billing").WithDefault(false),

		// Special instructions
		schema.Text("delivery_instructions").WithOptional().
			WithHelpText("Delivery notes, gate codes, etc."),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Address) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "addresses",
		VerboseName:       "Address",
		VerboseNamePlural: "Addresses",
		OrderBy:           []string{"-is_default_shipping", "-is_default_billing", "-created_at"},
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
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("addresses"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("customer_id").WithRequired(),

		schema.String("name").WithRequired().WithMaxLength(200).
			WithHelpText("Wish list name (e.g., 'Birthday Gifts')"),
		schema.Text("description").WithOptional(),

		// Visibility
		schema.Bool("is_public").WithDefault(false).
			WithHelpText("Share with others"),
		schema.String("share_token").WithMaxLength(100).WithOptional().
			WithHelpText("Token for sharing wish list"),

		// Status
		schema.Bool("is_default").WithDefault(false),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (WishList) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "wish_lists",
		VerboseName:       "Wish List",
		VerboseNamePlural: "Wish Lists",
		OrderBy:           []string{"-is_default", "-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_wishlist_customer", Fields: []string{"customer_id"}},
			{Name: "idx_wishlist_token", Fields: []string{"share_token"}},
		},
	}
}

func (WishList) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("wish_lists"),
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
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("wish_list_id").WithRequired(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("variant_id").WithOptional().
			WithHelpText("Specific product variant"),

		schema.Int32("desired_quantity").WithDefault(1),
		schema.Float64("price_when_added").WithOptional().
			WithHelpText("Price when added to wish list"),

		schema.Text("notes").WithOptional().
			WithHelpText("Personal notes about the item"),
		schema.Int32("priority").WithDefault(0).
			WithHelpText("Priority: 0=normal, 1=high, etc."),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (WishListItem) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "wish_list_items",
		VerboseName:       "Wish List Item",
		VerboseNamePlural: "Wish List Items",
		OrderBy:           []string{"-priority", "-created_at"},
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
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("items"),
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("wish_list_items"),
		schema.ForeignKey("variant_id", "ProductVariant").
			WithOnDelete(schema.CascadeCASCADE).
			WithOptional().
			WithRelatedName("wish_list_items"),
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
