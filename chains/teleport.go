package chains

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	"github.com/teleport-network/teleport/x/xibc/exported"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"github.com/teleport-network/teleport-sdk-go/client"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/teleport-network/teleport-data-analytics/config"
)

const DefaultDenom = "atele"

type Teleport struct {
	evm           *Evm
	tendermint    *TendermintClient
	teleportSDK   *client.TeleportClient
	Codec         *codec.ProtoCodec
	chainName     string
	frequency     int
	batchNumber   uint64
	revisedHeight uint64
	startHeight   uint64
}

func NewTeleport(cfg config.TendermintConfig, evmCli, tendermintClient BlockChain) *Teleport {
	cdc := makeCodec()
	evm, ok := evmCli.(*Evm)
	if !ok {
		panic("invalid evmCli")
	}
	tendermint, ok := tendermintClient.(*TendermintClient)
	if !ok {
		panic("invalid tendermintClient")
	}
	return &Teleport{
		evm:           evm,
		tendermint:    tendermint,
		Codec:         cdc,
		chainName:     cfg.ChainName,
		teleportSDK:   tendermint.tendermintSDK,
		frequency:     cfg.Frequency,
		batchNumber:   cfg.BatchNumber,
		startHeight:   cfg.StartHeight,
		revisedHeight: cfg.RevisedHeight,
	}
}

func (t *Teleport) GetBalance(address string) (string, error) {
	queryBalanceRequest := banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   DefaultDenom,
	}
	res, err := t.teleportSDK.BankQuery.Balance(context.Background(), &queryBalanceRequest)
	if err != nil {
		return "", err
	}
	return res.Balance.Amount.String(), nil
}

func (t *Teleport) GetFrequency() int {
	return t.frequency
}

func (t *Teleport) GetBatchNumber() uint64 {
	return t.batchNumber
}

func (t *Teleport) NativeToken() string {
	return t.evm.NativeToken()
}

func (t *Teleport) StartHeight() uint64 {
	return t.startHeight
}

func (t *Teleport) RevisedHeight() uint64 {
	return t.revisedHeight
}

func (t *Teleport) GetLatestHeight() (uint64, error) {
	block, err := t.teleportSDK.TMServiceQuery.GetLatestBlock(context.Background(), new(tmservice.GetLatestBlockRequest))
	if err != nil {
		return 0, err
	}
	var height = block.Block.Header.Height
	return uint64(height), err
}

func (t *Teleport) ChainName() string {
	return t.evm.ChainName()
}

func (t *Teleport) GetGasPrice() (*big.Int, error) {
	return t.evm.GetGasPrice()
}

func (t *Teleport) GetNativeDecimal() (uint8, error) {
	return t.evm.GetNativeDecimal()
}

func (t *Teleport) GetTokenLimit(addr common.Address, blockNumber *big.Int) (TokenLimit, error) {
	return t.evm.GetTokenLimit(addr, blockNumber)
}

func (t *Teleport) GetLightClientHeight(chainName string) (uint64, error) {
	ctx := context.Background()
	res, err := t.teleportSDK.XIBCClientQuery.ClientState(
		ctx,
		&clienttypes.QueryClientStateRequest{ChainName: chainName},
	)
	if err != nil {
		return 0, err
	}
	var clientState exported.ClientState
	if err := t.Codec.UnpackAny(res.ClientState, &clientState); err != nil {
		return 0, err
	}
	return clientState.GetLatestHeight().GetRevisionHeight(), nil
}

func (t *Teleport) GetPacketFee(srcChain, dstChain string, sequence int) (*PacketFee, error) {
	return t.evm.GetPacketFee(srcChain, dstChain, sequence)
}

func (t *Teleport) GetPackets(fromBlock, toBlock uint64) ([]*BaseBlockPackets, error) {
	times := toBlock - fromBlock + 1
	Packets := make([]*BaseBlockPackets, times)
	var l sync.Mutex
	var wg sync.WaitGroup
	wg.Add(int(times))
	var anyErr error
	for i := fromBlock; i <= toBlock; i++ {
		go func(h uint64) {
			defer wg.Done()
			var err error
			pkt, err := t.GetBlockPackets(h)
			if err != nil {
				anyErr = err
				return
			}
			l.Lock()
			Packets[h-fromBlock] = pkt
			l.Unlock()
		}(i)
	}
	wg.Wait()
	if anyErr != nil {
		return nil, anyErr
	}
	return Packets, nil
}

func (t *Teleport) GetPacketsByHash(txHash string) ([]BasePacketTx, error) {
	// Invalid query
	return t.evm.GetPacketsByHash(txHash)
}

