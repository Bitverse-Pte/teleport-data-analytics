package types

import "github.com/teleport-network/teleport-data-analytics/model"

// bridge tx
type BridgeTx struct {
	SrcChain string
	Txs      map[string][]model.CrossChainTransaction
}
