package datas

import (
	"fmt"

	"github.com/teleport-network/teleport-data-analytics/config"

	"github.com/teleport-network/teleport-data-analytics/tools"
)

func syncDatas(network string) (bridges, evmChains, cosmosChains, tokens []byte, err error) {
	evmChainListUrl := fmt.Sprintf(config.EvmChainsFormat, network)
	cosmosChainListUrl := fmt.Sprintf(config.CosmosChainsFormat, network)
	bridgeListUrl := fmt.Sprintf(config.BridgesFormat, network)
	tokenListUrl := fmt.Sprintf(config.TokensFormat, network)
	header := map[string]string{
		"Authorization": config.Authorization,
	}
	bridges, err = tools.HttpGet(bridgeListUrl, header)
	if err != nil {
		return
	}

	evmChains, err = tools.HttpGet(evmChainListUrl, header)
	if err != nil {
		return
	}

	cosmosChains, err = tools.HttpGet(cosmosChainListUrl, header)
	if err != nil {
		return
	}

	tokens, err = tools.HttpGet(tokenListUrl, header)
	if err != nil {
		return
	}
	return
}
