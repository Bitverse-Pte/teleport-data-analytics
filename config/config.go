package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/teleport-network/teleport-data-analytics/tools"
)

const (
	DefaultHomeDirName          = ".teleport-data-analytics"
	DefaultConfigDirName        = "configs" // TODO delete initialization.DefaultConfigDirName
	DefaultConfigName           = "config.toml"
)

var (
	Home                 string
	LocalConfig          string
	UserDir, _           = os.UserHomeDir()
	DefaultHomePath      = filepath.Join(UserDir, DefaultHomeDirName)
)

var (
	Authorization      = "token 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"
	BridgesFormat      = "https://raw.githubusercontent.com/teleport-network/teleport-bridge-lists/main/%s/bridgelist.json"
	EvmChainsFormat    = "https://raw.githubusercontent.com/teleport-network/teleport-bridge-lists/main/%s/evm.chains.json"
	CosmosChainsFormat = "https://raw.githubusercontent.com/teleport-network/teleport-bridge-lists/main/%s/cosmos.chains.json"
	TokensFormat       = "https://raw.githubusercontent.com/teleport-network/teleport-bridge-lists/main/%s/tokens.json"
)

type Status int8

type Config struct {
	MysqlAddr          string
	SyncEnable         bool
	ReconcileEnable    bool
	Network            string
	DBModel            string
	Authorization      string //option
	BridgesFormat      string //option
	EvmChainsFormat    string //option
	CosmosChainsFormat string //option
	TokensFormat       string //option
	ReloadPeriod       time.Duration
	EvmChains          []EvmConfig
	Teleport           TendermintConfig
}

type EvmConfig struct {
	EvmUrl         string
	PacketTopic    string
	AckTopic       string
	PacketContract string
	ChainName      string
	ChainID        string
	Frequency      int
	// Refer to block generation speed
	BatchNumber        uint64
	ProxyAddr          string
	ProxyTopic         string
	AgentAddr          string
	AgentTopic         string
	EndPointAddr       string
	NativeToken        string
	NativeDecimals     uint8
	RevisedHeight      uint64
	StartHeight        uint64
	BalanceMonitorings []BalanceMonitoring
}

type TendermintConfig struct {
	Url       string
	ChainName string
	ChainID   string
	AgentAddr string
	Frequency int
	// Refer to block generation speed
	BatchNumber        uint64
	RevisedHeight      uint64
	StartHeight        uint64
	BalanceMonitorings []BalanceMonitoring
}

type BalanceMonitoring struct {
	TokenName  string
	AddressHex string
	Accounts   []Account
}

type Account struct {
	AddressHex string
	Name       string
}

type AddressInfo struct {
	Name      string
	ChainName string
	Address   string
}

func LoadConfigs() *Config {
	//Default
	if Home == "" {
		Home = DefaultHomePath
	}
	if LocalConfig == "" {
		LocalConfig = filepath.Join(Home, DefaultConfigDirName, DefaultConfigName)
	}
	cfg := Config{}
	tools.InitTomlConfigs([]*tools.ConfigMap{
		{
			FilePath: LocalConfig,
			Pointer:  &cfg,
		},
	})
	if cfg.Authorization != "" {
		Authorization = cfg.Authorization
	}
	if cfg.BridgesFormat != "" && cfg.EvmChainsFormat != "" && cfg.CosmosChainsFormat != "" && cfg.TokensFormat != "" {
		BridgesFormat = cfg.BridgesFormat
		EvmChainsFormat = cfg.EvmChainsFormat
		CosmosChainsFormat = cfg.CosmosChainsFormat
		TokensFormat = cfg.TokensFormat
	}
	return &cfg
}
