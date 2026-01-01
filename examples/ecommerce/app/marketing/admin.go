package marketing

import (
	"context"
	
	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/db"
)

// RegisterAdmin registers marketing models with the admin interface
func RegisterAdmin(ctx context.Context, registry *admin.Registry, database *db.DB) {
	// Coupon admin
	couponConfig := &admin.ModelConfig{
		Name:             "Coupon",
		PluralName:       "Coupons",
		Icon:             "🎟️",
		ListDisplay:      []string{"id", "code", "name", "discount_type", "discount_value", "usage_count", "usage_limit", "valid_from", "valid_until", "is_active"},
		ListFilter:       []string{"discount_type", "is_active", "is_public"},
		SearchFields:     []string{"code", "name", "description"},
		OrderBy:          []string{"-priority", "-created_at"},
		PerPage:          20,
		Actions:          []string{"delete", "activate", "deactivate", "duplicate", "export"},
		ExportFormats:    []string{"csv", "json"},
		BulkActions:      true,
	}
	registry.Register("Coupon", &Coupon{}, couponConfig)
	
	// CouponUsage admin
	usageConfig := &admin.ModelConfig{
		Name:             "Coupon Usage",
		PluralName:       "Coupon Usage",
		Icon:             "📊",
		ListDisplay:      []string{"id", "coupon_id", "order_id", "customer_id", "discount_amount", "created_at"},
		ListFilter:       []string{"coupon_id", "customer_id"},
		SearchFields:     []string{"order_id", "customer_id"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"export"},
		ExportFormats:    []string{"csv", "json"},
	}
	registry.Register("CouponUsage", &CouponUsage{}, usageConfig)
	
	// Review admin
	reviewConfig := &admin.ModelConfig{
		Name:             "Product Review",
		PluralName:       "Product Reviews",
		Icon:             "⭐",
		ListDisplay:      []string{"id", "product_id", "customer_id", "title", "rating", "is_verified_purchase", "status", "helpful_count", "is_featured", "created_at"},
		ListFilter:       []string{"status", "rating", "is_verified_purchase", "is_featured"},
		SearchFields:     []string{"title", "content", "customer_id"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"approve", "reject", "feature", "unfeature", "delete", "export"},
		ExportFormats:    []string{"csv", "json"},
		BulkActions:      true,
	}
	registry.Register("Review", &Review{}, reviewConfig)
	
	// ReviewImage admin
	reviewImageConfig := &admin.ModelConfig{
		Name:             "Review Image",
		PluralName:       "Review Images",
		Icon:             "🖼️",
		ListDisplay:      []string{"id", "review_id", "alt_text", "sort_order", "created_at"},
		ListFilter:       []string{"review_id"},
		SearchFields:     []string{"alt_text"},
		OrderBy:          []string{"review_id", "sort_order"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("ReviewImage", &ReviewImage{}, reviewImageConfig)
	
	// ReviewHelpfulness admin
	helpfulnessConfig := &admin.ModelConfig{
		Name:             "Review Helpfulness",
		PluralName:       "Review Helpfulness",
		Icon:             "👍",
		ListDisplay:      []string{"id", "review_id", "customer_id", "is_helpful", "created_at"},
		ListFilter:       []string{"is_helpful", "review_id"},
		SearchFields:     []string{"customer_id"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"delete"},
	}
	registry.Register("ReviewHelpfulness", &ReviewHelpfulness{}, helpfulnessConfig)
	
	// ProductQuestion admin
	questionConfig := &admin.ModelConfig{
		Name:             "Product Question",
		PluralName:       "Product Questions",
		Icon:             "❓",
		ListDisplay:      []string{"id", "product_id", "customer_id", "question", "status", "is_public", "answered_at", "created_at"},
		ListFilter:       []string{"status", "is_public", "product_id"},
		SearchFields:     []string{"question", "answer"},
		OrderBy:          []string{"-created_at"},
		PerPage:          20,
		Actions:          []string{"answer", "hide", "delete"},
	}
	registry.Register("ProductQuestion", &ProductQuestion{}, questionConfig)
}
