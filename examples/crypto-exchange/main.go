package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"

	"github.com/forgego/forge/admin"
	"github.com/forgego/forge/api"
	"github.com/forgego/forge/config"
	"github.com/forgego/forge/db"
	"github.com/forgego/forge/identity"
	forgelog "github.com/forgego/forge/log"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/server"

	"examples/crypto-exchange/app/accounts"
	"examples/crypto-exchange/app/ledger"
	"examples/crypto-exchange/app/markets"
	"examples/crypto-exchange/app/trading"
	"examples/crypto-exchange/app/wallets"
)

//go:generate forge generate --models ./app --output ./app

func main() {
	ctx := context.Background()

	cfg := config.NewConfig()
	settings := config.LoadSettings(cfg)

	logger, err := forgelog.NewLogger(settings.App.Debug)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer logger.Sync()

	database, err := db.NewDBFromConfig(cfg)
	if err != nil {
		stdlog.Fatal(err)
	}
	defer database.Close()

	accounts.RegisterModels()
	markets.RegisterModels()
	wallets.RegisterModels()
	trading.RegisterModels()
	ledger.RegisterModels()

	orm.SetDBForAllManagers(database)

	identitySystem, err := identity.SetupIdentitySystem(database, nil)
	if err != nil {
		logger.Warn("identity system setup failed", forgelog.Error(err))
	}

	admin.DefaultSite.Title = settings.Admin.Title
	uiConfig := admin.DefaultUIConfig()
	uiConfig.BasePath = settings.Admin.Path
	admin.DefaultSite.WithUIConfig(uiConfig)

	accounts.RegisterAdmin(ctx)
	markets.RegisterAdmin(ctx)
	wallets.RegisterAdmin(ctx)
	trading.RegisterAdmin(ctx)
	ledger.RegisterAdmin(ctx)

	srv, err := server.NewServer(cfg, settings, logger)
	if err != nil {
		stdlog.Fatal(err)
	}

	api.SetupDefaultAPI()

	apiPath := cfg.GetString("api.path", "/api/v1")
	apiRouter := api.NewEnhancedRouter(apiPath)
	accounts.RegisterAPI(apiRouter)
	markets.RegisterAPI(apiRouter)
	wallets.RegisterAPI(apiRouter)
	trading.RegisterAPI(apiRouter)
	ledger.RegisterAPI(apiRouter)

	srv.RegisterRoutes(func(router *server.Router) {
		router.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, _ = w.Write([]byte("Forge Crypto Exchange running"))
		})

		if settings.Admin.Enabled {
			router.Mount(settings.Admin.Path, admin.DefaultSite.Handler())
		}

		if identitySystem != nil {
			identitySystem.RegisterRoutes(router)
		}

		apiRouter.RegisterRoutesEnhanced(router)
	})

	fmt.Printf("Forge Crypto Exchange running on %s:%s\n", settings.Server.Host, settings.Server.Port)
	if err := srv.Start(); err != nil {
		stdlog.Fatal(err)
	}
}
