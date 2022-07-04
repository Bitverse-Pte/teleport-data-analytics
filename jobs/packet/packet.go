package packet

import (
	"fmt"
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
	if err := db.AutoMigrate(&model.SyncState{}, &model.CrossPacket{}, &model.PacketRelation{}, &model.BridgeReconcileResult{}, &model.Record{}, &model.BlockReconcileResult{}, &model.CrossChainTransaction{}); err != nil {
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
				p.log.Infoln("GetPacketFee:%v", packetFee)
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
				//RelayChain:       ackPacket.RelayChain,
				DestChain:        ackPacket.DstChain,
				Sequence:         ackPacket.Sequence,
				//Receiver:         ackPacket.Receiver,
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
				p.log.Infoln("GetPacketFee:%v", packetFee)
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
				p.log.Infoln("GetPacketFee:%v", packetFee)
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
		if pt.SrcChain == chains.TeleportChain {
			key := fmt.Sprintf("%v/%v/%v", pt.SrcChain, pt.DstChain, tokenName)
			crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainBToken.Address()
		} else {
			key := fmt.Sprintf("%v/%v/%v", pt.DstChain, pt.SrcChain, tokenName)
			crossChainTx.ReceiveTokenAddress = p.BridegesManager.BridgeTokenMap[key].ChainAToken.Address()
		}
		srcChain := p.Chains[pt.SrcChain]
		packetFee, err := srcChain.GetPacketFee(pt.SrcChain, pt.DstChain, int(pt.Sequence))
		if err != nil {
			p.log.Errorf("GetPacketFee error:%+v", err)
			return err
		} else {
			p.log.Infoln("GetPacketFee:%v", packetFee)
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
