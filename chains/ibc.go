package chains

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/teleport-network/teleport-sdk-go/client"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/teleport-network/teleport-data-analytics/config"
)

const (
	EventTypeSendPacket        = "send_packet"
	EventTypeRecvPacket        = "recv_packet"
	EventTypeWriteAck          = "write_acknowledgement"
	EventTypeAcknowledgePacket = "acknowledge_packet"
)

type IBCClient struct {
	chainName     string
	nativeToken   string
	tokenDecimals uint8
	frequency     int
	batchNumber   uint64
	startHeight   uint64
	revisedHeight uint64
	tendermintSDK *client.TeleportClient
}

func NewIBCClient(cfg config.TendermintConfig) (*IBCClient, error) {
	cli, err := client.NewClient(cfg.Url, cfg.ChainID)
	if err != nil {
		panic(err)
	}
	return &IBCClient{
		chainName:     cfg.ChainName,
		tendermintSDK: cli,
		frequency:     cfg.Frequency,
		batchNumber:   cfg.BatchNumber,
		startHeight:   cfg.StartHeight,
		revisedHeight: cfg.RevisedHeight,
	}, err
}

func (c *IBCClient) GetPackcets(fromBlock, toBlock uint64) ([]*BlockPackets, error) {
	return nil, nil
}

