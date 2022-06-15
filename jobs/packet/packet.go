package packet

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/metrics"
	"github.com/teleport-network/teleport-data-analytics/model"
	"github.com/teleport-network/teleport-data-analytics/jobs/bridges"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/sirupsen/logrus"
	xibceth "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	xibctendermint "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

type PacketPool struct {
	Codec             *codec.ProtoCodec
	DB                *gorm.DB
	log               *logrus.Logger
	Chains            map[string]chains.BlockChain
	ChainMap          map[string]string
	IDToName          map[string]string
	NameToID          map[string]string
	ReconcileEnable   bool
	ReconciliationCli *bridges.Bridges
	MetricsManager    *metrics.MetricManager
}

const (
	Packet           = "packet"
	Ack              = "ack"
)

func NewPacketDBPool(db *gorm.DB, log *logrus.Logger, chain map[string]chains.BlockChain, chainMap map[string]string, reconciliationCli *bridges.Bridges, reconcileEnable bool, metricsManager *metrics.MetricManager) *PacketPool {
	cdc := makeCodec()
	if err := db.AutoMigrate(&model.SyncState{}, &model.CrossPacket{}, &model.PacketRelation{}, &model.BridgeReconcileResult{}, &model.Record{}, &model.BlockReconcileResult{}, &model.CrossChainTransaction{}); err != nil {
		panic(fmt.Errorf("db.AutoMigrate:%+v", err))
	}
	return &PacketPool{
		Codec:             cdc,
		DB:                db,
		log:               log,
		Chains:            chain,
		ChainMap:          chainMap,
		ReconcileEnable:   reconcileEnable,
		ReconciliationCli: reconciliationCli,
		MetricsManager:    metricsManager,
	}
}

var (
	_ = chains.BlockChain(&chains.Evm{})
	_ = chains.BlockChain(&chains.Teleport{})
)

func (p *PacketPool) SyncToDB(s *gocron.Scheduler, syncEnable bool) {
	for chainName, chain := range p.Chains {
		c := chain
		name := chainName
		syncState := model.SyncState{ChainName: name}
		if err := p.DB.Model(&syncState).Find(&syncState).Error; err != nil {
			p.log.Fatalf("get syncState error:%+v", err)
		}
		if chain.RevisedHeight() != 0 {
			p.log.Infof("chain %v database sync height:%v,reset height:%v", name, syncState.Height, chain.RevisedHeight())
			syncState.Height = chain.RevisedHeight()
			if err := p.DB.Save(&syncState).Error; err != nil {
				panic(err)
			}
		} else if syncState.Height == 0 {
			syncState.Height = chain.StartHeight()
			if err := p.DB.Save(&syncState).Error; err != nil {
				panic(err)
			}
		}

		if c.GetFrequency() <= 0 {
			p.log.Fatalf("invalid frequency:%v", c.GetFrequency())
		}
		if c.GetBatchNumber() == 0 {
			p.log.Fatalf("invalid batchNumber:%v", c.GetBatchNumber())
		}
		jobs, err := s.Every(c.GetFrequency()).Seconds().Do(func() {
			defer func() {
				if err := recover(); err != nil {
					p.log.Errorf("syncToDB panic:%+v", err)
				}
			}()
			if err := p.syncToDB(name, c); err != nil {
				time.Sleep(time.Second * 5)
			}
		})
		if err != nil {
			p.log.Fatalf("init jobs error:%+v", err)
		}
		jobs.SingletonMode()
	}
	if syncEnable {
		s.StartAsync()
	}
}

