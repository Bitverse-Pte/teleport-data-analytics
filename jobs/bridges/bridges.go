package bridges

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/jobs/datas"
)

type Bridges struct {
	log            *logrus.Logger
	DB             *gorm.DB
	BridgeTokenMap map[string]*BridgeToken                 //Total bridges token
	TokenMap       map[string]map[string]chains.TokenQuery // map[chainName][tokenAddress]chains.TokenQuery
	AddressToName  map[string]string                       // fmt.Sprintf("%v/%v", chainName, .Address)
}

type BridgeToken struct {
	ChainAToken chains.TokenQuery
	ChainBToken chains.TokenQuery
}

type ReconciliationPacket struct {
	SrcChain    string
	DestChain   string
	SrcHeight   uint64
	DestHeight  uint64
	SrcAddress  string
	DestAddress string
}

type SumAmount struct {
	Amount float64
}

func NewBridge(log *logrus.Logger, db *gorm.DB, bridges map[string]datas.BridgeInfo, chainMap map[string]chains.BlockChain) *Bridges {
	bridgeTokenMap := make(map[string]*BridgeToken)
	addressToName := make(map[string]string)
	tokenMap := make(map[string]map[string]chains.TokenQuery)
	for _, bridge := range bridges {
		srcChainID := bridge.SrcChain.ChainId
		chainA := chainMap[srcChainID]
		if chainA == nil {
			log.Infof("chainID: %v not support", srcChainID)
			continue
		}
		chainACli := chainA.(*chains.Evm)
		destChainID := bridge.DestChain.ChainId
		chainB := chainMap[destChainID]
		if chainB == nil {
			log.Infof("chainID: %v not support", destChainID)
			continue
		}
		chainBCli := chainB.(*chains.Evm)
		if bridge.SrcChain.IsTele == nil {
			continue
		}
		for _, bridgeToken := range bridge.Tokens {
			chainAToken := chains.NewTokenQuery(chainACli, bridgeToken.SrcToken.Address, bridgeToken.Name)
			chainBToken := chains.NewTokenQuery(chainBCli, bridgeToken.DestToken.Address, bridgeToken.Name)
			bridgeTokenMap[fmt.Sprintf("%v/%v/%v", chains.TeleportChain, chainB.ChainName(), bridgeToken.Name)] = &BridgeToken{chainAToken, chainBToken}
			addressToName[fmt.Sprintf("%v/%v", bridge.SrcChain.Name, bridgeToken.SrcToken.Address)] = bridgeToken.Name
			addressToName[fmt.Sprintf("%v/%v", bridge.DestChain.Name, bridgeToken.DestToken.Address)] = bridgeToken.Name
			if tokenMap[bridge.SrcChain.Name] != nil {
				tokenMap[bridge.SrcChain.Name][bridgeToken.SrcToken.Address] = chainAToken
			} else {
				tokenMap[bridge.SrcChain.Name] = map[string]chains.TokenQuery{bridgeToken.SrcToken.Address: chainAToken}
			}
			if tokenMap[bridge.DestChain.Name] != nil {
				tokenMap[bridge.DestChain.Name][bridgeToken.DestToken.Address] = chainBToken
			} else {
				tokenMap[bridge.DestChain.Name] = map[string]chains.TokenQuery{bridgeToken.DestToken.Address: chainBToken}
			}
		}
	}
	return &Bridges{
		log:            log,
		DB:             db,
		BridgeTokenMap: bridgeTokenMap,
		AddressToName:  addressToName,
		TokenMap:       tokenMap,
	}
}

func (r *Bridges) GetTokenNameByAddress(chainName, address string) string {
	return r.AddressToName[fmt.Sprintf("%v/%v", chainName, address)]
}

func (r *Bridges) GetBridgeTokenDecimals(chainName, tokenName, destChainName string) (uint8, uint8, error) {
	bridgeTokenMapKey := fmt.Sprintf("%v/%v/%v", chains.TeleportChain, chainName, tokenName)
	if chainName == chains.TeleportChain {
		bridgeTokenMapKey = fmt.Sprintf("%v/%v/%v", chains.TeleportChain, destChainName, tokenName)
	}
	bridgeToken := r.BridgeTokenMap[bridgeTokenMapKey]
	if bridgeToken == nil {
		return 0, 0, fmt.Errorf("invalid token,bridgeTokenMapKey:%v", bridgeTokenMapKey)
	}
	decimalsA, err := bridgeToken.ChainAToken.GetDecimals()
	if err != nil {
		return 0, 0, err
	}
	decimalsB, err := bridgeToken.ChainBToken.GetDecimals()
	if err != nil {
		return 0, 0, err
	}
	return decimalsA, decimalsB, nil
}

