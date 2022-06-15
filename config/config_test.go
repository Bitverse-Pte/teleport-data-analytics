package config

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/teleport-network/teleport-data-analytics/tools"
)

func TestLoadConfigs(t *testing.T) {
	cfg := Config{}
	tools.InitTomlConfigs([]*tools.ConfigMap{
		{
			FilePath: "./test.toml",
			Pointer:  &cfg,
		},
	})
	require.True(t, cfg.MysqlAddr != "")
	require.True(t, len(cfg.FreeChains) == 2)
}