func (p *PacketPool) syncToDB(name string, c chains.BlockChain) error {
	cn := name
	var delayBlock uint64 = 1
	syncState := model.SyncState{ChainName: cn}
	if err := p.DB.Model(&syncState).Find(&syncState).Error; err != nil {
		return fmt.Errorf("get syncState error:%s", err.Error())
	}
	systemGauge := p.MetricsManager.Gauge.With("chain_name", cn)
	chainHeight, err := c.GetLatestHeight()
	if err != nil {
		systemGauge.With("option", "interrupt").Add(1)
		return fmt.Errorf("GetLatestHeight error:%s", err.Error())
	} else {
		systemGauge.With("option", "interrupt").Set(0)
	}
	p.log.Infof("chainName:%v,latestHeight:%v", cn, chainHeight)
	if syncState.Height+delayBlock > chainHeight {
		return fmt.Errorf("sync height :%v > chain %v Height:%v", syncState.Height, cn, chainHeight)
	}
	p.log.Infof("sync chain %v data,height=%v", cn, syncState.Height)
	var updateHeight uint64
	//delay 1 block
	if chainHeight-c.GetBatchNumber()-delayBlock > syncState.Height {
		updateHeight = syncState.Height + c.GetBatchNumber()
		err = p.saveCrossChainPacketsByHeight(syncState.Height, syncState.Height+c.GetBatchNumber()-1, c, updateHeight)
	} else {
		updateHeight = chainHeight - delayBlock + 1
		err = p.saveCrossChainPacketsByHeight(syncState.Height, chainHeight-delayBlock, c, updateHeight)
	}
	if err != nil {
		p.log.Errorf("saveCrossChainPacketsByHeight error:%s", err)
		return err
	}
	return nil
}

