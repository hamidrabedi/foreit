package markets

import (
	"time"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Asset represents a tradeable asset.
type Asset struct {
	schema.BaseSchema

	ID        int64     `json:"id" db:"id"`
	Symbol    string    `json:"symbol" db:"symbol"`
	Name      string    `json:"name" db:"name"`
	AssetType string    `json:"asset_type" db:"asset_type"`
	Precision int32     `json:"precision" db:"precision"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (Asset) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("symbol").WithRequired().WithMaxLength(20).WithUnique(),
		schema.String("name").WithRequired().WithMaxLength(100),
		schema.String("asset_type").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("crypto", "Crypto", "fiat", "Fiat")...),
		schema.Int32("precision").WithDefault(8),
		schema.Bool("is_active").WithDefault(true),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Asset) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "assets",
		VerboseName:       "Asset",
		VerboseNamePlural: "Assets",
		Indexes: []schema.Index{
			{Name: "idx_assets_symbol", Fields: []string{"symbol"}, Unique: true},
		},
	}
}

func (Asset) Relations() []schema.Relation { return []schema.Relation{} }

func (Asset) Hooks() *schema.ModelHooks { return nil }

// TradingPair represents a market pair.
type TradingPair struct {
	schema.BaseSchema

	ID           int64     `json:"id" db:"id"`
	Symbol       string    `json:"symbol" db:"symbol"`
	BaseAssetID  int64     `json:"base_asset_id" db:"base_asset_id"`
	QuoteAssetID int64     `json:"quote_asset_id" db:"quote_asset_id"`
	Status       string    `json:"status" db:"status"`
	MinPrice     float64   `json:"min_price" db:"min_price"`
	MaxPrice     float64   `json:"max_price" db:"max_price"`
	PriceTick    float64   `json:"price_tick" db:"price_tick"`
	QtyStep      float64   `json:"qty_step" db:"qty_step"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (TradingPair) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.String("symbol").WithRequired().WithMaxLength(20).WithUnique(),
		schema.Int64("base_asset_id").WithRequired(),
		schema.Int64("quote_asset_id").WithRequired(),
		schema.String("status").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("active", "Active", "halted", "Halted")...),
		schema.Decimal("min_price").WithDefault(0),
		schema.Decimal("max_price").WithDefault(0),
		schema.Decimal("price_tick").WithDefault(0.01),
		schema.Decimal("qty_step").WithDefault(0.0001),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (TradingPair) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "trading_pairs",
		VerboseName:       "Trading Pair",
		VerboseNamePlural: "Trading Pairs",
		Indexes: []schema.Index{
			{Name: "idx_pairs_symbol", Fields: []string{"symbol"}, Unique: true},
			{Name: "idx_pairs_base", Fields: []string{"base_asset_id"}},
			{Name: "idx_pairs_quote", Fields: []string{"quote_asset_id"}},
		},
	}
}

func (TradingPair) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("base_asset_id", "Asset").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("base_pairs"),
		schema.ForeignKey("quote_asset_id", "Asset").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("quote_pairs"),
	}
}

func (TradingPair) Hooks() *schema.ModelHooks { return nil }

func RegisterModels() {
	_ = registry.RegisterModel(&Asset{})
	_ = registry.RegisterModel(&TradingPair{})
}