func (t *Teleport) GetBlockPackets(height uint64) (*BaseBlockPackets, error) {
	var bizPackets []BasePacketTx
	var ackPackets []BasePacketTx
	var receivedAcks []BasePacketTx
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := t.teleportSDK.TMServiceQuery.GetBlockByHeight(ctx, &tmservice.GetBlockByHeightRequest{
		Height: int64(height),
	})
	if err != nil {
		return nil, err
	}
	var packets BaseBlockPackets
	for _, tx := range res.Block.GetData().Txs {
		hash := hex.EncodeToString(tmhash.Sum(tx))
		// TODO
		res, err := t.teleportSDK.TxClient.GetTx(ctx, &typestx.GetTxRequest{
			Hash: hash,
		})
		if err != nil {
			continue
		}
		txTime, err := time.Parse("2006-01-02T15:04:05Z", res.TxResponse.Timestamp)
		if err != nil {
			return nil, err
		}
		if len(res.TxResponse.Logs) == 0 {
			continue
		}
		stringEvents := res.TxResponse.Logs[0].Events
		ethereumTxHash, _ := t.getEthereumTxHash(stringEvents)
		tmpPackets, err := t.getPackets(stringEvents)
		if err != nil {
			return nil, err
		}
		var packetId int
		for i := range tmpPackets {
			// only packet has ethereumTxHash
			tmpPackets[i].TxHash = ethereumTxHash
			tmpPackets[i].TimeStamp = txTime
			tmpPackets[i].Height = height
			// Avoid parsing whenever possible
			if len(tmpPackets) > 1 {
				if tmpPackets[i].Sender == t.evm.agentAddr {
					tmpPackets[i].MultiId = fmt.Sprintf("%v/%v", hash, packetId)
					packetId++
				}
			} else {
				tmpPackets[i].MultiId = hash
			}
		}
		bizPackets = append(bizPackets, tmpPackets...)
		tmpAckPacks, err := t.getAckPackets(stringEvents)
		if err != nil {
			return nil, err
		}
		var ackId int
		for i := range tmpAckPacks {
			tmpAckPacks[i].TxHash = hash
			tmpAckPacks[i].TimeStamp = txTime
			tmpAckPacks[i].Height = height
			// Avoid parsing whenever possible
			if len(tmpAckPacks) > 1 {
				if tmpAckPacks[i].Receiver == t.evm.agentAddr && tmpAckPacks[i].Code == 0 {
					tmpAckPacks[i].MultiId = fmt.Sprintf("%v/%v", hash, ackId)
					ackId++
				}
			} else {
				tmpAckPacks[i].MultiId = hash
			}
		}
		ackPackets = append(ackPackets, tmpAckPacks...)
		tmpReceivedAck, err := t.getReceivedAcks(stringEvents)
		if err != nil {
			return nil, err
		}
		for i := range tmpReceivedAck {
			tmpReceivedAck[i].TxHash = hash
			tmpReceivedAck[i].TimeStamp = txTime
			tmpReceivedAck[i].Height = height
		}
		receivedAcks = append(receivedAcks, tmpReceivedAck...)
	}
	ibcPacket, err := t.tendermint.GetBlockPackets(height)
	if err != nil {
		return nil, err
	}
	packets.Packets = append(bizPackets, ibcPacket.Packets...)
	packets.AckPackets = append(ackPackets, ibcPacket.AckPackets...)
	packets.RecivedAcks = append(receivedAcks, ibcPacket.RecivedAcks...)
	return &packets, nil
}

func (t *Teleport) getPackets(stringEvents sdk.StringEvents) ([]BasePacketTx, error) {
	protoEvents := getEventsVals("xibc.core.packet.v1.EventSendPacket", stringEvents)
	var packets []BasePacketTx
	for _, protoEvent := range protoEvents {
		event, ok := protoEvent.(*packettypes.EventSendPacket)
		if !ok {
			return nil, fmt.Errorf("invalid type")
		}
		var packet packettypes.Packet
		if err := packet.ABIDecode(event.Packet); err != nil {
			return nil, err
		}
		if packet.TransferData == nil {
			continue
		}
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(packet.TransferData); err != nil {
			return nil, err
		}
		a := big.Int{}
		amount := a.SetBytes(transferData.Amount)
		tmpPack := BasePacketTx{
			Sequence: packet.Sequence,
			SrcChain: packet.SrcChain,
			DstChain: packet.DstChain,
			Sender:   packet.Sender,
			Receiver: transferData.Receiver,
			Amount:   amount.String(),
			Token:    transferData.Token,
			OriToken: transferData.OriToken,
		}
		tmpPack.SrcChainId = chainMap[chainMap.GetXIBCChainKey(tmpPack.SrcChain)]
		tmpPack.DestChainId = chainMap[chainMap.GetXIBCChainKey(tmpPack.DstChain)]
		packets = append(packets, tmpPack)
	}
	return packets, nil
}

