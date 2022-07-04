package config

type EvmConfig struct {
	EvmUrl           string
	PacketTopic      string
	AckTopic         string
	ReceivedAckTopic string
	PacketContract   string
	ChainName        string
	ChainID          string
	Frequency        int
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
