package markets

import (
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/orm"
)

// RegisterAPI registers market API endpoints.
func RegisterAPI(router *api.EnhancedRouter) {
	serializer := func() api.Serializer { return api.NewBaseSerializer(nil) }

	assetQS, _ := AssetObjects.Filter(&orm.Q{})
	assetVS := api.CreateDefaultViewSet(serializer, assetQS, &Asset{})
	router.RegisterEnhanced("assets", assetVS)

	pairQS, _ := TradingPairObjects.Filter(&orm.Q{})
	pairVS := api.CreateDefaultViewSet(serializer, pairQS, &TradingPair{})
	router.RegisterEnhanced("pairs", pairVS)
}
