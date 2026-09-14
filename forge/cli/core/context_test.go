package core

import (
	"testing"

	"github.com/forgego/forge/config"
	"github.com/stretchr/testify/require"
)

func TestNewContextPreservesCallerConfig(t *testing.T) {
	cfg := config.NewConfig()
	cfg.Set("database.driver", "sqlite")

	ctx := NewContext(cfg)

	require.Same(t, cfg, ctx.Config)
	require.Equal(t, "sqlite", ctx.Config.GetDriver())
}
