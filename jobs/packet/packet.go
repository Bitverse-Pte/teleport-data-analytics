package packet

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/teleport-network/teleport-data-analytics/jobs/datas"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/jobs/bridges"
	"github.com/teleport-network/teleport-data-analytics/metrics"
	"github.com/teleport-network/teleport-data-analytics/model"

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
	Codec           *codec.ProtoCodec
	DB              *gorm.DB
	log             *logrus.Logger
	Chains          map[string]chains.BlockChain
	ChainMap        map[string]string
	IDToName        map[string]string
	NameToID        map[string]string
	ReconcileEnable bool
	BridegesManager *bridges.Bridges
	MetricsManager  *metrics.MetricManager
}

const (
	Packet = "packet"
	Ack    = "ack"
)

func NewPacketDBPool(db *gorm.DB, log *logrus.Logger, chain map[string]chains.BlockChain, chainMap map[string]string, bridgesManager *bridges.Bridges, reconcileEnable bool, metricsManager *metrics.MetricManager) *PacketPool {
	cdc := makeCodec()
	if err := db.AutoMigrate(&model.SyncState{}, &model.CrossPacket{}, &model.PacketRelation{}, &model.BridgeReconcileResult{}, &model.Record{}, &model.BlockReconcileResult{}, &model.CrossChainTransaction{}, &model.GlobalMetrics{}, &model.SingleDirectionBridgeMetrics{}, &model.Token{}); err != nil {
		panic(fmt.Errorf("db.AutoMigrate:%+v", err))
	}
	return &PacketPool{
		Codec:           cdc,
		DB:              db,
		log:             log,
		Chains:          chain,
		ChainMap:        chainMap,
		ReconcileEnable: reconcileEnable,
		BridegesManager: bridgesManager,
		MetricsManager:  metricsManager,
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
		p.log.Errorf("saveCrossChainPacketsByHeight error:%+v", err)
		return err
	}
	return nil
}

