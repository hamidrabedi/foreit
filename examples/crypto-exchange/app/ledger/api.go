package ledger

import (
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/orm"
)

// RegisterAPI registers ledger API endpoints.
func RegisterAPI(router *api.EnhancedRouter) {
	serializer := func() api.Serializer { return api.NewBaseSerializer(nil) }

	entryQS, _ := LedgerEntryObjects.Filter(&orm.Q{})
	entryVS := api.CreateDefaultViewSet(serializer, entryQS, &LedgerEntry{})
	router.RegisterEnhanced("ledger", entryVS)
}
