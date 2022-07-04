package bridges

import (
	"fmt"
	"math"
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/model"
)

func (r *Bridges) Reconcile(teleportPacket, otherPacket ReconciliationPacket) (model.BridgeReconcileResult, error) {
	var bridgeReconcileResult model.BridgeReconcileResult
	var chainName, destChainName, address string
	if teleportPacket.SrcAddress != "" {
		chainName = teleportPacket.SrcChain
		destChainName = otherPacket.SrcChain
		address = teleportPacket.SrcAddress

	} else if otherPacket.SrcAddress != "" {
		chainName = otherPacket.SrcChain
		destChainName = teleportPacket.SrcChain
		address = otherPacket.SrcAddress
	}
	tokenName := r.AddressToName[fmt.Sprintf("%v/%v", chainName, address)]
	if tokenName == "" {
		return bridgeReconcileResult, fmt.Errorf("invalid token address,chainName:%v,tokenAddress:%v", chainName, address)
	}
	var otherAmount, defaultTokenAmount float64
	bridgeTokenMapKey := fmt.Sprintf("%v/%v/%v", chains.TeleportChain, chainName, tokenName)
	if chainName == chains.TeleportChain {
		bridgeTokenMapKey = fmt.Sprintf("%v/%v/%v", chains.TeleportChain, destChainName, tokenName)
	}
	bridgeToken := r.BridgeTokenMap[bridgeTokenMapKey]
	if bridgeToken == nil {
		return bridgeReconcileResult, fmt.Errorf("bridgeToken not found,bridgeTokenMapKey:%v", bridgeTokenMapKey)
	}
	teleportPacket.SrcAddress = bridgeToken.ChainAToken.Address()
	otherPacket.SrcAddress = bridgeToken.ChainBToken.Address()
	sumA, err := r.getSumAmount(teleportPacket)
	if err != nil {
		r.log.Errorf("get teleportPacket sum amount error:%+v", err)
		return bridgeReconcileResult, err
	}
	sumB, err := r.getSumAmount(otherPacket)
	if err != nil {
		r.log.Errorf("get %vPacket sum amount error:%+v", otherPacket.SrcChain, err)
		return bridgeReconcileResult, err
	}
	var getChainATokenAmount func(addressHex string, srcChain string, blockNumber *big.Int) (chains.TokenAmount, error)
	var getChainBTokenAmount func(blockNumber *big.Int) (chains.TokenAmount, error)
	if tokenName == "tele" {
		getChainATokenAmount = bridgeToken.ChainAToken.GetOutTokenAmount
		getChainBTokenAmount = bridgeToken.ChainBToken.GetErc20TotalSupply

	} else {
		getChainATokenAmount = bridgeToken.ChainAToken.GetInTokenAmount
		getChainBTokenAmount = bridgeToken.ChainBToken.GetTransferBalance
	}
	var coefficientA, coefficientB float64
	if bridgeToken.ChainAToken.NativeToken() == bridgeToken.ChainAToken.TokenName() {
		coefficientA = -1
		coefficientB = 1
	} else {
		coefficientA = 1
		coefficientB = -1
	}
	teleportAmount, err := getChainATokenAmount(bridgeToken.ChainAToken.Address(), bridgeToken.ChainBToken.ChainName(), big.NewInt(int64(teleportPacket.SrcHeight))) //TODO do not use
	if err != nil {
		r.log.Errorf("get chainA token amount error:%v\n,chainName:%v,tokenName:%v,tokenAddress:%v,height:%v",
			err, bridgeToken.ChainAToken.ChainName(), tokenName, bridgeToken.ChainAToken.Address(), teleportPacket.SrcHeight)
		return bridgeReconcileResult, err
	}
	r.log.Infof("teleport chain %v amount:%v", bridgeToken.ChainAToken.ChainName(), teleportAmount)
	defaultTokenAmount, err = teleportAmount.Float64()
	if err != nil {
		r.log.Errorf("teleportAmount.Float64 error:%+v", err)
		return bridgeReconcileResult, err
	}
	amount, err := getChainBTokenAmount(big.NewInt(int64(otherPacket.SrcHeight)))
	if err != nil {
		r.log.Errorf("get chainA  token amount error:%v\\n,chainName:%v,tokenName:%v,tokenAddress:%v,height:%v",
			err, bridgeToken.ChainBToken.ChainName(), tokenName, bridgeToken.ChainBToken.Address(), otherPacket.SrcHeight)
		return bridgeReconcileResult, err
	}
	r.log.Infof("teleport chain %v amount:%v", bridgeToken.ChainBToken.ChainName(), amount)
	otherAmount, err = amount.Float64()
	if err != nil {
		r.log.Errorf("amount.Float64 error:%+v", err)
		return bridgeReconcileResult, err
	}
	// decimals convert
	sumA = sumA * math.Pow10(int(teleportAmount.Decimals)-int(amount.Decimals))
	sumB = sumB * math.Pow10(int(amount.Decimals)-int(teleportAmount.Decimals))
	bridgeReconcileResult = model.BridgeReconcileResult{
		TokenName:      tokenName,
		SrcChain:       bridgeToken.ChainAToken.ChainName(),
		SrcToken:       bridgeToken.ChainAToken.Address(),
		SrcAmount:      defaultTokenAmount,
		SrcTrimAmount:  sumA,
		DestChain:      bridgeToken.ChainBToken.ChainName(),
		DestToken:      bridgeToken.ChainBToken.Address(),
		DestAmount:     otherAmount,
		DestTrimAmount: sumB,
		Result:         "success",
		ResultMsg:      "success",
	}
	decimal.DivisionPrecision = 2
	sub := decimal.NewFromFloat(defaultTokenAmount).Add(decimal.NewFromFloat(coefficientA * sumA)).Sub(decimal.NewFromFloat(otherAmount * math.Pow10(int(teleportAmount.Decimals-amount.Decimals))).Add(decimal.NewFromFloat(coefficientB * sumB * math.Pow10(int(teleportAmount.Decimals-amount.Decimals)))))
	if !sub.IsZero() {
		//config.ReconciliationStatus = config.Failed
		bridgeReconcileResult.Result = "failed"
		bridgeReconcileResult.ResultMsg = "failed"
		newSub, _ := sub.Float64()
		bridgeReconcileResult.DifferenceValue = newSub
	}
	return bridgeReconcileResult, nil
}

