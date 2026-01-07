package wallets

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers wallet models with the admin interface.
func RegisterAdmin(ctx context.Context) {
	admin.Register(&admin.Config[Wallet]{
		Icon: "Wallet",
		ListDisplay: []admin.Field{
			WalletFieldsInstance.Id,
			WalletFieldsInstance.UserId,
			WalletFieldsInstance.AssetId,
			WalletFieldsInstance.Available,
			WalletFieldsInstance.Locked,
			WalletFieldsInstance.UpdatedAt,
		},
		ListFilter: []admin.Field{
			WalletFieldsInstance.AssetId,
		},
		SearchFields: []admin.Field{
			WalletFieldsInstance.UserId,
		},
		Ordering: []admin.Field{
			WalletFieldsInstance.UpdatedAt.Desc(),
		},
	})

	admin.Register(&admin.Config[Transfer]{
		Icon: "ArrowDownUp",
		ListDisplay: []admin.Field{
			TransferFieldsInstance.Id,
			TransferFieldsInstance.UserId,
			TransferFieldsInstance.AssetId,
			TransferFieldsInstance.Direction,
			TransferFieldsInstance.Status,
			TransferFieldsInstance.Amount,
			TransferFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			TransferFieldsInstance.Direction,
			TransferFieldsInstance.Status,
		},
		SearchFields: []admin.Field{
			TransferFieldsInstance.TxRef,
		},
		Ordering: []admin.Field{
			TransferFieldsInstance.CreatedAt.Desc(),
		},
	})
}
