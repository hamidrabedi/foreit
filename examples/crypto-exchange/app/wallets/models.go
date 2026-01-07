package wallets

import (
	"time"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Wallet tracks balances per user and asset.
type Wallet struct {
	schema.BaseSchema

	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	AssetID   int64     `json:"asset_id" db:"asset_id"`
	Available float64   `json:"available" db:"available"`
	Locked    float64   `json:"locked" db:"locked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Wallet) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.Int64("asset_id").WithRequired(),
		schema.Decimal("available").WithDefault(0),
		schema.Decimal("locked").WithDefault(0),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Wallet) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "wallets",
		VerboseName:       "Wallet",
		VerboseNamePlural: "Wallets",
		Indexes: []schema.Index{
			{Name: "idx_wallets_user", Fields: []string{"user_id"}},
			{Name: "idx_wallets_asset", Fields: []string{"asset_id"}},
		},
		UniqueTogether: [][]string{{"user_id", "asset_id"}},
	}
}

func (Wallet) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("asset_id", "Asset").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("wallets"),
	}
}

func (Wallet) Hooks() *schema.ModelHooks { return nil }

// Transfer represents a deposit/withdrawal.
type Transfer struct {
	schema.BaseSchema

	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	AssetID   int64     `json:"asset_id" db:"asset_id"`
	Direction string    `json:"direction" db:"direction"`
	Status    string    `json:"status" db:"status"`
	Amount    float64   `json:"amount" db:"amount"`
	TxRef     string    `json:"tx_ref" db:"tx_ref"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Transfer) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.Int64("asset_id").WithRequired(),
		schema.String("direction").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("deposit", "Deposit", "withdrawal", "Withdrawal")...),
		schema.String("status").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("pending", "Pending", "completed", "Completed", "failed", "Failed")...),
		schema.Decimal("amount").WithDefault(0),
		schema.String("tx_ref").WithMaxLength(128),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Transfer) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "transfers",
		VerboseName:       "Transfer",
		VerboseNamePlural: "Transfers",
		Indexes: []schema.Index{
			{Name: "idx_transfers_user", Fields: []string{"user_id"}},
			{Name: "idx_transfers_asset", Fields: []string{"asset_id"}},
		},
	}
}

func (Transfer) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("asset_id", "Asset").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("transfers"),
	}
}

func (Transfer) Hooks() *schema.ModelHooks { return nil }

func RegisterModels() {
	_ = registry.RegisterModel(&Wallet{})
	_ = registry.RegisterModel(&Transfer{})
}