func (p *PacketPool) saveCrossChainPacketsByHeight(fromBlock, toBlock uint64, chain chains.BlockChain, updateHeight uint64) error {
	packets, err := chain.GetPackets(fromBlock, toBlock)
	if err != nil {
		return err
	}
	var (
		crossChainTxs    []model.CrossChainTransaction
		ackcrossChainTxs []model.CrossChainTransaction
	)

	nativeDecimal, err := chain.GetNativeDecimal()
	if err != nil {
		return err
	}
	for _, pkts := range packets {
		for _, pt := range pkts.BizPackets {
			if pt.Packet.TransferData == nil {
				continue
			}
			var transferData packettypes.TransferData
			if pt.Packet.TransferData == nil {
				if err := transferData.ABIDecode(pt.Packet.TransferData); err != nil {
					return err
				}
			}else {
				return nil
			}
			a := big.Int{}
			amount := a.SetBytes(transferData.Amount)
			sender := pt.Packet.Sender
			if pt.Signer != common.BytesToAddress([]byte{0x00}) {
				sender = pt.Signer.Hex()
			}
			crossChainTx := model.CrossChainTransaction{
				SrcChain: pt.Packet.SrcChain,
				//RelayChain:       pt.Packet.RelayChain,
				DestChain:        pt.Packet.DstChain,
				Sequence:         pt.Packet.Sequence,
				Sender:           sender,
				Receiver:         transferData.Receiver,
				SendTokenAddress: transferData.Token,
				Status:           int8(model.Pending),
				SendTxHash:       pt.TxHash,
				SendTxTime:       pt.TimeStamp,
				SrcHeight:        pt.Height,
				AmountRaw:        a.String(),
			}
			var crossChainTransaction model.CrossChainTransaction
			if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", pt.Packet.SrcChain, pt.Packet.DstChain, pt.Packet.Sequence).Find(&crossChainTransaction).Error; err != nil {
				return err
			}
			var tokenName string
			if crossChainTransaction.TokenName == "" {
				tokenName = p.ReconciliationCli.GetTokenNameByAddress(pt.Packet.SrcChain, transferData.Token)
				if tokenName == "" {
					p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", pt.Packet.SrcChain, transferData.Token, pt.Packet.DstChain, pt.Packet.Sequence)
					continue
				}
				crossChainTx.TokenName = tokenName
			} else {
				tokenName = crossChainTransaction.TokenName
			}
			if pt.Packet.SrcChain == chains.TeleportChain {
				key := fmt.Sprintf("%v/%v/%v", pt.Packet.SrcChain, pt.Packet.DstChain, tokenName)
				crossChainTx.ReceiveTokenAddress = p.ReconciliationCli.BridgeTokenMap[key].ChainBToken.AddressHex()
			} else {
				key := fmt.Sprintf("%v/%v/%v", pt.Packet.DstChain, pt.Packet.SrcChain, tokenName)
				crossChainTx.ReceiveTokenAddress = p.ReconciliationCli.BridgeTokenMap[key].ChainAToken.AddressHex()
			}
			srcChain := chain
			packetFee, err := srcChain.GetPacketFee(pt.Packet.SrcChain, pt.Packet.DstChain, int(pt.Packet.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infoln("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.ReconciliationCli.GetSingleTokenDecimals(srcChain.ChainName(), packetFee.TokenAddress)
				if err != nil {
					p.log.Errorf("GetSingleTokenDecimals error:%+v\n,chainName:%v,tokenAddr:%v", err, srcChain.ChainName(), packetFee.TokenAddress)
					return err
				}
				feeBigFloat := new(big.Float).SetInt(packetFee.FeeAmount)
				div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
				afterDecimals := new(big.Float).Quo(feeBigFloat, new(big.Float).SetInt(div)) // evm base token decimals 1e18
				fee, _ := afterDecimals.Float64()
				crossChainTx.PacketFeePaid = fee
			}
			if crossChainTransaction.AmountFloat == 0 {
				teleportDecimal, otherDecimal, err := p.ReconciliationCli.GetBridgeTokenDecimals(pt.Packet.SrcChain, tokenName, pt.Packet.DstChain)
				if err != nil {
					p.log.Errorf("ReconciliationCli.GetBridgeTokenDecimals failed:%+v\n,packet:%v,packet type:%v,txHash:%v", err, pt.Packet, Packet, pt.TxHash)
					return err
				}
				var amountFloat float64
				if pt.Packet.DstChain == chains.TeleportChain {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
					amountFloat = amountFloat * math.Pow10(int(teleportDecimal)-int(otherDecimal))
				} else {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
				}
				if err != nil {
					p.log.Errorf("tokenAmount.Float64 failed:%+v\n,packet:%v", err, pt.Packet)
					return err
				}
				amountStr := fmt.Sprintf("%.0f", amountFloat)
				amountFloat, _ = strconv.ParseFloat(amountStr, 64)
				crossChainTx.Amount = amountStr
				crossChainTx.AmountFloat = amountFloat
			}
			if pt.Packet.Sender == chains.AgentContract {
				crossChainTx.MultiID2 = pt.MultiId
			}
			crossChainTxs = append(crossChainTxs, crossChainTx)
		}
		for _, ackPacket := range pkts.AckPackets {
			if ackPacket.Ack.Packet.TransferData == nil {
				continue
			}
			var transferData packettypes.TransferData
			if err := transferData.ABIDecode(ackPacket.Ack.Packet.TransferData); err != nil {
				return err
			}
			var ack packettypes.Acknowledgement
			if err := ack.ABIDecode(ackPacket.Ack.Acknowledgement); err != nil {
				return err
			}
			a := big.Int{}
			amount := a.SetBytes(transferData.Amount)
			status, msg := getStatus(ack)
			crossChainTx := model.CrossChainTransaction{
				SrcChain: ackPacket.Ack.Packet.SrcChain,
				//RelayChain:       ackPacket.Ack.Packet.RelayChain,
				DestChain:        ackPacket.Ack.Packet.DstChain,
				Sequence:         ackPacket.Ack.Packet.Sequence,
				Receiver:         transferData.Receiver,
				Status:           int8(status),
				ErrMessage:       msg,
				ReceiveTxHash:    ackPacket.TxHash,
				SendTokenAddress: transferData.Token,
				ReceiveTxTime:    ackPacket.TimeStamp,
				DestHeight:       ackPacket.Height,
				// An Ack is generated where a packet is received
				PacketGas:      float64(ackPacket.Gas),
				PacketGasPrice: ackPacket.GasPrice,
				PacketFee:      float64(ackPacket.Gas) * ackPacket.GasPrice / math.Pow10(int(nativeDecimal)),
			}
			var crossChainTransaction model.CrossChainTransaction
			if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", ackPacket.Ack.Packet.SrcChain, ackPacket.Ack.Packet.DstChain, ackPacket.Ack.Packet.Sequence).Find(&crossChainTransaction).Error; err != nil {
				p.log.Errorf("query db error:%+v\n,where src_chain = %v and dest_chain = %v and sequence = %v", err, ackPacket.Ack.Packet.SrcChain, ackPacket.Ack.Packet.DstChain, ackPacket.Ack.Packet.Sequence)
				return err
			}
			var tokenName string
			if crossChainTransaction.TokenName == "" {
				tokenName = p.ReconciliationCli.GetTokenNameByAddress(ackPacket.Ack.Packet.SrcChain, transferData.Token)
				if tokenName == "" {
					p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", ackPacket.Ack.Packet.SrcChain, transferData.Token, ackPacket.Ack.Packet.DstChain, ackPacket.Ack.Packet.Sequence)
					continue
				}
				crossChainTx.TokenName = tokenName
			} else {
				tokenName = crossChainTransaction.TokenName
			}

			//tokenLimit, err := chain.GetTokenLimit(common.HexToAddress(ft.Token), big.NewInt(int64(ackPacket.Height)))
			//if err != nil {
			//	return err
			//}
			//if tokenLimit.Enable {
			//	totalCrossAmount, err := p.ReconciliationCli.GetDestReceivedSumAmount(ackPacket.Ack.Packet.SrcChain, ackPacket.Ack.Packet.DstChain, ft.Token, time.Duration(tokenLimit.TimePeriod.Uint64()))
			//	if err != nil {
			//		return err
			//	}
			//	limit := new(big.Float).SetInt(tokenLimit.TimeBasedLimit)
			//	limitFloat, _ := limit.Float64()
			//	if totalCrossAmount > limitFloat {
			//		p.MetricsManager.Gauge.With("chain_name", chain.ChainName()).With("option", "over_limit").Set(totalCrossAmount - limitFloat)
			//	}
			//}
			srcChain := chain
			if srcChain == nil {
				return fmt.Errorf("invalid chain,chainName:%s", ackPacket.Ack.Packet.SrcChain)
			}
			packetFee, err := srcChain.GetPacketFee(ackPacket.Ack.Packet.SrcChain, ackPacket.Ack.Packet.DstChain, int(ackPacket.Ack.Packet.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infoln("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.ReconciliationCli.GetSingleTokenDecimals(ackPacket.Ack.Packet.SrcChain, packetFee.TokenAddress)
				if err != nil {
					p.log.Errorf("GetSingleTokenDecimals error:%+v\n,chainName:%v,tokenAddr:%v", err, srcChain.ChainName(), packetFee.TokenAddress)
					return err
				}
				feeBigFloat := new(big.Float).SetInt(packetFee.FeeAmount)
				div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
				afterDecimals := new(big.Float).Quo(feeBigFloat, new(big.Float).SetInt(div)) // evm base token decimals 1e18
				fee, _ := afterDecimals.Float64()
				crossChainTx.PacketFeePaid = fee
			}
			if crossChainTransaction.AmountFloat == 0 {
				teleportDecimal, otherDecimal, err := p.ReconciliationCli.GetBridgeTokenDecimals(ackPacket.Ack.Packet.SrcChain, tokenName, ackPacket.Ack.Packet.DstChain)
				if err != nil {
					p.log.Errorf("ReconciliationCli.GetBridgeTokenDecimals failed:%+v\n,packet:%v,packet type:%v,txHash:%v", err, ackPacket.Ack.Packet, Ack, ackPacket.TxHash)
					return err
				}
				var amountFloat float64
				if ackPacket.Ack.Packet.DstChain == chains.TeleportChain {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
					amountFloat = amountFloat * math.Pow10(int(teleportDecimal)-int(otherDecimal))
				} else {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
				}
				if err != nil {
					p.log.Errorf("tokenAmount.Float64 failed:%+v\n,packet:%v", err, ackPacket.Ack.Packet)
					return err
				}
				amountStr := fmt.Sprintf("%.0f", amountFloat)
				amountFloat, _ = strconv.ParseFloat(amountStr, 64)
				crossChainTx.Amount = amountStr
				crossChainTx.AmountFloat = amountFloat
			}
			var callData  packettypes.CallData
			if ackPacket.Ack.Packet.CallData == nil {
				if err := callData.ABIDecode(ackPacket.Ack.Packet.CallData);err != nil {
					return err
				}
				if callData.ContractAddress == chains.AgentContract {
					crossChainTx.MultiID1 = ackPacket.MultiId
				}
			}
			crossChainTxs = append(crossChainTxs, crossChainTx)
			ackcrossChainTxs = append(ackcrossChainTxs, crossChainTx)
		}
		// received ack
		for _, receivedAckPacket := range pkts.RecivedAcks {
			if receivedAckPacket.Ack.Packet.TransferData == nil {
				continue
			}
			var transferData packettypes.TransferData
			if err := transferData.ABIDecode(receivedAckPacket.Ack.Packet.TransferData); err != nil {
				return err
			}
			var ack packettypes.Acknowledgement
			if err := ack.ABIDecode(receivedAckPacket.Ack.Acknowledgement); err != nil {
				return err
			}
			status, msg := getStatus(ack)
			if status != model.Fail {
				continue
			}
			crossChainTx := model.CrossChainTransaction{
				SrcChain: receivedAckPacket.Ack.Packet.SrcChain,
				//RelayChain:   receivedAckPacket.Ack.Packet.RelayChain,
				DestChain:    receivedAckPacket.Ack.Packet.DstChain,
				Sequence:     receivedAckPacket.Ack.Packet.Sequence,
				Status:       int8(model.Refund),
				ErrMessage:   msg,
				RefundTxHash: receivedAckPacket.TxHash,
				RefundTxTime: receivedAckPacket.TimeStamp,
				RefundHeight: receivedAckPacket.Height,
				// A ReceivedAck is generated at the place where the Ack is received
				AckGas:      float64(receivedAckPacket.Gas),
				AckGasPrice: receivedAckPacket.GasPrice,
				AckFee:      float64(receivedAckPacket.Gas) * receivedAckPacket.GasPrice / math.Pow10(int(nativeDecimal)),
			}
			srcChain := p.Chains[receivedAckPacket.Ack.Packet.SrcChain]
			if srcChain == nil {
				return fmt.Errorf("invalid chain,chainName:%s", receivedAckPacket.Ack.Packet.SrcChain)
			}
			packetFee, err := srcChain.GetPacketFee(receivedAckPacket.Ack.Packet.SrcChain, receivedAckPacket.Ack.Packet.DstChain, int(receivedAckPacket.Ack.Packet.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infoln("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.ReconciliationCli.GetSingleTokenDecimals(receivedAckPacket.Ack.Packet.SrcChain, packetFee.TokenAddress)
				if err != nil {
					p.log.Errorf("GetSingleTokenDecimals error:%+v\n,chainName:%v,tokenAddr:%v", err, srcChain.ChainName(), packetFee.TokenAddress)
					return err
				}
				feeBigFloat := new(big.Float).SetInt(packetFee.FeeAmount)
				div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
				afterDecimals := new(big.Float).Quo(feeBigFloat, new(big.Float).SetInt(div)) // evm base token decimals 1e18
				fee, _ := afterDecimals.Float64()
				crossChainTx.PacketFeePaid = fee
			}
			crossChainTxs = append(crossChainTxs, crossChainTx)
		}
	}
	if err := p.saveToDB(crossChainTxs, chain.ChainName(), updateHeight); err != nil {
		return fmt.Errorf("sync transaction fail:%v", err.Error())
	}
	return p.reconcile(ackcrossChainTxs)
}

func (p *PacketPool) HandlePacket(bizPackets []chains.PacketTx, crossChainTxs []model.CrossChainTransaction) error {
	for _, pt := range bizPackets {
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(pt.Packet.TransferData); err != nil {
			return err
		}
		a := big.Int{}
		amount := a.SetBytes(transferData.Amount)
		sender := pt.Packet.Sender
		if pt.Signer != common.BytesToAddress([]byte{0x00}) {
			sender = pt.Signer.Hex()
		}
		crossChainTx := model.CrossChainTransaction{
			SrcChain: pt.Packet.SrcChain,
			//RelayChain:       pt.Packet.RelayChain,
			DestChain:        pt.Packet.DstChain,
			Sequence:         pt.Packet.Sequence,
			Sender:           sender,
			Receiver:         transferData.Receiver,
			SendTokenAddress: transferData.Token,
			Status:           int8(model.Pending),
			SendTxHash:       pt.TxHash,
			SendTxTime:       pt.TimeStamp,
			SrcHeight:        pt.Height,
			AmountRaw:        a.String(),
		}
		var crossChainTransaction model.CrossChainTransaction
		if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", pt.Packet.SrcChain, pt.Packet.DstChain, pt.Packet.Sequence).Find(&crossChainTransaction).Error; err != nil {
			return err
		}
		var tokenName string
		if crossChainTransaction.TokenName == "" {
			tokenName = p.ReconciliationCli.GetTokenNameByAddress(pt.Packet.SrcChain, transferData.Token)
			if tokenName == "" {
				p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", pt.Packet.SrcChain, transferData.Token, pt.Packet.DstChain, pt.Packet.Sequence)
				continue
			}
			crossChainTx.TokenName = tokenName
		} else {
			tokenName = crossChainTransaction.TokenName
		}
		if pt.Packet.SrcChain == chains.TeleportChain {
			key := fmt.Sprintf("%v/%v/%v", pt.Packet.SrcChain, pt.Packet.DstChain, tokenName)
			crossChainTx.ReceiveTokenAddress = p.ReconciliationCli.BridgeTokenMap[key].ChainBToken.AddressHex()
		} else {
			key := fmt.Sprintf("%v/%v/%v", pt.Packet.DstChain, pt.Packet.SrcChain, tokenName)
			crossChainTx.ReceiveTokenAddress = p.ReconciliationCli.BridgeTokenMap[key].ChainAToken.AddressHex()
		}
		srcChain := p.Chains[pt.Packet.SrcChain]
		packetFee, err := srcChain.GetPacketFee(pt.Packet.SrcChain, pt.Packet.DstChain, int(pt.Packet.Sequence))
		if err != nil {
			p.log.Errorf("GetPacketFee error:%+v", err)
			return err
		} else {
			p.log.Infoln("GetPacketFee:%v", packetFee)
		}
		if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
			decimals, err := p.ReconciliationCli.GetSingleTokenDecimals(srcChain.ChainName(), packetFee.TokenAddress)
			if err != nil {
				p.log.Errorf("GetSingleTokenDecimals error:%+v\n,chainName:%v,tokenAddr:%v", err, srcChain.ChainName(), packetFee.TokenAddress)
				return err
			}
			feeBigFloat := new(big.Float).SetInt(packetFee.FeeAmount)
			div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
			afterDecimals := new(big.Float).Quo(feeBigFloat, new(big.Float).SetInt(div)) // evm base token decimals 1e18
			fee, _ := afterDecimals.Float64()
			crossChainTx.PacketFeePaid = fee
		}
		if crossChainTransaction.AmountFloat == 0 {
			teleportDecimal, otherDecimal, err := p.ReconciliationCli.GetBridgeTokenDecimals(pt.Packet.SrcChain, tokenName, pt.Packet.DstChain)
			if err != nil {
				p.log.Errorf("ReconciliationCli.GetBridgeTokenDecimals failed:%+v\n,packet:%v,packet type:%v,txHash:%v", err, pt.Packet, Packet, pt.TxHash)
				return err
			}
			var amountFloat float64
			if pt.Packet.DstChain == chains.TeleportChain {
				tokenAmount := &chains.TokenAmount{Amount: amount}
				amountFloat, err = tokenAmount.Float64()
				amountFloat = amountFloat * math.Pow10(int(teleportDecimal)-int(otherDecimal))
			} else {
				tokenAmount := &chains.TokenAmount{Amount: amount}
				amountFloat, err = tokenAmount.Float64()
			}
			if err != nil {
				p.log.Errorf("tokenAmount.Float64 failed:%+v\n,packet:%v", err, pt.Packet)
				return err
			}
			amountStr := fmt.Sprintf("%.0f", amountFloat)
			amountFloat, _ = strconv.ParseFloat(amountStr, 64)
			crossChainTx.Amount = amountStr
			crossChainTx.AmountFloat = amountFloat
		}
		if pt.Packet.Sender == chains.AgentContract {
			crossChainTx.MultiID2 = pt.MultiId
		}
		crossChainTxs = append(crossChainTxs, crossChainTx)
	}
	return nil
}

func (p *PacketPool) syncPacketByHash(chainName, txHash string) ([]model.CrossChainTransaction, error) {
	var (
		crossChainTxs []model.CrossChainTransaction
	)
	chain := p.Chains[chainName]
	if chain == nil {
		return nil, fmt.Errorf("invalid chainName")
	}
	// teleport is not support
	bizPackets, err := chain.GetPacketsByHash(txHash)
	if err != nil {
		return nil, err
	}
	if err := p.HandlePacket(bizPackets, crossChainTxs); err != nil {
		return nil, err
	}
	if err := p.DB.Transaction(func(tx *gorm.DB) error {
		for _, cctx := range crossChainTxs {
			var ccTransaction model.CrossChainTransaction
			if err := tx.Where("src_chain = ? and dest_chain = ? and sequence = ? ", cctx.SrcChain, cctx.DestChain, cctx.Sequence).First(&ccTransaction).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if cctx.MultiID2 != "" {
						continue
					}
					if err := tx.Create(&cctx).Error; err != nil {
						p.log.Errorf("create CrossChainTransaction  error!!!,result:%+v", cctx)
						return err
					}
					p.log.Infof("CrossChainTransaction create success!!!,result:%+v", cctx)
					continue
				} else {
					return err
				}
			}
			if ccTransaction.ID != 0 {
				cctx.ID = ccTransaction.ID
				if cctx.Status < ccTransaction.Status {
					cctx.Status = ccTransaction.Status
				}
				if err := tx.Updates(&cctx).Error; err != nil {
					p.log.Errorf("Update CrossChainTransaction  error!!!,result:%+v", cctx)
					return err
				}
				p.log.Infof("Update CrossChainTransaction  success!!!,result:%+v", cctx)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return crossChainTxs, nil
}

func (p *PacketPool) saveToDB(crossChainTxs []model.CrossChainTransaction, name string, updateHeight uint64) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		for _, cctx := range crossChainTxs {
			var ccTransaction model.CrossChainTransaction
			if err := tx.Where("src_chain = ? and dest_chain = ? and sequence = ? ", cctx.SrcChain, cctx.DestChain, cctx.Sequence).First(&ccTransaction).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&cctx).Error; err != nil {
						p.log.Errorf("create CrossChainTransaction  error!!!,result:%+v", cctx)
						return err
					}
					p.log.Infof("CrossChainTransaction create success!!!,result:%+v", cctx)
					continue
				} else {
					p.log.Errorf("tx.Where error:%+v", err)
					return err
				}
			}
			if ccTransaction.ID != 0 {
				cctx.ID = ccTransaction.ID
				if cctx.Status < ccTransaction.Status {
					cctx.Status = ccTransaction.Status
				}
				if err := tx.Updates(&cctx).Error; err != nil {
					p.log.Errorf("Update CrossChainTransaction  error!!!,result:%+v", cctx)
					return err
				}
				p.log.Infof("Update CrossChainTransaction  success!!!,result:%+v", cctx)
			}
		}
		if err := tx.Updates(&model.SyncState{ChainName: name, Height: updateHeight}).Error; err != nil {
			return err
		}
		p.log.Infoln("tx.Updates success,sync state:", model.SyncState{ChainName: name, Height: updateHeight})
		return nil
	})
}

