package bridges

import "math/big"

type PacketToken struct {
	AddressHex  string
	Coefficient int64
	Amount      *big.Int
	Height      *big.Int
}

type Chain struct {
	Name   string
	Tokens []ChainToken
}
type ChainToken struct {
	Name         string
	Address      string
	ABI          string
	ChainName    string
	EndPointAddr string
}
