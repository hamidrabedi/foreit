package ledger

import (
	"time"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// LedgerEntry records every balance movement.
type LedgerEntry struct {
	schema.BaseSchema

	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	AssetID      int64     `json:"asset_id" db:"asset_id"`
	EntryType    string    `json:"entry_type" db:"entry_type"`
	Amount       float64   `json:"amount" db:"amount"`
	BalanceAfter float64   `json:"balance_after" db:"balance_after"`
	RefType      string    `json:"ref_type" db:"ref_type"`
	RefID        int64     `json:"ref_id" db:"ref_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

func (LedgerEntry) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.Int64("asset_id").WithRequired(),
		schema.String("entry_type").WithRequired().WithMaxLength(30).
			WithChoices(schema.WithChoices("deposit", "Deposit", "withdrawal", "Withdrawal", "trade", "Trade", "fee", "Fee")...),
		schema.Decimal("amount").WithDefault(0),
		schema.Decimal("balance_after").WithDefault(0),
		schema.String("ref_type").WithMaxLength(50),
		schema.Int64("ref_id").WithDefault(0),
		schema.Time("created_at").WithAutoNowAdd(),
	}
}

func (LedgerEntry) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "ledger_entries",
		VerboseName:       "Ledger Entry",
		VerboseNamePlural: "Ledger Entries",
		Indexes: []schema.Index{
			{Name: "idx_ledger_user", Fields: []string{"user_id"}},
			{Name: "idx_ledger_asset", Fields: []string{"asset_id"}},
			{Name: "idx_ledger_ref", Fields: []string{"ref_type", "ref_id"}},
		},
	}
}

func (LedgerEntry) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("asset_id", "Asset").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("ledger_entries"),
	}
}

func (LedgerEntry) Hooks() *schema.ModelHooks { return nil }

func RegisterModels() {
	_ = registry.RegisterModel(&LedgerEntry{})
}