func (p *PacketPool) reconcile(ackcrossChainTxs []model.CrossChainTransaction) error {
	if p.ReconcileEnable {
		for _, ackCrossChainTx := range ackcrossChainTxs {
			var crossChainTx model.CrossChainTransaction
			if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", ackCrossChainTx.SrcChain, ackCrossChainTx.DestChain, ackCrossChainTx.Sequence).Find(&crossChainTx).Error; err != nil {
				return err
			}
			packetA := bridges.ReconciliationPacket{
				SrcChain:   crossChainTx.SrcChain,
				DestChain:  crossChainTx.DestChain,
				SrcHeight:  crossChainTx.SrcHeight,
				DestHeight: crossChainTx.DestHeight,
				SrcAddress: crossChainTx.SendTokenAddress,
			}
			packetB := bridges.ReconciliationPacket{
				SrcChain:   crossChainTx.DestChain,
				DestChain:  crossChainTx.SrcChain,
				SrcHeight:  crossChainTx.DestHeight,
				DestHeight: crossChainTx.SrcHeight,
			}
			var teleportPacket, otherPacket bridges.ReconciliationPacket
			if crossChainTx.SrcChain == "teleport" {
				teleportPacket = packetA
				otherPacket = packetB
			} else {
				teleportPacket = packetB
				otherPacket = packetA
			}
			if crossChainTx.SrcHeight != 0 && crossChainTx.DestHeight != 0 {
				bridgeReconcileResult, err := p.ReconciliationCli.Reconcile(teleportPacket, otherPacket)
				if err != nil {
					return err
				}
				crossChainTransaction := model.CrossChainTransaction{
					TokenName:       bridgeReconcileResult.TokenName,
					SrcAmount:       bridgeReconcileResult.SrcAmount,
					SrcTrimAmount:   bridgeReconcileResult.SrcTrimAmount,
					DestAmount:      bridgeReconcileResult.DestAmount,
					DestTrimAmount:  bridgeReconcileResult.DestTrimAmount,
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
				if err := p.DB.Updates(&crossChainTransaction).Error; err != nil {
					p.log.Errorf("update Reconcile result error:%+v", err)
					return err
				}
			}
		}
	}
	return nil
}


func getStatus(ack packettypes.Acknowledgement) (model.PacketStatus, string) {
	var status model.PacketStatus
	var ackMsg string
	if ack.Code == 0 {
		status = model.Success
	} else {
		status = model.Fail
		ackMsg = ack.GetMessage()
	}
	return status, ackMsg
}

func makeCodec() *codec.ProtoCodec {
	ir := codectypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(ir)
	govtypes.RegisterInterfaces(ir)
	xibctendermint.RegisterInterfaces(ir)
	xibceth.RegisterInterfaces(ir)
	packettypes.RegisterInterfaces(ir)
	ir.RegisterInterface("cosmos.v1beta1.Msg", (*sdk.Msg)(nil))
	tx.RegisterInterfaces(ir)
	cryptocodec.RegisterInterfaces(ir)
	return codec.NewProtoCodec(ir)
}
