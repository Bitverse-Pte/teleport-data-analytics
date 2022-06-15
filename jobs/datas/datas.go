package datas

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/teleport-network/teleport-data-analytics/chains"

	"github.com/sirupsen/logrus"

	"github.com/go-co-op/gocron"

	types "github.com/teleport-network/teleport-data-analytics/types"
	cosmostypes "github.com/teleport-network/teleport-data-analytics/types/cosmos"
	evmtypes "github.com/teleport-network/teleport-data-analytics/types/evm"
)

type Token struct {
	// Name corresponds to the JSON schema field "name".
	Name string `json:"name" yaml:"name"`

	// PairType  corresponds to the JSON schema field "name".
	PairType string `json:"pairType" yaml:"pairType"`

	// DestToken corresponds to the JSON schema field "destToken".
	DestToken types.TokenInfo `json:"destToken" yaml:"destToken"`

	// RelayToken
	RelayToken string `json:"relayToken" yaml:"relayToken"`

	// SrcToken corresponds to the JSON schema field "srcToken".
	SrcToken types.TokenInfo `json:"srcToken" yaml:"srcToken"`
}

type BridgeInfo struct {
	// DestChain corresponds to the JSON schema field "destChain".
	DestChain types.BridgeChain `json:"destChain,omitempty" yaml:"destChain,omitempty"`

	// SrcChain corresponds to the JSON schema field "srcChain".
	SrcChain types.BridgeChain `json:"srcChain,omitempty" yaml:"srcChain,omitempty"`

	// AgentAddress
	AgentAddress string `json:"agent_address"`

	// Tokens corresponds to the JSON schema field "tokens".
	Tokens []Token `json:"tokens,omitempty" yaml:"tokens,omitempty"`
}

var (
	// listing all chians
	ChainList types.Chains
	// mapping the counterparty chain list for each chain
	CounterpartyChains = make(map[string]types.Chains)
	// mapping token list for each two cross-chain
	// {src_chain}/{dest_chain} => Tokens
	Bridges = make(map[string]BridgeInfo)

	TokenMap = make(map[string]types.TokenInfo)
)

func ReloadState(chainType string, period time.Duration) {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(period * time.Minute).Do(func() {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("syncToDB panic:%+v",err)
			}
		}()
		if err := LoadState(chainType); err != nil {
			logrus.Errorf("LoadState Error:%+v", err)
		}
	})
	scheduler.StartAsync()
}

