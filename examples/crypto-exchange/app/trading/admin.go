package trading

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers trading models with the admin interface.
func RegisterAdmin(ctx context.Context) {
	admin.Register(&admin.Config[Order]{
		Icon: "ListOrdered",
		ListDisplay: []admin.Field{
			OrderFieldsInstance.Id,
			OrderFieldsInstance.UserId,
			OrderFieldsInstance.PairId,
			OrderFieldsInstance.Side,
			OrderFieldsInstance.OrderType,
			OrderFieldsInstance.Status,
			OrderFieldsInstance.Price,
			OrderFieldsInstance.Quantity,
			OrderFieldsInstance.FilledQuantity,
			OrderFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			OrderFieldsInstance.Side,
			OrderFieldsInstance.OrderType,
			OrderFieldsInstance.Status,
		},
		SearchFields: []admin.Field{
			OrderFieldsInstance.UserId,
			OrderFieldsInstance.PairId,
		},
		Ordering: []admin.Field{
			OrderFieldsInstance.CreatedAt.Desc(),
		},
		Actions: []admin.Action[Order]{
			admin.NewAction("cancel", "Cancel", func(ctx context.Context, instances []*Order) error {
				for _, order := range instances {
					order.Status = "canceled"
				}
				return nil
			}).WithDescription("Cancel selected orders").WithDangerous(true),
		},
	})

	admin.Register(&admin.Config[Trade]{
		Icon: "ArrowLeftRight",
		ListDisplay: []admin.Field{
			TradeFieldsInstance.Id,
			TradeFieldsInstance.PairId,
			TradeFieldsInstance.BuyOrderId,
			TradeFieldsInstance.SellOrderId,
			TradeFieldsInstance.Price,
			TradeFieldsInstance.Quantity,
			TradeFieldsInstance.ExecutedAt,
		},
		ListFilter: []admin.Field{
			TradeFieldsInstance.PairId,
		},
		Ordering: []admin.Field{
			TradeFieldsInstance.ExecutedAt.Desc(),
		},
	})
}