func (r *Bridges) BatchBridgeReconcile() {
	var crossChainTxs []model.CrossChainTransaction
	if err := r.DB.Where("src_chain = ? and src_height > ?  and dest_height <> 0 and reconcile_result = '' ", "teleport", 537835).Limit(1000).Find(&crossChainTxs).Error; err != nil {
		return
	}
	for _, crossChainTx := range crossChainTxs {
		packetA := ReconciliationPacket{
			SrcChain:   crossChainTx.SrcChain,
			DestChain:  crossChainTx.DestChain,
			SrcHeight:  crossChainTx.SrcHeight,
			DestHeight: crossChainTx.DestHeight,
			SrcAddress: crossChainTx.SendTokenAddress,
		}
		packetB := ReconciliationPacket{
			SrcChain:   crossChainTx.DestChain,
			DestChain:  crossChainTx.SrcChain,
			SrcHeight:  crossChainTx.DestHeight,
			DestHeight: crossChainTx.SrcHeight,
		}
		var teleportPacket, otherPacket ReconciliationPacket
		if crossChainTx.SrcChain == "teleport" {
			teleportPacket = packetA
			otherPacket = packetB
		} else {
			teleportPacket = packetB
			otherPacket = packetA
		}
		if crossChainTx.SrcHeight != 0 && crossChainTx.DestHeight != 0 {
			bridgeReconcileResult, err := r.Reconcile(teleportPacket, otherPacket)
			if err != nil {
				return
			}
			crossChainTransaction := model.CrossChainTransaction{
				TokenName:       bridgeReconcileResult.TokenName,
				SrcAmount:       bridgeReconcileResult.DestAmount,
				SrcTrimAmount:   bridgeReconcileResult.DestTrimAmount,
				DestAmount:      bridgeReconcileResult.SrcAmount,
				DestTrimAmount:  bridgeReconcileResult.SrcTrimAmount,
				ReconcileResult: bridgeReconcileResult.Result,
				ReconcileMsg:    bridgeReconcileResult.ResultMsg,
				DifferenceValue: bridgeReconcileResult.DifferenceValue,
			}
			crossChainTransaction.ID = crossChainTx.ID
			if bridgeReconcileResult.SrcChain != crossChainTx.SrcChain {
				crossChainTransaction.SrcAmount = bridgeReconcileResult.DestAmount
				crossChainTransaction.SrcTrimAmount = bridgeReconcileResult.DestTrimAmount
				crossChainTransaction.DestAmount = bridgeReconcileResult.SrcAmount
				crossChainTransaction.DestTrimAmount = bridgeReconcileResult.SrcTrimAmount
			}
			if err := r.DB.Updates(&crossChainTransaction).Error; err != nil {
				r.log.Errorf("update Reconcile result error:%+v", err)
				return
			}
		}

	}

}
