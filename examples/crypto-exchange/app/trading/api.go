package trading

import (
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/orm"
)

// RegisterAPI registers trading API endpoints.
func RegisterAPI(router *api.EnhancedRouter) {
	serializer := func() api.Serializer { return api.NewBaseSerializer(nil) }

	orderQS, _ := OrderObjects.Filter(&orm.Q{})
	orderVS := api.CreateDefaultViewSet(serializer, orderQS, &Order{})
	router.RegisterEnhanced("orders", orderVS)

	tradeQS, _ := TradeObjects.Filter(&orm.Q{})
	tradeVS := api.CreateDefaultViewSet(serializer, tradeQS, &Trade{})
	router.RegisterEnhanced("trades", tradeVS)
}