func (t *Teleport) getEthereumTxHash(stringEvents sdk.StringEvents) (string, error) {
	return getValue("ethereum_tx", "ethereumTxHash", stringEvents)
}

func (t *Teleport) getAckPackets(stringEvents sdk.StringEvents) ([]BasePacketTx, error) {
	protoEvents := getEventsVals("xibc.core.packet.v1.EventWriteAck", stringEvents)
	var ackPackets []BasePacketTx
	for _, protoEvent := range protoEvents {
		event, ok := protoEvent.(*packettypes.EventWriteAck)
		if !ok {
			return nil, fmt.Errorf("proto parse failed")
		}
		var packet packettypes.Packet
		if err := packet.ABIDecode(event.Packet); err != nil {
			return nil, err
		}
		if packet.TransferData == nil {
			continue
		}
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(packet.TransferData); err != nil {
			return nil, err
		}
		var ack packettypes.Acknowledgement
		if err := ack.ABIDecode(event.Ack); err != nil {
			return nil, err
		}
		status, msg := t.getStatus(ack)
		a := big.Int{}
		amount := a.SetBytes(transferData.Amount)
		ackPacket := BasePacketTx{
			Sequence: packet.Sequence,
			SrcChain: packet.SrcChain,
			DstChain: packet.DstChain,
			Sender:   packet.Sender,
			// transfer data. keep empty if not used.
			Receiver: transferData.Receiver,
			Amount:   amount.String(),
			Token:    transferData.Token,
			OriToken: transferData.OriToken,
			Code:     status,
			ErrMsg:   msg,
		}
		ackPackets = append(ackPackets, ackPacket)
	}
	return ackPackets, nil
}

func (t *Teleport) getReceivedAcks(stringEvents sdk.StringEvents) ([]BasePacketTx, error) {
	protoEvents := getEventsVals("xibc.core.packet.v1.EventAcknowledgePacket", stringEvents)
	var ackPackets []BasePacketTx
	for _, protoEvent := range protoEvents {
		event, ok := protoEvent.(*packettypes.EventAcknowledgePacket)
		if !ok {
			return nil, fmt.Errorf("proto parse failed")
		}
		var packet packettypes.Packet
		if err := packet.ABIDecode(event.Packet); err != nil {
			return nil, err
		}
		if packet.TransferData == nil {
			continue
		}
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(packet.TransferData); err != nil {
			return nil, err
		}
		var ack packettypes.Acknowledgement
		if err := ack.ABIDecode(event.Ack); err != nil {
			return nil, err
		}
		status, _ := t.getStatus(ack)
		if status != Fail {
			continue
		}
		a := big.Int{}
		amount := a.SetBytes(transferData.Amount)
		ackPacket := BasePacketTx{
			Sequence: packet.Sequence,
			SrcChain: packet.SrcChain,
			DstChain: packet.DstChain,
			//Sender:   packet.Sender,
			// transfer data. keep empty if not used.
			//TransferData: transferData.,
			//Receiver:  transferData.Receiver,
			Amount:   amount.String(),
			Token:    transferData.Token,
			OriToken: transferData.OriToken,
		}
		ackPackets = append(ackPackets, ackPacket)
	}
	return ackPackets, nil
}

func getEventsVals(typ string, stringEvents sdk.StringEvents) []proto.Message {
	var events []proto.Message
	for _, e := range stringEvents {
		abciEvent := abci.Event{}
		if e.Type == typ {
			abciEvent.Type = e.Type
			for _, attr := range e.Attributes {
				abciEvent.Attributes = append(abciEvent.Attributes, abci.EventAttribute{
					Key:   []byte(attr.Key),
					Value: []byte(attr.Value),
				})
			}
			protoEvent, err := sdk.ParseTypedEvent(abciEvent)
			if err != nil {
				return nil
			}
			events = append(events, protoEvent)
		}
	}
	return events
}

func getValue(typ, key string, se sdk.StringEvents) (string, error) {
	for _, e := range se {
		if e.Type == typ {
			for _, attr := range e.Attributes {
				if attr.Key == key {
					return attr.Value, nil
				}
			}
		}
	}
	return "", fmt.Errorf("not found type:%s key:%s", typ, key)
}

func (t *Teleport) getStatus(ack packettypes.Acknowledgement) (int8, string) {
	var status int8
	var ackMsg string
	if ack.Code == 0 {
		status = Success
	} else {
		status = Fail
		ackMsg = ack.GetMessage()
	}
	return status, ackMsg
}
