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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		
		// Identification
		schema.String("code").Required().MaxLength(50).Unique().
			HelpText("Coupon code (e.g., SAVE20)").Build(),
		schema.String("name").Required().MaxLength(200).
			HelpText("Internal name for coupon").Build(),
		schema.Text("description").Optional().
			HelpText("Description shown to customers").Build(),
		
		// Discount type
		schema.String("discount_type").Required().MaxLength(20).
			HelpText("Type: percentage, fixed_amount, free_shipping, buy_x_get_y").Build(),
		schema.Float64("discount_value").Required().
			HelpText("Discount value (percentage or fixed amount)").Build(),
		
		// Constraints
		schema.Float64("minimum_purchase_amount").Optional().
			HelpText("Minimum cart value to apply coupon").Build(),
		schema.Float64("maximum_discount_amount").Optional().
			HelpText("Cap on discount amount for percentage coupons").Build(),
		
		// Usage limits
		schema.Int32("usage_limit").Optional().
			HelpText("Total number of times coupon can be used (null = unlimited)").Build(),
		schema.Int32("usage_limit_per_customer").Default(1).
			HelpText("Max uses per customer").Build(),
		schema.Int32("usage_count").Default(0).
			HelpText("Number of times coupon has been used").Build(),
		
		// Applicability
		schema.Bool("applies_to_all_products").Default(true).Build(),
		schema.String("product_ids").MaxLength(500).Optional().
			HelpText("Comma-separated product IDs (if not all products)").Build(),
		schema.String("category_ids").MaxLength(500).Optional().
			HelpText("Comma-separated category IDs").Build(),
		schema.String("excluded_product_ids").MaxLength(500).Optional().Build(),
		
		// Customer restrictions
		schema.Bool("applies_to_all_customers").Default(true).Build(),
		schema.String("customer_group_ids").MaxLength(500).Optional().
			HelpText("Comma-separated customer group IDs").Build(),
		schema.Text("customer_email_list").Optional().
			HelpText("Specific email addresses (one per line)").Build(),
		
		// Validity period
		schema.Time("valid_from").Required().
			HelpText("Start date/time").Build(),
		schema.Time("valid_until").Optional().
			HelpText("End date/time (null = no expiry)").Build(),
		
		// Status
		schema.Bool("is_active").Default(true).Build(),
		schema.Bool("is_public").Default(true).
			HelpText("Show in promotions list").Build(),
		
		// Priority (for multiple coupons)
		schema.Int32("priority").Default(0).
			HelpText("Higher priority coupons apply first").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (Coupon) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "coupons",
		VerboseName:      "Coupon",
		VerboseNamePlural: "Coupons",
		OrderBy:          []string{"-priority", "-created_at"},
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("coupon_id").Required().Build(),
		schema.Int64("order_id").Required().Build(),
		schema.Int64("customer_id").Required().Build(),
		
		schema.Float64("discount_amount").Required().
			HelpText("Actual discount applied").Build(),
		
		schema.Time("created_at").AutoNowAdd().Build(),
	}
}

func (CouponUsage) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "coupon_usage",
		VerboseName:      "Coupon Usage",
		VerboseNamePlural: "Coupon Usage",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("usage_records").Build(),
		schema.ForeignKey("order_id", "Order").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("coupon_usage").Build(),
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("coupon_usage").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("customer_id").Required().Build(),
		schema.Int64("order_id").Optional().
			HelpText("Order where product was purchased (verified purchase)").Build(),
		
		// Review content
		schema.String("title").Required().MaxLength(255).
			HelpText("Review title/headline").Build(),
		schema.Text("content").Required().
			HelpText("Review body").Build(),
		schema.Int32("rating").Required().
			HelpText("Rating from 1 to 5").Build(),
		
		// Verification
		schema.Bool("is_verified_purchase").Default(false).
			HelpText("Customer purchased this product").Build(),
		
		// Moderation
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, approved, rejected, flagged").Build(),
		schema.Bool("is_featured").Default(false).
			HelpText("Feature this review").Build(),
		
		// Helpfulness
		schema.Int32("helpful_count").Default(0).
			HelpText("Number of users who found this helpful").Build(),
		schema.Int32("not_helpful_count").Default(0).Build(),
		
		// Response
		schema.Text("merchant_response").Optional().
			HelpText("Response from store owner").Build(),
		schema.Time("merchant_response_at").Optional().Build(),
		schema.Int64("merchant_response_by_user_id").Optional().Build(),
		
		// Reporting
		schema.Int32("report_count").Default(0).
			HelpText("Number of times review was reported").Build(),
		schema.Text("report_reasons").Optional().
			HelpText("Reasons for reports (JSON)").Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
		schema.Time("approved_at").Optional().Build(),
	}
}

