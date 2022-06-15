package types

import (
	"encoding/json"

	"github.com/teleport-network/teleport-data-analytics/types/cosmos"
	"github.com/teleport-network/teleport-data-analytics/types/evm"
)

type Chains []ChainI

type ChainI interface {
	Type() string
}

type EvmChain evm.Chain
type CosmosChain cosmos.Chain
type TeleportChain cosmos.Chain

func (EvmChain) Type() string {
	return "evm"
}

func (ec EvmChain) MarshalJSON() ([]byte, error) {
	chain := Chain{
		Type:  ec.Type(),
		Value: evm.Chain(ec),
	}
	return json.Marshal(chain)
}

func (CosmosChain) Type() string {
	return "cosmos"
}

func (cc CosmosChain) MarshalJSON() ([]byte, error) {
	chain := Chain{
		Type:  cc.Type(),
		Value: cosmos.Chain(cc),
	}
	return json.Marshal(chain)
}

func (TeleportChain) Type() string {
	return "teleport"
}

func (tc TeleportChain) MarshalJSON() ([]byte, error) {
	chain := Chain{
		Type:  tc.Type(),
		Value: cosmos.Chain(tc),
	}
	return json.Marshal(chain)
}

type Chain struct {
	Type  string      `json:"type" yaml:"type"`
	Value interface{} `json:"value" yaml:"value"`
}