func (r *Bridges) GetSingleTokenDecimals(chainName, tokenAddress string) (uint8, error) {
	if r.TokenMap[chainName] == nil {
		return 0, fmt.Errorf("invalid chainName")
	}
	tokenQuery := r.TokenMap[chainName][strings.ToLower(tokenAddress)]
	if tokenQuery == nil {
		return 0, fmt.Errorf("invalid token")
	}
	decimals, err := tokenQuery.GetDecimals()
	if err != nil {
		return 0, fmt.Errorf("tokenQuery.GetDecimals error:%+v", err)
	}
	return decimals, nil
}

func (r *Bridges) GetSrcChainTokenAmount(chainName, address, destChainName string) (uint8, uint8, error) {
	tokenName := r.AddressToName[fmt.Sprintf("%v/%v", chainName, strings.ToLower(address))]
	bridgeTokenMapKey := fmt.Sprintf("%v/%v/%v", chains.TeleportChain, chainName, tokenName)
	if chainName == chains.TeleportChain {
		bridgeTokenMapKey = fmt.Sprintf("%v/%v/%v", chains.TeleportChain, destChainName, tokenName)
	}
	bridgeToken := r.BridgeTokenMap[bridgeTokenMapKey]
	decimalsA, err := bridgeToken.ChainAToken.GetDecimals()
	if err != nil {
		return 0, 0, err
	}
	decimalsB, err := bridgeToken.ChainBToken.GetDecimals()
	if err != nil {
		return 0, 0, err
	}
	return decimalsA, decimalsB, nil
}

func (r *Bridges) getSumAmount(rPacket ReconciliationPacket) (float64, error) {
	var sumAmount SumAmount
	// unReceived packet
	// base：(src_height <= rPacket.SrcHeight)
	// case 1: status 2 && dest_height > rPacket.DestHeight
	// case 2: status == 1 || status == 3
	// case 3: status == 4 && refund_height > rPacket.SrcHeight
	if err := r.DB.Table("cross_chain_transactions").Select("sum(amount_float) as amount").
		Where("src_chain = ? and dest_chain = ? and send_token_address = ? ", rPacket.SrcChain, rPacket.DestChain, rPacket.SrcAddress).
		Where("src_height <= ? and ((status = 2 and dest_height > ?) or (status = 1 or status = 3) or (status = 4 and refund_height > ?)) ", rPacket.SrcHeight, rPacket.DestHeight, rPacket.SrcHeight).
		Find(&sumAmount).Error; err != nil {
		return 0, err
	}
	return sumAmount.Amount, nil
}

func (r *Bridges) GetDestReceivedSumAmount(destChain, receiveTokenAddr string, interval time.Duration) (float64, error) {
	if destChain == "" || receiveTokenAddr == "" || interval <= 0 {
		return 0, fmt.Errorf("invalid param,destChain:%v,sendTokenAddr:%v,interval:%v", destChain, receiveTokenAddr, interval)
	}
	t := time.Now().Add(-1 * interval)
	var sumAmount SumAmount
	if err := r.DB.Table("cross_chain_transactions").Select("sum(dest_amount) as amount").
		Where(" dest_chain = ? and receive_token_address = ? ", destChain, receiveTokenAddr).
		Where(" send_tx_time > ?", t).
		Find(&sumAmount).Error; err != nil {
		return 0, err
	}
	return sumAmount.Amount, nil
}

func (r *Bridges) GetDestPendingSumAmount(destChain, receiveTokenAddr string, interval time.Duration) (float64, error) {
	if destChain == "" || receiveTokenAddr == "" || interval <= 0 {
		return 0, fmt.Errorf("invalid param,destChain:%v,sendTokenAddr:%v,interval:%v", destChain, receiveTokenAddr, interval)
	}
	t := time.Now().Add(-1 * interval)
	var sumAmount SumAmount
	if err := r.DB.Table("cross_chain_transactions").Select("sum(dest_amount) as amount").
		Where(" dest_chain = ? and receive_token_address = ? ", destChain, receiveTokenAddr).
		Where(" send_tx_time > ?", t).
		Find(&sumAmount).Error; err != nil {
		return 0, err
	}
	return sumAmount.Amount, nil
}