// loadState loads the latest data
func LoadState(chainType string) error {
	var tmpBridges = new(types.Bridges)
	var tmpEvmChains = new(evmtypes.ChainlistSchemaJson)
	var tmpCosmosChains = new(cosmostypes.CosmosChains)
	var tmpTokens = new(types.TokenlistSchemaJson)
	var chainList types.Chains
	counterpartyChains := make(map[string]types.Chains)
	bridges := make(map[string]BridgeInfo)
	switch chainType {
	case "qanet":
		network := "QA"
		jsonBridges, evmChains, cosmosChains, tokens, err := syncDatas(network)
		if err != nil {
			return fmt.Errorf("failed to load data from teleport-bridge-lists: %s, network: %s", err.Error(), network)
		}
		if err := tmpBridges.UnmarshalJSON(jsonBridges); err != nil {
			return fmt.Errorf("failed to load bridges, err: %s", err.Error())
		}
		if err := tmpEvmChains.UnmarshalJSON(evmChains); err != nil {
			return fmt.Errorf("failed to load chains, err: %s", err.Error())
		}
		if err := tmpCosmosChains.UnmarshalJSON(cosmosChains); err != nil {
			return fmt.Errorf("failed to load chains, err: %s", err.Error())
		}
		if err := tmpTokens.UnmarshalJSON(tokens); err != nil {
			return fmt.Errorf("failed to load tokens, err: %s", err.Error())
		}
	case "testnet":
		network := "Testnet"
		jsonBridges, evmChains, cosmosChains, tokens, err := syncDatas(network)
		if err != nil {
			return fmt.Errorf("failed to load data from teleport-bridge-lists: %s, network: %s", err.Error(), network)
		}
		if err := tmpBridges.UnmarshalJSON(jsonBridges); err != nil {
			return fmt.Errorf("failed to load bridges, err: %s", err.Error())
		}
		if err := tmpEvmChains.UnmarshalJSON(evmChains); err != nil {
			return fmt.Errorf("failed to load chains, err: %s", err.Error())
		}
		if err := tmpCosmosChains.UnmarshalJSON(cosmosChains); err != nil {
			return fmt.Errorf("failed to load chains, err: %s", err.Error())
		}
		if err := tmpTokens.UnmarshalJSON(tokens); err != nil {
			return fmt.Errorf("failed to load tokens, err: %w", err)
		}
	default:
		return fmt.Errorf("invalid chain type %s", chainType)
	}
	tmpChainsMap := make(map[string]types.ChainI)
	var teleportChainID string
	var teleportChain types.EvmChain
	for _, chain := range tmpCosmosChains.Chains {
		chainId, ok := ConverTeleportChainID(chain.ChainId)
		if ok {
			tmpChainsMap[chain.ChainId] = types.TeleportChain(chain)
			chainList = append(chainList, types.TeleportChain(chain))
			teleportChainID = chainId
		} else {
			tmpChainsMap[chain.ChainId] = types.CosmosChain(chain)
			chainList = append(chainList, types.CosmosChain(chain))
		}
	}

	for _, chain := range tmpEvmChains.Chains {
		tmpChainsMap[fmt.Sprintf("%d", chain.ChainId)] = types.EvmChain(chain)
		chainList = append(chainList, types.EvmChain(chain))
		if strconv.Itoa(chain.ChainId) == teleportChainID {
			teleportChain = types.EvmChain(chain)
		}
	}
	tmpTokensMap := make(map[string]types.TokenInfo)
	for _, token := range tmpTokens.Tokens {
		tmpTokensMap[fmt.Sprintf("%s/%s", token.ChainId, token.Address)] = token
	}
	for _, bridge := range tmpBridges.Bridges {
		if _, ok := tmpChainsMap[bridge.SrcChain.ChainId]; !ok {
			continue
		}
		if _, ok := tmpChainsMap[bridge.DestChain.ChainId]; !ok {
			continue
		}
		if srcChainId, isTeleport := ConverTeleportChainID(bridge.SrcChain.ChainId); isTeleport {
			counterpartyChains[srcChainId] = append(counterpartyChains[srcChainId], tmpChainsMap[bridge.DestChain.ChainId])
		}
		if destChainId, isTeleport := ConverTeleportChainID(bridge.DestChain.ChainId); isTeleport {
			counterpartyChains[destChainId] = append(counterpartyChains[destChainId], tmpChainsMap[bridge.SrcChain.ChainId])
		}
		counterpartyChains[bridge.SrcChain.ChainId] = append(counterpartyChains[bridge.SrcChain.ChainId], tmpChainsMap[bridge.DestChain.ChainId])
		counterpartyChains[bridge.DestChain.ChainId] = append(counterpartyChains[bridge.DestChain.ChainId], tmpChainsMap[bridge.SrcChain.ChainId])
		chainPair := fmt.Sprintf("%s/%s", bridge.SrcChain.ChainId, bridge.DestChain.ChainId)
		reversedChainPair := fmt.Sprintf("%s/%s", bridge.DestChain.ChainId, bridge.SrcChain.ChainId)
		var tokens []Token
		var reversedTokens []Token
		for _, token := range bridge.Tokens {
			relayToken := ""
			if token.RelayToken != nil {
				relayToken = *token.RelayToken
			}
			tokens = append(tokens, Token{
				Name:       token.Name,
				PairType:   token.PairType,
				SrcToken:   tmpTokensMap[fmt.Sprintf("%s/%s", bridge.SrcChain.ChainId, token.SrcToken)],
				DestToken:  tmpTokensMap[fmt.Sprintf("%s/%s", bridge.DestChain.ChainId, token.DestToken)],
				RelayToken: relayToken,
			})
			pairTypes := strings.Split(token.PairType, "-")
			var reversedTokenPairType string
			if len(pairTypes) == 2 {
				reversedTokenPairType = fmt.Sprintf("%s-%s", pairTypes[1], pairTypes[0])
			}
			reversedTokens = append(reversedTokens, Token{
				Name:       token.Name,
				PairType:   reversedTokenPairType,
				SrcToken:   tmpTokensMap[fmt.Sprintf("%s/%s", bridge.DestChain.ChainId, token.DestToken)],
				DestToken:  tmpTokensMap[fmt.Sprintf("%s/%s", bridge.SrcChain.ChainId, token.SrcToken)],
				RelayToken: relayToken,
			})

		}
		agentAddress := ""
		if bridge.AgentAddress != nil {
			agentAddress = *bridge.AgentAddress
		}
		bridges[chainPair] = BridgeInfo{
			SrcChain:     bridge.SrcChain,
			DestChain:    bridge.DestChain,
			AgentAddress: agentAddress,
			Tokens:       tokens,
		}
		bridges[reversedChainPair] = BridgeInfo{
			SrcChain:     bridge.DestChain,
			DestChain:    bridge.SrcChain,
			AgentAddress: agentAddress,
			Tokens:       reversedTokens,
		}
	}

	for key, counterpartyChain := range counterpartyChains {
		if IsCosmosChainId(key) {
			for i, v := range counterpartyChain {
				if v.Type() == chains.TeleportChain {
					counterpartyChains[key][i] = teleportChain
				}
			}
		}
	}
	ChainList = chainList
	CounterpartyChains = counterpartyChains
	Bridges = bridges
	return nil
}

func ConverTeleportChainID(chainId string) (string, bool) {
	var IsTeleport bool
	chainIdSliA := strings.Split(chainId, "_")
	if len(chainIdSliA) == 2 && chainIdSliA[0] == chains.TeleportChain {
		chainIdSliB := strings.Split(chainIdSliA[1], "-")
		if len(chainIdSliB) == 2 {
			chainId = chainIdSliB[0]
			IsTeleport = true
		} else {
			logrus.Errorf("invalid teleport chain id %s", chainId)
		}
	}
	return chainId, IsTeleport
}

func IsCosmosChainId(chainId string) bool {
	_, err := strconv.Atoi(chainId)
	if err != nil {
		_, isTeleChainId := ConverTeleportChainID(chainId)
		if !isTeleChainId {
			return true
		}
	}
	return false
}
