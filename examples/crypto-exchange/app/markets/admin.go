package markets

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers market models with the admin interface.
func RegisterAdmin(ctx context.Context) {
	admin.Register(&admin.Config[Asset]{
		Icon: "Coins",
		ListDisplay: []admin.Field{
			AssetFieldsInstance.Id,
			AssetFieldsInstance.Symbol,
			AssetFieldsInstance.Name,
			AssetFieldsInstance.AssetType,
			AssetFieldsInstance.Precision,
			AssetFieldsInstance.IsActive,
		},
		ListFilter: []admin.Field{
			AssetFieldsInstance.AssetType,
			AssetFieldsInstance.IsActive,
		},
		SearchFields: []admin.Field{
			AssetFieldsInstance.Symbol,
			AssetFieldsInstance.Name,
		},
		Ordering: []admin.Field{
			AssetFieldsInstance.Symbol.Asc(),
		},
	})

	admin.Register(&admin.Config[TradingPair]{
		Icon: "Repeat",
		ListDisplay: []admin.Field{
			TradingPairFieldsInstance.Id,
			TradingPairFieldsInstance.Symbol,
			TradingPairFieldsInstance.Status,
			TradingPairFieldsInstance.PriceTick,
			TradingPairFieldsInstance.QtyStep,
		},
		ListFilter: []admin.Field{
			TradingPairFieldsInstance.Status,
		},
		SearchFields: []admin.Field{
			TradingPairFieldsInstance.Symbol,
		},
		Ordering: []admin.Field{
			TradingPairFieldsInstance.Symbol.Asc(),
		},
		Actions: []admin.Action[TradingPair]{
			admin.NewAction("halt", "Halt", func(ctx context.Context, instances []*TradingPair) error {
				for _, pair := range instances {
					pair.Status = "halted"
				}
				return nil
			}).WithDescription("Halt selected pairs").WithDangerous(true),
		},
	})
}
