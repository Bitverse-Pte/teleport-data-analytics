package monitoring

import (
	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/config"
)

type BalanceMonitoring struct {
	TokensQuery chains.TokenQuery
	Chains      map[string]chains.BlockChain
	Accounts    []config.Account
}

func NewBalanceMonitoring(tokensQuery chains.TokenQuery, accounts []config.Account) *BalanceMonitoring {
	return &BalanceMonitoring{TokensQuery: tokensQuery, Accounts: accounts}
}
