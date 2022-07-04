package config

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
	ChannelInfo         []ChannelInfo
}

type ChannelInfo struct {
	ChainName   string
	ChannelID   string
}