func (p *PacketPool) saveCrossChainPacketsByHeight(fromBlock, toBlock uint64, chain chains.BlockChain, updateHeight uint64) error {
	packets, err := chain.GetPackets(fromBlock, toBlock)
	if err != nil {
		return fmt.Errorf("get packets error:%s", err.Error())
	}
	var (
		crossChainTxs    []model.CrossChainTransaction
		ackcrossChainTxs []model.CrossChainTransaction
	)

	nativeDecimal, err := chain.GetNativeDecimal()
	if err != nil {
		return fmt.Errorf("get nativeDecimal error:%s", err.Error())
	}
	for _, pkts := range packets {
		for _, pt := range pkts.Packets {
			sender := pt.Sender
			if pt.Signer != common.BytesToAddress([]byte{0x00}) {
				sender = pt.Signer.Hex()
			}
			crossChainTx := model.CrossChainTransaction{
				SrcChain:         pt.SrcChain,
				DestChain:        pt.DstChain,
				Sequence:         pt.Sequence,
				Sender:           sender,
				Receiver:         pt.Receiver,
				SendTokenAddress: pt.Token,
				Status:           int8(model.Pending),
				SendTxHash:       pt.TxHash,
				SendTxTime:       pt.TimeStamp,
				SrcHeight:        pt.Height,
				AmountRaw:        pt.Amount,
				MultiID2:         pt.MultiId,
				SrcChainId:       pt.SrcChainId,
				DestChainId:      pt.DestChainId,
			}
			var crossChainTransaction model.CrossChainTransaction
			if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", pt.SrcChain, pt.DstChain, pt.Sequence).Find(&crossChainTransaction).Error; err != nil {
				return err
			}
			var tokenName string
			if crossChainTransaction.TokenName == "" {
				tokenName = p.BridegesManager.GetTokenNameByAddress(pt.SrcChain, pt.Token)
				if tokenName == "" {
					p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", pt.SrcChain, pt.Token, pt.DstChain, pt.Sequence)
					continue
				}
				crossChainTx.TokenName = tokenName
			} else {
				tokenName = crossChainTransaction.TokenName
			}
			if pt.SrcChain == chains.TeleportChain {
				key := fmt.Sprintf("%v/%v/%v", pt.SrcChain, pt.DstChain, tokenName)
				crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainBToken.Address()
			} else {
				key := fmt.Sprintf("%v/%v/%v", pt.DstChain, pt.SrcChain, tokenName)
				crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainAToken.Address()
			}
			srcChain := chain
			packetFee, err := srcChain.GetPacketFee(pt.SrcChain, pt.DstChain, int(pt.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infof("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.BridegesManager.GetSingleTokenDecimals(srcChain.ChainName(), packetFee.TokenAddress)
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
				teleportDecimal, otherDecimal, err := p.BridegesManager.GetBridgeTokenDecimals(pt.SrcChain, tokenName, pt.DstChain)
				if err != nil {
					p.log.Errorf("BridegesManager.GetBridgeTokenDecimals failed:%+v\n,packet:%v,packet type:%v,txHash:%v", err, pt, Packet, pt.TxHash)
					return err
				}
				var amountFloat float64
				amount := new(big.Int)
				amount.SetString(pt.Amount, 10)
				if pt.DstChain == chains.TeleportChain {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
					amountFloat = amountFloat * math.Pow10(int(teleportDecimal)-int(otherDecimal))
				} else {
					tokenAmount := &chains.TokenAmount{Amount: amount}
					amountFloat, err = tokenAmount.Float64()
				}
				if err != nil {
					p.log.Errorf("tokenAmount.Float64 failed:%+v\n,packet:%v", err, pt)
					return err
				}
				amountStr := fmt.Sprintf("%.0f", amountFloat)
				amountFloat, _ = strconv.ParseFloat(amountStr, 64)
				crossChainTx.Amount = amountStr
				crossChainTx.AmountFloat = amountFloat
			}
			if pt.Sender == chains.AgentContract {
				crossChainTx.MultiID2 = pt.MultiId
			}
			crossChainTxs = append(crossChainTxs, crossChainTx)
		}
		for _, ackPacket := range pkts.AckPackets {
			crossChainTx := model.CrossChainTransaction{
				SrcChain: ackPacket.SrcChain,
				DestChain:        ackPacket.DstChain,
				Sequence:         ackPacket.Sequence,
				Status:           ackPacket.Code,
				ErrMessage:       ackPacket.ErrMsg,
				ReceiveTxHash:    ackPacket.TxHash,
				SendTokenAddress: ackPacket.Token,
				ReceiveTxTime:    ackPacket.TimeStamp,
				DestHeight:       ackPacket.Height,
				// An Ack is generated where a packet is received
				PacketGas:      float64(ackPacket.Gas),
				PacketGasPrice: ackPacket.GasPrice,
				PacketFee:      float64(ackPacket.Gas) * ackPacket.GasPrice / math.Pow10(int(nativeDecimal)),
				MultiID1:       ackPacket.MultiId,
			}
			var crossChainTransaction model.CrossChainTransaction
			if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", ackPacket.SrcChain, ackPacket.DstChain, ackPacket.Sequence).Find(&crossChainTransaction).Error; err != nil {
				p.log.Errorf("query db error:%+v\n,where src_chain = %v and dest_chain = %v and sequence = %v", err, ackPacket.SrcChain, ackPacket.DstChain, ackPacket.Sequence)
				return err
			}
			var tokenName string
			if crossChainTransaction.TokenName == "" {
				tokenName = p.BridegesManager.GetTokenNameByAddress(ackPacket.SrcChain, ackPacket.Token)
				if tokenName == "" {
					p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", ackPacket.SrcChain, ackPacket.Token, ackPacket.DstChain, ackPacket.Sequence)
					continue
				}
				crossChainTx.TokenName = tokenName
			} else {
				tokenName = crossChainTransaction.TokenName
			}
			srcChain := chain
			if srcChain == nil {
				return fmt.Errorf("invalid chain,chainName:%s", ackPacket.SrcChain)
			}
			packetFee, err := srcChain.GetPacketFee(ackPacket.SrcChain, ackPacket.DstChain, int(ackPacket.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infof("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.BridegesManager.GetSingleTokenDecimals(ackPacket.SrcChain, packetFee.TokenAddress)
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
			ackcrossChainTxs = append(ackcrossChainTxs, crossChainTx)
		}
		// received ack
		for _, receivedAckPacket := range pkts.RecivedAcks {
			crossChainTx := model.CrossChainTransaction{
				SrcChain:     receivedAckPacket.SrcChain,
				DestChain:    receivedAckPacket.DstChain,
				Sequence:     receivedAckPacket.Sequence,
				Status:       int8(model.Refund),
				ErrMessage:   receivedAckPacket.ErrMsg,
				RefundTxHash: receivedAckPacket.TxHash,
				RefundTxTime: receivedAckPacket.TimeStamp,
				RefundHeight: receivedAckPacket.Height,
				// A ReceivedAck is generated at the place where the Ack is received
				AckGas:      float64(receivedAckPacket.Gas),
				AckGasPrice: receivedAckPacket.GasPrice,
				AckFee:      float64(receivedAckPacket.Gas) * receivedAckPacket.GasPrice / math.Pow10(int(nativeDecimal)),
			}
			srcChain := p.Chains[receivedAckPacket.SrcChain]
			if srcChain == nil {
				return fmt.Errorf("invalid chain,chainName:%s", receivedAckPacket.SrcChain)
			}
			packetFee, err := srcChain.GetPacketFee(receivedAckPacket.SrcChain, receivedAckPacket.DstChain, int(receivedAckPacket.Sequence))
			if err != nil {
				p.log.Errorf("GetPacketFee error:%+v", err)
				return err
			} else {
				p.log.Infof("GetPacketFee:%v", packetFee)
			}
			if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
				decimals, err := p.BridegesManager.GetSingleTokenDecimals(receivedAckPacket.SrcChain, packetFee.TokenAddress)
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
	// judge if need to init singleDirectionMeritsTable

	// classify crossChainTxs by different bridge
	//classifiedTxs := p.classifyCrossChainTxs(chain.ChainName(), crossChainTxs)
	executed := false
	if executed, err = p.InitSingleDirectionBridgeMetrics(); err != nil {
		return err
	}

	// update metrics table by classifiedTxs
	if !executed {
		if err = p.UpdateMetrics(crossChainTxs); err != nil {
			return fmt.Errorf("update metrics fail:%v", err.Error())
		}
	}

	p.log.Infof("update metrics success,chainName:%s,updateHeight:%d", chain.ChainName(), updateHeight)
	return p.reconcile(ackcrossChainTxs)
}

func (p *PacketPool) HandlePacket(bizPackets []chains.BasePacketTx, crossChainTxs []model.CrossChainTransaction) error {
	for _, pt := range bizPackets {
		a := new(big.Int)
		amount, ok := a.SetString(pt.Amount, 10)
		if !ok {
			return fmt.Errorf("invalid amount")
		}
		sender := pt.Sender
		if pt.Signer != common.BytesToAddress([]byte{0x00}) {
			sender = pt.Signer.Hex()
		}
		crossChainTx := model.CrossChainTransaction{
			SrcChain: pt.SrcChain,
			//RelayChain:       pt.RelayChain,
			DestChain:        pt.DstChain,
			Sequence:         pt.Sequence,
			Sender:           sender,
			Receiver:         pt.Receiver,
			SendTokenAddress: pt.Token,
			Status:           int8(model.Pending),
			SendTxHash:       pt.TxHash,
			SendTxTime:       pt.TimeStamp,
			SrcHeight:        pt.Height,
			AmountRaw:        a.String(),
		}
		var crossChainTransaction model.CrossChainTransaction
		if err := p.DB.Where("src_chain = ? and dest_chain = ? and sequence = ?", pt.SrcChain, pt.DstChain, pt.Sequence).Find(&crossChainTransaction).Error; err != nil {
			return err
		}
		var tokenName string
		if crossChainTransaction.TokenName == "" {
			tokenName = p.BridegesManager.GetTokenNameByAddress(pt.SrcChain, pt.Token)
			if tokenName == "" {
				p.log.Errorf("skip invalid token address,chainName:%v,tokenAddress:%v\n,destChain:%v,sequence:%v", pt.SrcChain, pt.Token, pt.DstChain, pt.Sequence)
				continue
			}
			crossChainTx.TokenName = tokenName
		} else {
			tokenName = crossChainTransaction.TokenName
		}
		key := fmt.Sprintf("%v/%v/%v", pt.DstChain, pt.SrcChain, tokenName)
		if pt.SrcChain == chains.TeleportChain {
			key = fmt.Sprintf("%v/%v/%v", pt.SrcChain, pt.DstChain, tokenName)
			crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainBToken.Address()
		}
		crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainAToken.Address()
		p.log.Infof("srcChain:%v,destChain:%v,ReceiveTokenAddress:%v",pt.SrcChain,pt.DstChain,crossChainTx.ReceiveTokenAddress)
		srcChain := p.Chains[pt.SrcChain]
		packetFee, err := srcChain.GetPacketFee(pt.SrcChain, pt.DstChain, int(pt.Sequence))
		if err != nil {
			p.log.Errorf("GetPacketFee error:%+v", err)
			return err
		} else {
			p.log.Infof("GetPacketFee:%v", packetFee)
		}
		if packetFee != nil && packetFee.FeeAmount != nil && packetFee.TokenAddress != "" {
			decimals, err := p.BridegesManager.GetSingleTokenDecimals(srcChain.ChainName(), packetFee.TokenAddress)
			if err != nil {
				p.log.Errorf("GetSingleTokenDecimals error:%+v\n,chainName:%v,tokenAddr:%v", err, srcChain.ChainName(), packetFee.TokenAddress)
				return err
			}
			feeBigFloat := new(big.Float).SetInt(packetFee.FeeAmount)
			div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
			afterDecimals := new(big.Float).Quo(feeBigFloat, new(big.Float).SetInt(div)) // evm base token decimals 1e18
			fee, _ := afterDecimals.Float64()
			crossChainTx.PacketFeePaid = fee
		}else {
			p.log.Warnf("packetFee query failed:%+v",packetFee)
		}
		if crossChainTransaction.AmountFloat == 0 {
			teleportDecimal, otherDecimal, err := p.BridegesManager.GetBridgeTokenDecimals(pt.SrcChain, tokenName, pt.DstChain)
			if err != nil {
				p.log.Errorf("BridegesManager.GetBridgeTokenDecimals failed:%+v\n,packet:%v,packet type:%v,txHash:%v", err, pt, Packet, pt.TxHash)
				return err
			}
			var amountFloat float64
			if pt.DstChain == chains.TeleportChain {
				tokenAmount := &chains.TokenAmount{Amount: amount}
				amountFloat, err = tokenAmount.Float64()
				amountFloat = amountFloat * math.Pow10(int(teleportDecimal)-int(otherDecimal))
			} else {
				tokenAmount := &chains.TokenAmount{Amount: amount}
				amountFloat, err = tokenAmount.Float64()
			}
			if err != nil {
				p.log.Errorf("tokenAmount.Float64 failed:%+v\n,packet:%v", err, pt)
				return err
			}
			amountStr := fmt.Sprintf("%.0f", amountFloat)
			amountFloat, _ = strconv.ParseFloat(amountStr, 64)
			crossChainTx.Amount = amountStr
			crossChainTx.AmountFloat = amountFloat
		}
		if pt.Sender == chains.AgentContract {
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
		var (
			//	txInc         uint64
			//	txFailedCount uint64
			senderSet = mapset.NewSet[string]()
		)
		// scan all crossChainTxs
		for _, cctx := range crossChainTxs {
			// get all sender
			senderSet.Add(cctx.Sender)
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

		// fill globalMetrics
		dUserAmt := 0
		//query allsender from crossChainTxs
		for _, sender := range senderSet.ToSlice() {
			// query sender if exist from crossChainTxs
			var cctx model.CrossChainTransaction
			if err := tx.Where("sender = ?", sender).First(&cctx).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					dUserAmt++
				} else {
					p.log.Errorf("tx.Where error:%+v", err)
					return err
				}
			}
		}
		//fill globalMetrics
		glbalMetrics := &model.GlobalMetrics{
			UserAmt: int64(dUserAmt),
		}
		//if !markGlobal {
		if err := p.insertGlobalMetrics(glbalMetrics); err != nil {
			return err
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
			if !strings.Contains(ackCrossChainTx.SendTokenAddress, "0x") {
				continue
			}
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
				bridgeReconcileResult, err := p.BridegesManager.Reconcile(teleportPacket, otherPacket)
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

// init SingleDirectionBridgeMetrics table by CrossChainTransaction table
// Notice: this function must run after init CrossChainTransaction table
// Warn: You must clean SingleDirectionBridgeMetrics table after you clean the CrossChainTransaction table
// todo: packet fee token may not be the same as the token of the bridge
func (p *PacketPool) InitSingleDirectionBridgeMetrics() (executed bool, err error) {
	// query SingleDirectionBridgeMetrics if it is empty then init otherwise update it
	if err = p.DB.Last(&model.SingleDirectionBridgeMetrics{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// if SingleDirectionBridgeMetrics is empty then init it with crossChainTransaction table
			var src, dest string // src and dest chain name
			var amt string
			if err = p.DB.Transaction(func(tx *gorm.DB) error {
				// init SingleDirectionBridgeMetrics table by scan CrossChainTransaction table
				for srcChainName, targetChainList := range datas.BridgeNameAdjList {
					for _, destChainName := range targetChainList {
						// every bridge has two direction ,so we need to scan twice to construct the single direction bridge metrics
						for i := 0; i < 2; i++ {
							if i == 0 {
								src = srcChainName
								dest = destChainName
							} else {
								src = destChainName
								dest = srcChainName
							}

							st := &model.SingleDirectionBridgeMetrics{
								SrcChain:  src,
								DestChain: dest,
								//PktAmt:
							}
							// calc bridge info
							//  calc pktAmount
							// select count(*) as pktAmount from cross_chain_transactions where src_chain='bsctest' and dest_chain='teleport' ;
							if err := tx.Table("cross_chain_transactions").Where("src_chain = ? and dest_chain = ?", src, dest).Count(&st.PktAmt).Error; err != nil {
								p.log.Errorf("query pktAmount error:%+v", err)
								return err // return error if query pktAmount error
							}

							// calc FailedAmt
							//SELECT count(*) as failedAmt FROM backend.cross_chain_transactions where src_chain='teleport' and dest_chain='bsctest' and  (status=3 or status =4);
							if err := tx.Table("cross_chain_transactions").Where("src_chain = ? and dest_chain = ? and (status = ? or status = ?)", src, dest, model.Fail, model.Refund).Count(&st.FailedAmt).Error; err != nil {
								p.log.Errorf("query FailedAmount error:%+v", err)
								return err // return error if query FailedAmount error
							}

							// create single direction bridge metrics
							if err := tx.Create(st).Error; err != nil {
								p.log.Errorf("create SingleDirectionBridgeMetrics error:%+v", err)
								return err // return error if create SingleDirectionBridgeMetrics error
							}

							// init token table
							// todo: Fee Token is not necessarily in the token list

							bridgeKey := p.ChainMap[src] + "/" + p.ChainMap[dest]
							tokens := datas.Bridges[bridgeKey].Tokens
							// travel all tokens in the bridge
							for _, token := range tokens {
								// for fee
								// calc FeeAmt
								// select coalesce(sum(packet_fee),0)  as fee_Amount from cross_chain_transactions where src_chain='teleport' and dest_chain='bsctest' ;
								if err := tx.Table("cross_chain_transactions").Select("coalesce(sum(packet_fee),0) as fee_Amount").Where("src_chain = ? and dest_chain = ? and packet_fee_demand_token = ?", src, dest, token.Name).Find(&amt).Error; err != nil {
									p.log.Errorf("query feeAmount error:%+v", err)
									return err // return error if query feeAmount error
								}
								// insert
								if err = tx.Create(&model.Token{
									BridgeID:  st.ID,
									TokenName: token.Name,
									TokenType: model.PacketFee,
									Amt:       amt,
								}).Error; err != nil {
									p.log.Errorf("create Token error:%+v", err)
								}

								// for value
								//select coalesce(sum(amount),0) as usdt_amount from cross_chain_transactions where src_chain='teleport' and dest_chain='rinkeby' and token_name = ?;
								if err := tx.Table("cross_chain_transactions").Select("coalesce(sum(amount),0) as amt").Where("src_chain = ? and dest_chain = ? and token_name = ?", src, dest, token.Name).Find(&amt).Error; err != nil {
									p.log.Errorf("query %+v amt error:%+v", token.Name, err)
									return err // return error if query usdtAmount error

								}
								// insert
								if err = tx.Create(&model.Token{
									BridgeID:  st.ID,
									TokenName: token.Name,
									TokenType: model.PacketValue,
									Amt:       amt,
								}).Error; err != nil {
									p.log.Errorf("create Token error:%+v", err)
									return err
								}

							}

							// get all packet fee token
							// select distinct(packet_fee_demand_token) from cross_chain_transactions where src_chain='teleport' and dest_chain='rinkeby' ;
							if err := tx.Table("cross_chain_transactions").Select("distinct(packet_fee_demand_token)").Where("src_chain = ? and dest_chain = ?", src, dest).Find(&tokens).Error; err != nil {

							}
							// init all packet fee
						}

					}
				}
				return nil
			}); err != nil {
				p.log.Errorf("init SingleDirectionBridgeMetrics error:%+v", err)
				return
			}
			executed = true
		} else {
			p.log.Errorf("query SingleDirectionBridgeMetrics error:%+v", err)
			return
		}
	}
	return

}

// init GlobalMetrics table by SingleDirectionBridgeMetrics table
func (p *PacketPool) initGlobalMetrics() (err error) {
	// query GlobalMetrics if it is empty then init otherwise return
	globalMetres := &model.GlobalMetrics{}
	if err = p.DB.First(&model.GlobalMetrics{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// init GlobalMetrics table by scan CrossChainTransaction table
			// SELECT count(distinct sender) FROM backend.cross_chain_transactions where length(sender)>0;
			if err = p.DB.Table("cross_chain_transactions").Select("count(distinct sender) as user_amt").Where("length(sender) > ?", 0).Find(&globalMetres.UserAmt).Error; err != nil {
				if err != gorm.ErrRecordNotFound {
					p.log.Errorf("count userAmount error:%+v", err)
					return // return error if count userAmount error
				}
			}
			// insert into GlobalMetrics table
			if err = p.DB.Create(globalMetres).Error; err != nil {
				p.log.Errorf("insert GlobalMetrics error:%+v", err)
				return // return error if insert GlobalMetrics error
			}
		} else {
			p.log.Errorf("query GlobalMetrics table error:%+v", err)
			return
		}
	}
	return nil
}

// update globalMetrics table
func (p *PacketPool) insertGlobalMetrics(dData *model.GlobalMetrics) error {
	// query GlobalMetrics
	oldData := &model.GlobalMetrics{}
	if err := p.DB.Model(&model.GlobalMetrics{}).Last(oldData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// if globalMetrics is empty then init it with crossChainTransaction table
			err = p.initGlobalMetrics()
			if err != nil {
				return err
			}
			return nil
		} else {
			p.log.Errorf("query GlobalMetrics error:%+v", err)
			return err
		}
	}

	// get new GlobalMetrics
	newData := &model.GlobalMetrics{
		UserAmt: oldData.UserAmt + dData.UserAmt,
	}

	// insert new GlobalMetrics
	if err := p.DB.Create(newData).Error; err != nil {
		p.log.Errorf("update GlobalMetrics error:%+v", err)
		return err
	}
	return nil
}

// TODO: delete @Meta
//// get latest bridge metrics and add delta to it and insert it to SingleDirectionBridgeMetrics table
//func (p *PacketPool) insertSingleDirectionBridgeMetrics(bridgeMetrics *model.SingleDirectionBridgeMetrics) error {
//	// query SingleDirectionBridgeMetrics
//	oldData := &model.SingleDirectionBridgeMetrics{}
//	if err := p.DB.Where("src_chain=? and dest_chain = ?", bridgeMetrics.SrcChain, bridgeMetrics.DestChain).Last(oldData).Error; err != nil {
//		if err != gorm.ErrRecordNotFound {
//			p.log.Errorf("query SingleDirectionBridgeMetrics error:%+v", err)
//			return err
//		}
//	}
//
//	// get new SingleDirectionBridgeMetrics
//	feeSum := tools.Sum(bridgeMetrics.FeeAmt, oldData.FeeAmt)
//	if feeSum == "" {
//		p.log.Error("sum fee error")
//		return errors.New("sum fee error")
//	}
//
//	newData := &model.SingleDirectionBridgeMetrics{
//		SrcChain:  oldData.SrcChain,
//		DestChain: oldData.DestChain,
//		PktAmt:    bridgeMetrics.PktAmt + oldData.PktAmt,
//		FeeAmt:    bridgeMetrics.FeeAmt + oldData.FeeAmt,
//		FailedAmt: bridgeMetrics.FailedAmt + oldData.FailedAmt,
//		UAmt:      bridgeMetrics.UAmt + oldData.UAmt,
//		TeleAmt:   bridgeMetrics.TeleAmt + oldData.TeleAmt,
//	}
//
//	// insert new SingleDirectionBridgeMetrics
//	if err := p.DB.Create(newData).Error; err != nil {
//		p.log.Errorf("update SingleDirectionBridgeMetrics error:%+v", err)
//		return err
//	}
//	return nil
//}

// update metrics by calc bridgeTx
func (p *PacketPool) UpdateMetrics(txs []model.CrossChainTransaction) error {
	return p.DB.Transaction(func(tx *gorm.DB) error {
		// travel all txs
		for _, t := range txs {

			// query SingleDirectionBridgeMetrics id by src_chain and dest_chain
			bridgeMetrics := &model.SingleDirectionBridgeMetrics{}
			if err := tx.Where("src_chain=? and dest_chain = ?", t.SrcChain, t.DestChain).First(bridgeMetrics).Error; err != nil {
				// if it is new bridge then insert it
				if err == gorm.ErrRecordNotFound {
					bridgeMetrics.SrcChain = t.SrcChain
					bridgeMetrics.DestChain = t.DestChain
					// insert
					if err := tx.Create(bridgeMetrics).Error; err != nil {
						p.log.Errorf("insert SingleDirectionBridgeMetrics error:%+v", err)
						return err
					}
				} else {
					p.log.Errorf("query SingleDirectionBridgeMetrics error:%+v", err)
					return err
				}
			}
			// pkt amt +1
			//SQL update single_direction_bridge_metrics set pkt_amt = pkt_amt+1 where src_chain= ? and dest_chain = ? ;
			if err := tx.Model(&model.SingleDirectionBridgeMetrics{}).Where("id=?", bridgeMetrics.ID).Update("pkt_amt", gorm.Expr("pkt_amt+?", 1)).Error; err != nil {
				p.log.Errorf("update SingleDirectionBridgeMetrics error:%+v", err)
				return err
			}

			// if tx is failed then failed_amt +1
			//SQL update single_direction_bridge_metrics set failed_amt = failed_amt+1 where src_chain= ? and dest_chain = ? ;
			if model.PacketStatus(t.Status) == model.Fail || model.PacketStatus(t.Status) == model.Refund {
				if err := p.DB.Model(&model.SingleDirectionBridgeMetrics{}).Where("id= ?", bridgeMetrics.ID).Update("failed_amt", gorm.Expr("failed_amt+?", 1)).Error; err != nil {
					p.log.Errorf("update SingleDirectionBridgeMetrics error:%+v", err)
					return err
				}
			}

			token := &model.Token{}
			// if not exist the token then insert
			if err := tx.Where("bridge_id=? and token_name=?", bridgeMetrics.ID, t.TokenName).First(token).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					token = &model.Token{
						BridgeID:  bridgeMetrics.ID,
						TokenName: t.TokenName,
						TokenType: model.PacketFee,
					}
					// insert
					if err := tx.Create(token).Error; err != nil {
						p.log.Errorf("insert Token error:%+v", err)
						return err
					}
					token = &model.Token{
						TokenName: t.TokenName,
						BridgeID:  bridgeMetrics.ID,
						TokenType: model.PacketValue,
					}
					// insert
					if err := tx.Create(token).Error; err != nil {
						p.log.Errorf("insert Token error:%+v", err)
						return err
					}

				} else {
					p.log.Errorf("query Token error:%+v", err)
					return err
				}
			}

			// update token table
			// update token table of fee amount
			//SQL update token set amt = amt+? where bridge_id = (select id from single_direction_bridge_metrics where src_chain= ? and dest_chain = ? ) and token_type = ?  and token_name=?;
			if err := tx.Model(&model.Token{}).Where("bridge_id = ? and token_type = ? and token_name=?", bridgeMetrics.ID, model.PacketFee, t.TokenName).Update("amt", gorm.Expr("amt+?", t.PacketFee)).Error; err != nil {
				p.log.Errorf("update Token error:%+v", err)
				return err
			}

			// update token table of value amount
			//SQL update token set amt = amt+? where bridge_id = (select id from single_direction_bridge_metrics where src_chain= ? and dest_chain = ? ) and token_type = ?  and token_name=?;
			if err := tx.Model(&model.Token{}).Where("bridge_id = ? and token_type = ? and token_name=?", bridgeMetrics.ID, model.PacketValue, t.TokenName).Update("amt", gorm.Expr("amt+?", t.Amount)).Error; err != nil {
				p.log.Errorf("update Token error:%+v", err)
				return err
			}

		}
		return nil
	})

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