func (c *IBCClient) GetBlockPackets(height uint64) (*BlockPackets, error) {
	var bizPackets []PacketTx
	var ackPackets []AckTx
	var recivedAcks []AckTx
	res, err := c.tendermintSDK.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{
		Height: int64(height),
	})
	if err != nil {
		return nil, err
	}
	var packets BlockPackets
	for _, tx := range res.Block.GetData().Txs {
		hash := hex.EncodeToString(tmhash.Sum(tx))
		res, err := c.tendermintSDK.TxClient.GetTx(context.Background(), &typestx.GetTxRequest{
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
		//ethereumTxHash, _ := c.getEthereumTxHash(stringEvents)
		tmpPacket, err := c.getPacket(stringEvents)
		if err != nil {
			return nil, err
		}
		tmpPacket.TxHash = hash
		tmpPacket.TimeStamp = txTime
		tmpPacket.Height = height
		bizPackets = append(bizPackets, tmpPacket)
		//tmpAckPacks, err := c.getAckPackets(stringEvents)
		//if err != nil {
		//	return nil, err
		//}
		//for i := range tmpAckPacks {
		//	tmpAckPacks[i].TxHash = hash
		//	tmpAckPacks[i].TimeStamp = txTime
		//	tmpAckPacks[i].Height = height
		//}
		//ackPackets = append(ackPackets, tmpAckPacks...)
		tmpReceivedAck, err := c.getRecievedAcks(stringEvents)
		if err != nil {
			return nil, err
		}
		for i := range tmpReceivedAck {
			tmpReceivedAck[i].TxHash = hash
			tmpReceivedAck[i].TimeStamp = txTime
			tmpReceivedAck[i].Height = height
		}
		recivedAcks = append(recivedAcks, tmpReceivedAck...)
	}
	packets.BizPackets = bizPackets
	packets.AckPackets = ackPackets
	packets.RecivedAcks = recivedAcks
	return &packets, nil
}

func (c *IBCClient) GetMultiInfo(fromBlock, toBlock uint64) ([]MultiInfo, error) {
	return nil, nil
}
func (c *IBCClient) GetLightClientHeight(chainName string) (uint64, error) {
	// TODO
	return 0, nil
}
func (c *IBCClient) GetNativeDecimal() (uint8, error) {
	return c.tokenDecimals, nil
}
func (c *IBCClient) GetGasPrice() (*big.Int, error) {
	// TODO
	return nil, nil
}
func (c *IBCClient) ChainName() string {
	return c.chainName
}
func (c *IBCClient) GetLatestHeight() (uint64, error) {
	// TODO
	return 0, nil
}
func (c *IBCClient) GetFrequency() int {
	return c.frequency
}
func (c *IBCClient) NativeToken() string {
	return c.nativeToken
}
func (c *IBCClient) GetBatchNumber() uint64 {
	return c.batchNumber
}
func (c *IBCClient) StartHeight() uint64 {
	return c.startHeight
}
func (c *IBCClient) RevisedHeight() uint64 {
	return c.revisedHeight
}

func (c *IBCClient) getPacket(stringEvents sdk.StringEvents) (PacketTx, error) {
	events := getEvents(channeltypes.EventTypeSendPacket, stringEvents)
	if len(events) == 0 {
		return PacketTx{}, nil
	}
	//packet, err := ibctesting.ParsePacketFromEvents(events)
	//if err != nil {
	//	return PacketTx{}, err
	//}
	//var packets []PacketTx
	tmpPack := PacketTx{
		Packet: packettypes.Packet{
			//Sequence:         packet.Sequence,
			//SourceChain:      packet.SourceChannel,
			//DestinationChain: packet.DestinationChannel,
			////Ports:            event.GetPorts(),
			////RelayChain:       packet.,
			//DataList: [][]byte{packet.Data},
		},
	}
	return tmpPack, nil
}

func (c *IBCClient) getAckPackets(stringEvents sdk.StringEvents) (AckTx, error) {
	events := getEvents(channeltypes.EventTypeWriteAck, stringEvents)
	if len(events) == 0 {
		return AckTx{}, nil
	}
	//channeltypes.Packet{}
	//channeltypes.Acknowledgement{}
	//ack,err := ibctesting.ParseAckFromEvents(events)
	//if err != nil {
	//	return AckTx{},err
	//}

	//protoEvents := getEventsVals("xibc.core.packet.v1.EventWriteAck", stringEvents)
	//var ackPackets []AckTx
	//for _, protoEvent := range protoEvents {
	//	event, ok := protoEvent.(*packettypes.EventWriteAck)
	//	if !ok {
	//		return nil, fmt.Errorf("proto parse failed")
	//	}
	//	sequence, err := strconv.Atoi(event.GetSequence())
	//	if err != nil {
	//		return nil, err
	//	}
	//	tmpPack := packettypes.Packet{
	//		Sequence:         uint64(sequence),
	//		SourceChain:      event.GetSrcChain(),
	//		DestinationChain: event.GetDstChain(),
	//		Ports:            event.GetPorts(),
	//		RelayChain:       event.RelayChain,
	//		DataList:         event.GetDataList(),
	//	}
	//	var ackPacket AckTx
	//	ackPacket.Ack.Packet = tmpPack
	//	ackPacket.Ack.Acknowledgement = event.Ack
	//	ackPackets = append(ackPackets, ackPacket)
	////}
	//return ackPackets, nil
	return AckTx{}, nil
}

func (c *IBCClient) getRecievedAcks(stringEvents sdk.StringEvents) ([]AckTx, error) {
	protoEvents := getEventsVals("xibc.core.packet.v1.EventAcknowledgePacket", stringEvents)
	var ackPackets []AckTx
	for _, protoEvent := range protoEvents {
		event, ok := protoEvent.(*packettypes.EventAcknowledgePacket)
		if !ok {
			return nil, fmt.Errorf("proto parse failed")
		}
		//sequence, err := strconv.Atoi(event.GetSequence())
		//if err != nil {
		//	return nil, err
		//}
		tmpPack := packettypes.Packet{
			//Sequence:         uint64(sequence),
			//SourceChain:      event.GetSrcChain(),
			//DestinationChain: event.GetDstChain(),
			//Ports:            event.GetPorts(),
			//RelayChain:       event.RelayChain,
			//DataList:         event.GetDataList(),
		}
		var ackPacket AckTx
		ackPacket.Ack.Packet = tmpPack
		ackPacket.Ack.Acknowledgement = event.Ack
		ackPackets = append(ackPackets, ackPacket)
	}
	return ackPackets, nil
}

func getEvents(typ string, stringEvents sdk.StringEvents) sdk.Events {
	var events sdk.Events
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
			events = append(events, sdk.Event(abciEvent))
		}
	}
	return events
}
