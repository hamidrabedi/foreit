package accounts

import (
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/orm"
)

// RegisterAPI registers account API endpoints.
func RegisterAPI(router *api.EnhancedRouter) {
	serializer := func() api.Serializer { return api.NewBaseSerializer(nil) }

	profileQS, _ := UserProfileObjects.Filter(&orm.Q{})
	profileVS := api.CreateDefaultViewSet(serializer, profileQS, &UserProfile{})
	router.RegisterEnhanced("profiles", profileVS)

	keyQS, _ := APIKeyObjects.Filter(&orm.Q{})
	keyVS := api.CreateDefaultViewSet(serializer, keyQS, &APIKey{})
	router.RegisterEnhanced("api-keys", keyVS)
}
