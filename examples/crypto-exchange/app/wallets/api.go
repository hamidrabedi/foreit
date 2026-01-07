package wallets

import (
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/orm"
)

// RegisterAPI registers wallet API endpoints.
func RegisterAPI(router *api.EnhancedRouter) {
	serializer := func() api.Serializer { return api.NewBaseSerializer(nil) }

	walletQS, _ := WalletObjects.Filter(&orm.Q{})
	walletVS := api.CreateDefaultViewSet(serializer, walletQS, &Wallet{})
	router.RegisterEnhanced("wallets", walletVS)

	transferQS, _ := TransferObjects.Filter(&orm.Q{})
	transferVS := api.CreateDefaultViewSet(serializer, transferQS, &Transfer{})
	router.RegisterEnhanced("transfers", transferVS)
}
