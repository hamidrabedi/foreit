package marketing

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers marketing models with the admin interface
func RegisterAdmin(ctx context.Context) {
	// Coupon admin
	admin.Register(&admin.Config[Coupon]{
		Icon: "Ticket",
		ListDisplay: []admin.Field{
			CouponFields.Code,
			CouponFields.Name,
			CouponFields.DiscountType,
			CouponFields.DiscountValue,
			CouponFields.IsActive,
			CouponFields.ValidUntil,
		},
		ListFilter: []admin.Field{
			CouponFields.DiscountType,
			CouponFields.IsActive,
		},
	})

	// CouponUsage admin
	admin.Register(&admin.Config[CouponUsage]{
		Icon: "BarChart3",
		ListDisplay: []admin.Field{
			CouponUsageFields.CouponID,
			CouponUsageFields.OrderID,
			CouponUsageFields.CustomerID,
			CouponUsageFields.DiscountAmount,
			CouponUsageFields.CreatedAt,
		},
	})

	// Review admin
	admin.Register(&admin.Config[Review]{
		Icon: "Star",
		ListDisplay: []admin.Field{
			ReviewFields.ProductID,
			ReviewFields.CustomerID,
			ReviewFields.Rating,
			ReviewFields.Status,
			ReviewFields.IsFeatured,
			ReviewFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			ReviewFields.Status,
			ReviewFields.Rating,
			ReviewFields.IsFeatured,
		},
	})

	// ReviewImage admin
	admin.Register(&admin.Config[ReviewImage]{
		Icon: "Image",
		ListDisplay: []admin.Field{
			ReviewImageFields.ReviewID,
			ReviewImageFields.AltText,
			ReviewImageFields.SortOrder,
		},
	})

	// ReviewHelpfulness admin
	admin.Register(&admin.Config[ReviewHelpfulness]{
		Icon: "ThumbsUp",
		ListDisplay: []admin.Field{
			ReviewHelpfulnessFields.ReviewID,
			ReviewHelpfulnessFields.CustomerID,
			ReviewHelpfulnessFields.IsHelpful,
		},
	})

	// ProductQuestion admin
	admin.Register(&admin.Config[ProductQuestion]{
		Icon: "HelpCircle",
		ListDisplay: []admin.Field{
			ProductQuestionFields.ProductID,
			ProductQuestionFields.CustomerID,
			ProductQuestionFields.Status,
			ProductQuestionFields.IsPublic,
			ProductQuestionFields.CreatedAt,
		},
		ListFilter: []admin.Field{
			ProductQuestionFields.Status,
			ProductQuestionFields.IsPublic,
		},
	})
}
