package ledger

import (
	"context"

	"github.com/forgego/forge/admin"
)

// RegisterAdmin registers ledger models with the admin interface.
func RegisterAdmin(ctx context.Context) {
	admin.Register(&admin.Config[LedgerEntry]{
		Icon: "BookOpen",
		ListDisplay: []admin.Field{
			LedgerEntryFieldsInstance.Id,
			LedgerEntryFieldsInstance.UserId,
			LedgerEntryFieldsInstance.AssetId,
			LedgerEntryFieldsInstance.EntryType,
			LedgerEntryFieldsInstance.Amount,
			LedgerEntryFieldsInstance.BalanceAfter,
			LedgerEntryFieldsInstance.CreatedAt,
		},
		ListFilter: []admin.Field{
			LedgerEntryFieldsInstance.EntryType,
			LedgerEntryFieldsInstance.AssetId,
		},
		SearchFields: []admin.Field{
			LedgerEntryFieldsInstance.RefType,
			LedgerEntryFieldsInstance.RefId,
		},
		Ordering: []admin.Field{
			LedgerEntryFieldsInstance.CreatedAt.Desc(),
		},
	})
}
