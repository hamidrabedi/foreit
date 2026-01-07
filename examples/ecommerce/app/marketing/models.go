package marketing

import (
	"context"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Coupon represents discount codes and promotions
type Coupon struct {
	schema.BaseSchema
}

func (Coupon) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),

		// Identification
		schema.String("code").WithRequired().WithMaxLength(50).WithUnique().
			WithHelpText("Coupon code (e.g., SAVE20)"),
		schema.String("name").WithRequired().WithMaxLength(200).
			WithHelpText("Internal name for coupon"),
		schema.Text("description").WithOptional().
			WithHelpText("Description shown to customers"),

		// Discount type
		schema.String("discount_type").WithRequired().WithMaxLength(20).
			WithHelpText("Type: percentage, fixed_amount, free_shipping, buy_x_get_y"),
		schema.Float64("discount_value").WithRequired().
			WithHelpText("Discount value (percentage or fixed amount)"),

		// Constraints
		schema.Float64("minimum_purchase_amount").WithOptional().
			WithHelpText("Minimum cart value to apply coupon"),
		schema.Float64("maximum_discount_amount").WithOptional().
			WithHelpText("Cap on discount amount for percentage coupons"),

		// Usage limits
		schema.Int32("usage_limit").WithOptional().
			WithHelpText("Total number of times coupon can be used (null = unlimited)"),
		schema.Int32("usage_limit_per_customer").WithDefault(1).
			WithHelpText("Max uses per customer"),
		schema.Int32("usage_count").WithDefault(0).
			WithHelpText("Number of times coupon has been used"),

		// Applicability
		schema.Bool("applies_to_all_products").WithDefault(true),
		schema.String("product_ids").WithMaxLength(500).WithOptional().
			WithHelpText("Comma-separated product IDs (if not all products)"),
		schema.String("category_ids").WithMaxLength(500).WithOptional().
			WithHelpText("Comma-separated category IDs"),
		schema.String("excluded_product_ids").WithMaxLength(500).WithOptional(),

		// Customer restrictions
		schema.Bool("applies_to_all_customers").WithDefault(true),
		schema.String("customer_group_ids").WithMaxLength(500).WithOptional().
			WithHelpText("Comma-separated customer group IDs"),
		schema.Text("customer_email_list").WithOptional().
			WithHelpText("Specific email addresses (one per line)"),

		// Validity period
		schema.Time("valid_from").WithRequired().
			WithHelpText("Start date/time"),
		schema.Time("valid_until").WithOptional().
			WithHelpText("End date/time (null = no expiry)"),

		// Status
		schema.Bool("is_active").WithDefault(true),
		schema.Bool("is_public").WithDefault(true).
			WithHelpText("Show in promotions list"),

		// Priority (for multiple coupons)
		schema.Int32("priority").WithDefault(0).
			WithHelpText("Higher priority coupons apply first"),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Coupon) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "coupons",
		VerboseName:       "Coupon",
		VerboseNamePlural: "Coupons",
		OrderBy:           []string{"-priority", "-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_coupon_code", Fields: []string{"code"}, Unique: true},
			{Name: "idx_coupon_active", Fields: []string{"is_active"}},
			{Name: "idx_coupon_valid_from", Fields: []string{"valid_from"}},
			{Name: "idx_coupon_valid_until", Fields: []string{"valid_until"}},
		},
	}
}

func (Coupon) Relations() []schema.Relation {
	return []schema.Relation{}
}

func (Coupon) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Validate discount_type
			// Uppercase coupon code
			// Validate dates
			return nil
		},
	}
}

// CouponUsage tracks coupon usage
type CouponUsage struct {
	schema.BaseSchema
}

func (CouponUsage) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("coupon_id").WithRequired(),
		schema.Int64("order_id").WithRequired(),
		schema.Int64("customer_id").WithRequired(),

		schema.Float64("discount_amount").WithRequired().
			WithHelpText("Actual discount applied"),

		schema.Time("created_at").WithAutoNowAdd(),
	}
}

