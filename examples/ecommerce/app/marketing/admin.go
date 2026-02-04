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
			CouponFieldsInstance.Code,
			CouponFieldsInstance.Name,
			CouponFieldsInstance.DiscountType,
			CouponFieldsInstance.DiscountValue,
			CouponFieldsInstance.IsActive,
			CouponFieldsInstance.ValidUntil,
		},
		ListFilter: []admin.Field{
			CouponFieldsInstance.DiscountType,
			CouponFieldsInstance.IsActive,
		},
	})

	// CouponUsage admin
	admin.Register(&admin.Config[CouponUsage]{
		Icon: "BarChart3",
		ListDisplay: []admin.Field{
			CouponUsageFieldsInstance.CouponId,
			CouponUsageFieldsInstance.OrderId,
			CouponUsageFieldsInstance.CustomerId,
			CouponUsageFieldsInstance.DiscountAmount,
			CouponUsageFieldsInstance.CreatedAt,
		},
	})

	// Review admin
	admin.Register(&admin.Config[Review]{
		Icon: "Star",
		ListDisplay: []admin.Field{
			ReviewFieldsInstance.ProductId,
			ReviewFieldsInstance.CustomerId,
			ReviewFieldsInstance.Rating,
			ReviewFieldsInstance.Status,
			ReviewFieldsInstance.IsFeatured,
			ReviewFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			ReviewFieldsInstance.Status,
			ReviewFieldsInstance.Rating,
			ReviewFieldsInstance.IsFeatured,
		},
	})

	// ReviewImage admin
	admin.Register(&admin.Config[ReviewImage]{
		Icon: "Image",
		ListDisplay: []admin.Field{
			ReviewImageFieldsInstance.ReviewId,
			ReviewImageFieldsInstance.AltText,
			ReviewImageFieldsInstance.SortOrder,
		},
	})

	// ReviewHelpfulness admin
	admin.Register(&admin.Config[ReviewHelpfulness]{
		Icon: "ThumbsUp",
		ListDisplay: []admin.Field{
			ReviewHelpfulnessFieldsInstance.ReviewId,
			ReviewHelpfulnessFieldsInstance.CustomerId,
			ReviewHelpfulnessFieldsInstance.IsHelpful,
		},
	})

	// ProductQuestion admin
	admin.Register(&admin.Config[ProductQuestion]{
		Icon: "HelpCircle",
		ListDisplay: []admin.Field{
			ProductQuestionFieldsInstance.ProductId,
			ProductQuestionFieldsInstance.CustomerId,
			ProductQuestionFieldsInstance.Status,
			ProductQuestionFieldsInstance.IsPublic,
			ProductQuestionFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			ProductQuestionFieldsInstance.Status,
			ProductQuestionFieldsInstance.IsPublic,
		},
	})
}