func (Review) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "reviews",
		VerboseName:      "Product Review",
		VerboseNamePlural: "Product Reviews",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("reviews").Build(),
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("reviews").Build(),
		schema.ForeignKey("order_id", "Order").
			OnDelete(schema.CascadeSET_NULL).
			Optional().
			RelatedName("reviews").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("review_id").Required().Build(),
		
		schema.String("image_url").Required().MaxLength(500).Build(),
		schema.String("thumbnail_url").MaxLength(500).Optional().Build(),
		schema.String("alt_text").MaxLength(255).Optional().Build(),
		
		schema.Int32("sort_order").Default(0).Build(),
		
		schema.Time("created_at").AutoNowAdd().Build(),
	}
}

func (ReviewImage) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "review_images",
		VerboseName:      "Review Image",
		VerboseNamePlural: "Review Images",
		OrderBy:          []string{"review_id", "sort_order"},
		Indexes: []schema.Index{
			{Name: "idx_review_image_review", Fields: []string{"review_id"}},
		},
	}
}

func (ReviewImage) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("review_id", "Review").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("images").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("review_id").Required().Build(),
		schema.Int64("customer_id").Optional().Build(),
		
		schema.Bool("is_helpful").Required().
			HelpText("True = helpful, False = not helpful").Build(),
		
		schema.String("ip_address").MaxLength(45).Optional().
			HelpText("For anonymous votes").Build(),
		
		schema.Time("created_at").AutoNowAdd().Build(),
	}
}

func (ReviewHelpfulness) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "review_helpfulness",
		VerboseName:      "Review Helpfulness",
		VerboseNamePlural: "Review Helpfulness",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("helpfulness_votes").Build(),
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Optional().
			RelatedName("review_helpfulness_votes").Build(),
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
		schema.Int64("id").Primary().AutoIncrement().Build(),
		schema.Int64("product_id").Required().Build(),
		schema.Int64("customer_id").Required().Build(),
		
		// Question
		schema.Text("question").Required().Build(),
		
		// Answer
		schema.Text("answer").Optional().Build(),
		schema.Time("answered_at").Optional().Build(),
		schema.Int64("answered_by_user_id").Optional().Build(),
		schema.String("answered_by_user_name").MaxLength(200).Optional().Build(),
		
		// Status
		schema.String("status").Required().MaxLength(20).Default("pending").
			HelpText("Status: pending, answered, hidden").Build(),
		schema.Bool("is_public").Default(true).Build(),
		
		// Helpfulness
		schema.Int32("helpful_count").Default(0).Build(),
		
		// Timestamps
		schema.Time("created_at").AutoNowAdd().Build(),
		schema.Time("updated_at").AutoNow().Build(),
	}
}

func (ProductQuestion) Meta() schema.Meta {
	return schema.Meta{
		TableName:        "product_questions",
		VerboseName:      "Product Question",
		VerboseNamePlural: "Product Questions",
		OrderBy:          []string{"-created_at"},
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
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("questions").Build(),
		schema.ForeignKey("customer_id", "Customer").
			OnDelete(schema.CascadeCASCADE).
			Required().
			RelatedName("questions").Build(),
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