func (CouponUsage) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "coupon_usage",
		VerboseName:       "Coupon Usage",
		VerboseNamePlural: "Coupon Usage",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_usage_coupon", Fields: []string{"coupon_id"}},
			{Name: "idx_usage_order", Fields: []string{"order_id"}},
			{Name: "idx_usage_customer", Fields: []string{"customer_id"}},
		},
	}
}

func (CouponUsage) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("coupon_id", "Coupon").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("usage_records"),
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("coupon_usage"),
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("coupon_usage"),
	}
}

func (CouponUsage) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Increment coupon usage_count
			return nil
		},
	}
}

// Review represents product reviews
type Review struct {
	schema.BaseSchema
}

func (Review) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("customer_id").WithRequired(),
		schema.Int64("order_id").WithOptional().
			WithHelpText("Order where product was purchased (verified purchase)"),

		// Review content
		schema.String("title").WithRequired().WithMaxLength(255).
			WithHelpText("Review title/headline"),
		schema.Text("content").WithRequired().
			WithHelpText("Review body"),
		schema.Int32("rating").WithRequired().
			WithHelpText("Rating from 1 to 5"),

		// Verification
		schema.Bool("is_verified_purchase").WithDefault(false).
			WithHelpText("Customer purchased this product"),

		// Moderation
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, approved, rejected, flagged"),
		schema.Bool("is_featured").WithDefault(false).
			WithHelpText("Feature this review"),

		// Helpfulness
		schema.Int32("helpful_count").WithDefault(0).
			WithHelpText("Number of users who found this helpful"),
		schema.Int32("not_helpful_count").WithDefault(0),

		// Response
		schema.Text("merchant_response").WithOptional().
			WithHelpText("Response from store owner"),
		schema.Time("merchant_response_at").WithOptional(),
		schema.Int64("merchant_response_by_user_id").WithOptional(),

		// Reporting
		schema.Int32("report_count").WithDefault(0).
			WithHelpText("Number of times review was reported"),
		schema.Text("report_reasons").WithOptional().
			WithHelpText("Reasons for reports (JSON)"),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
		schema.Time("approved_at").WithOptional(),
	}
}

func (Review) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "reviews",
		VerboseName:       "Product Review",
		VerboseNamePlural: "Product Reviews",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_review_product", Fields: []string{"product_id"}},
			{Name: "idx_review_customer", Fields: []string{"customer_id"}},
			{Name: "idx_review_order", Fields: []string{"order_id"}},
			{Name: "idx_review_status", Fields: []string{"status"}},
			{Name: "idx_review_rating", Fields: []string{"rating"}},
			{Name: "idx_review_featured", Fields: []string{"is_featured"}},
		},
		UniqueTogether: [][]string{
			{"product_id", "customer_id", "order_id"},
		},
	}
}

func (Review) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("reviews"),
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("reviews"),
		schema.ForeignKey("order_id", "Order").
			WithOnDelete(schema.CascadeSET_NULL).
			WithOptional().
			WithRelatedName("reviews"),
	}
}

func (Review) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		BeforeSave: func(ctx context.Context, instance interface{}) error {
			// Validate rating (1-5)
			// Check for duplicate review
			return nil
		},
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Update product rating statistics
			// Send notification to merchant
			return nil
		},
		AfterUpdate: func(ctx context.Context, instance interface{}) error {
			// Update product rating statistics
			// Notify customer of status change
			return nil
		},
	}
}

// ReviewImage represents images attached to reviews
type ReviewImage struct {
	schema.BaseSchema
}

func (ReviewImage) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("review_id").WithRequired(),

		schema.String("image_url").WithRequired().WithMaxLength(500),
		schema.String("thumbnail_url").WithMaxLength(500).WithOptional(),
		schema.String("alt_text").WithMaxLength(255).WithOptional(),

		schema.Int32("sort_order").WithDefault(0),

		schema.Time("created_at").WithAutoNowAdd(),
	}
}

