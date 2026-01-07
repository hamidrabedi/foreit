package trading

import (
	"time"

	"github.com/forgego/forge/registry"
	"github.com/forgego/forge/schema"
)

// Order represents a spot order.
type Order struct {
	schema.BaseSchema

	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	PairID         int64     `json:"pair_id" db:"pair_id"`
	Side           string    `json:"side" db:"side"`
	OrderType      string    `json:"order_type" db:"order_type"`
	Status         string    `json:"status" db:"status"`
	Price          float64   `json:"price" db:"price"`
	Quantity       float64   `json:"quantity" db:"quantity"`
	FilledQuantity float64   `json:"filled_quantity" db:"filled_quantity"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (Order) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("user_id").WithRequired(),
		schema.Int64("pair_id").WithRequired(),
		schema.String("side").WithRequired().WithMaxLength(10).
			WithChoices(schema.WithChoices("buy", "Buy", "sell", "Sell")...),
		schema.String("order_type").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("market", "Market", "limit", "Limit")...),
		schema.String("status").WithRequired().WithMaxLength(20).
			WithChoices(schema.WithChoices("open", "Open", "filled", "Filled", "canceled", "Canceled")...),
		schema.Decimal("price").WithDefault(0),
		schema.Decimal("quantity").WithDefault(0),
		schema.Decimal("filled_quantity").WithDefault(0),
		schema.Time("created_at").WithAutoNowAdd(),
		schema.Time("updated_at").WithAutoNow(),
	}
}

func (Order) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "orders",
		VerboseName:       "Order",
		VerboseNamePlural: "Orders",
		Indexes: []schema.Index{
			{Name: "idx_orders_user", Fields: []string{"user_id"}},
			{Name: "idx_orders_pair", Fields: []string{"pair_id"}},
			{Name: "idx_orders_status", Fields: []string{"status"}},
		},
	}
}

func (Order) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("pair_id", "TradingPair").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("orders"),
	}
}

func (Order) Hooks() *schema.ModelHooks { return nil }

// Trade represents a matched trade.
type Trade struct {
	schema.BaseSchema

	ID          int64     `json:"id" db:"id"`
	PairID      int64     `json:"pair_id" db:"pair_id"`
	BuyOrderID  int64     `json:"buy_order_id" db:"buy_order_id"`
	SellOrderID int64     `json:"sell_order_id" db:"sell_order_id"`
	Price       float64   `json:"price" db:"price"`
	Quantity    float64   `json:"quantity" db:"quantity"`
	ExecutedAt  time.Time `json:"executed_at" db:"executed_at"`
}

func (Trade) Fields() []schema.Field {
	return []schema.Field{
		schema.Int64("id").WithPrimary().WithAutoIncrement(),
		schema.Int64("pair_id").WithRequired(),
		schema.Int64("buy_order_id").WithRequired(),
		schema.Int64("sell_order_id").WithRequired(),
		schema.Decimal("price").WithDefault(0),
		schema.Decimal("quantity").WithDefault(0),
		schema.Time("executed_at").WithAutoNowAdd(),
	}
}

func (Trade) Meta() schema.Meta {
	return schema.Meta{
		TableName:         "trades",
		VerboseName:       "Trade",
		VerboseNamePlural: "Trades",
		Indexes: []schema.Index{
			{Name: "idx_trades_pair", Fields: []string{"pair_id"}},
			{Name: "idx_trades_buy", Fields: []string{"buy_order_id"}},
			{Name: "idx_trades_sell", Fields: []string{"sell_order_id"}},
		},
	}
}

func (Trade) Relations() []schema.Relation {
	return []schema.Relation{
		schema.ForeignKey("pair_id", "TradingPair").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("trades"),
		schema.ForeignKey("buy_order_id", "Order").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("buy_trades"),
		schema.ForeignKey("sell_order_id", "Order").WithOnDelete(schema.CascadePROTECT).WithRequired().WithRelatedName("sell_trades"),
	}
}

func (Trade) Hooks() *schema.ModelHooks { return nil }

func RegisterModels() {
	_ = registry.RegisterModel(&Order{})
	_ = registry.RegisterModel(&Trade{})
}