func (ReviewImage) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "review_images",
		VerboseName:       "Review Image",
		VerboseNamePlural: "Review Images",
		OrderBy:           []string{"review_id", "sort_order"},
		Indexes: []schema.Index{
			{Name: "idx_review_image_review", Fields: []string{"review_id"}},
		},
	}
}

func (ReviewImage) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("review_id", "Review").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("images"),
	}
}

func (ReviewImage) Hooks() *schema.ModelHooks {
	return nil
}

// ReviewHelpfulness tracks if users found reviews helpful
type ReviewHelpfulness struct {
	schema.BaseSchema
}

func (ReviewHelpfulness) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("review_id").WithRequired(),
		schema.Int64("customer_id").WithOptional(),

		schema.Bool("is_helpful").WithRequired().
			WithHelpText("True = helpful, False = not helpful"),

		schema.String("ip_address").WithMaxLength(45).WithOptional().
			WithHelpText("For anonymous votes"),

		schema.Time("created_at").WithAutoNowAdd(),
	}
}

func (ReviewHelpfulness) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "review_helpfulness",
		VerboseName:       "Review Helpfulness",
		VerboseNamePlural: "Review Helpfulness",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_helpfulness_review", Fields: []string{"review_id"}},
			{Name: "idx_helpfulness_customer", Fields: []string{"customer_id"}},
		},
		UniqueTogether: [][]string{
			{"review_id", "customer_id"},
			{"review_id", "ip_address"},
		},
	}
}

func (ReviewHelpfulness) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("review_id", "Review").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("helpfulness_votes"),
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithOptional().
			WithRelatedName("review_helpfulness_votes"),
	}
}

func (ReviewHelpfulness) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterCreate: func(ctx context.Context, instance interface{}) error {
			// Update review helpful_count or not_helpful_count
			return nil
		},
	}
}

// ProductQuestion represents Q&A for products
type ProductQuestion struct {
	schema.BaseSchema
}

func (ProductQuestion) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("product_id").WithRequired(),
		schema.Int64("customer_id").WithRequired(),

		// Question
		schema.Text("question").WithRequired(),

		// Answer
		schema.Text("answer").WithOptional(),
		schema.Time("answered_at").WithOptional(),
		schema.Int64("answered_by_user_id").WithOptional(),
		schema.String("answered_by_user_name").WithMaxLength(200).WithOptional(),

		// Status
		schema.String("status").WithRequired().WithMaxLength(20).WithDefault("pending").
			WithHelpText("Status: pending, answered, hidden"),
		schema.Bool("is_public").WithDefault(true),

		// Helpfulness
		schema.Int32("helpful_count").WithDefault(0),

		// Timestamps
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (ProductQuestion) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "product_questions",
		VerboseName:       "Product Question",
		VerboseNamePlural: "Product Questions",
		OrderBy:           []string{"-created_at"},
		Indexes: []schema.Index{
			{Name: "idx_question_product", Fields: []string{"product_id"}},
			{Name: "idx_question_customer", Fields: []string{"customer_id"}},
			{Name: "idx_question_status", Fields: []string{"status"}},
		},
	}
}

func (ProductQuestion) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("product_id", "Product").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("questions"),
		schema.ForeignKey("customer_id", "Customer").
			WithOnDelete(schema.CascadeCASCADE).
			WithRequired().
			WithRelatedName("questions"),
	}
}

func (ProductQuestion) Hooks() *schema.ModelHooks {
	return &schema.ModelHooks{
		AfterUpdate: func(ctx context.Context, instance interface{}) error {
			// Send notification when answered
			return nil
		},
	}
}

// RegisterModels registers marketing models with the framework
func RegisterModels() {
	registry.RegisterModel(&Coupon{})
	registry.RegisterModel(&CouponUsage{})
	registry.RegisterModel(&Review{})
	registry.RegisterModel(&ReviewImage{})
	registry.RegisterModel(&ReviewHelpfulness{})
	registry.RegisterModel(&ProductQuestion{})
}
