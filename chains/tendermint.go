package chains

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/teleport-network/teleport-sdk-go/client"
	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/teleport-network/teleport-data-analytics/config"
)

type TendermintClient struct {
	chainName          string
	nativeToken        string
	tokenDecimals      uint8
	frequency          int
	batchNumber        uint64
	startHeight        uint64
	revisedHeight      uint64
	tendermintSDK      *client.TeleportClient
	channelToChainName map[string]string
	chainNameToChannel map[string]string
}

func NewTendermintClient(cfg config.TendermintConfig) (*TendermintClient, error) {
	cli, err := client.NewClient(cfg.Url, cfg.ChainID)
	if err != nil {
		panic(err)
	}
	channelToChainName := make(map[string]string)
	chainNameToChannel := make(map[string]string)
	for _, channelInfo := range cfg.ChannelInfo {
		channelToChainName[channelInfo.ChannelID] = channelInfo.ChainName
		chainNameToChannel[channelInfo.ChainName] = channelInfo.ChannelID
	}
	return &TendermintClient{
		chainName:          cfg.ChainName,
		tendermintSDK:      cli,
		frequency:          cfg.Frequency,
		batchNumber:        cfg.BatchNumber,
		startHeight:        cfg.StartHeight,
		revisedHeight:      cfg.RevisedHeight,
		channelToChainName: channelToChainName,
		chainNameToChannel: chainNameToChannel,
	}, err
}

func (c *TendermintClient) GetPackets(fromBlock, toBlock uint64) ([]*BaseBlockPackets, error) {
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
			pkt, err := c.GetBlockPackets(h)
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

func (c *TendermintClient) GetPackcetsByEvent() ([]*BaseBlockPackets, error) {
	fmt.Println(fmt.Sprintf("%s.%s=%s", channeltypes.EventTypeSendPacket, channeltypes.AttributeKeySequence, "1"))
	res, err := c.tendermintSDK.TxClient.GetTxsEvent(context.Background(), &typestx.GetTxsEventRequest{
		Events: []string{"send_packet.packet_sequence=1"},
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("res:", res)
	return nil, nil
}

func (c *TendermintClient) GetTokenLimit(addr common.Address, blockNumber *big.Int) (TokenLimit, error) {
	return TokenLimit{}, nil
}

func (c *TendermintClient) GetPacketsByHash(txHash string) ([]BasePacketTx, error) {
	return nil, nil
}

func (c *TendermintClient) GetTokenBalance(denom string, address string) (TokenAmount, error) {
	var tokenAmount TokenAmount
	res, err := c.tendermintSDK.BankQuery.Balance(context.Background(), &types.QueryBalanceRequest{
		Denom: denom, Address: address,
	})
	if err != nil {
		return TokenAmount{}, err
	}
	tokenAmount.Amount = res.Balance.Amount.BigInt()
	// TODO decimals
	return tokenAmount, nil
}

//func (c *TendermintClient)

func (c *TendermintClient) GetBlockPackets(height uint64) (*BaseBlockPackets, error) {
	var bizPacketTxs []BasePacketTx
	var ackPacketTxs []BasePacketTx
	var recivedAckTxs []BasePacketTx
	res, err := c.tendermintSDK.TMServiceQuery.GetBlockByHeight(context.Background(), &tmservice.GetBlockByHeightRequest{
		Height: int64(height),
	})
	if err != nil {
		return nil, err
	}
	var packets BaseBlockPackets
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
		tmpPackets, err := c.getPackets(stringEvents)
		if err != nil {
			return nil, err
		}
		for _, tmpPacket := range tmpPackets {
			var ibcPacketTx BasePacketTx
			if tmpPacket.SourcePort != "transfer" || tmpPacket.DestinationPort != "transfer" {
				continue
			}
			var ft ibctransfertypes.FungibleTokenPacketData
			if err := json.Unmarshal(tmpPacket.Data, &ft); err != nil {
				panic(err)
			}
			ibcPacketTx.SrcChain, err = c.getChainName(tmpPacket.SourceChannel)
			if err != nil {
				continue
			}
			ibcPacketTx.DstChain, err = c.getChainName(tmpPacket.DestinationChannel)
			if err != nil {
				continue
			}
			ibcPacketTx.SrcChainId = chainMap[chainMap.GetIBCChainKey(ibcPacketTx.SrcChain)]
			ibcPacketTx.DestChainId = chainMap[chainMap.GetIBCChainKey(ibcPacketTx.DstChain)]
			ibcPacketTx.Sender = ft.Sender
			ibcPacketTx.Receiver = ft.Receiver
			ibcPacketTx.Amount = ft.Amount
			ibcPacketTx.Token = ft.Denom
			ibcPacketTx.TxHash = hash
			ibcPacketTx.TimeStamp = txTime
			ibcPacketTx.Height = height
			bizPacketTxs = append(bizPacketTxs, ibcPacketTx)
		}

		ackPackets, err := c.getAckPackets(stringEvents)
		if err != nil {
			return nil, err
		}
		for _, ackPacket := range ackPackets {
			var ibcAckTx BasePacketTx
			ibcAckTx.SrcChain, err = c.getChainName(ackPacket.Packet.SourceChannel)
			if err != nil {
				continue
			}
			ibcAckTx.DstChain, err = c.getChainName(ackPacket.Packet.DestinationChannel)
			if err != nil {
				continue
			}
			var ack channeltypes.Acknowledgement
			if err := ack.Unmarshal(ackPacket.Ack); err != nil {
				return nil, err
			}
			if ack.Response != nil {
				switch ack.Response.(type) {
				case *channeltypes.Acknowledgement_Error:
					ibcAckTx.Code = 1
					ibcAckTx.ErrMsg = ack.Response.(*channeltypes.Acknowledgement_Error).Error
				default:
					ibcAckTx.Code = 0
				}
			}
			// TODO Gas
			ibcAckTx.TxHash = hash
			ibcAckTx.TimeStamp = txTime
			ibcAckTx.Height = height
			ackPacketTxs = append(ackPacketTxs, ibcAckTx)
		}
		tmpReceivedAcks, err := c.getRecievedAcks(stringEvents)
		if err != nil {
			return nil, fmt.Errorf("getRecievedAcks error:%v,chainName:%v,height:%v", err, c.chainName, height)
		}
		for _, tmpReceivedAck := range tmpReceivedAcks {
			var ibcAckReceivedTx BasePacketTx
			ibcAckReceivedTx.SrcChain, err = c.getChainName(tmpReceivedAck.Packet.SourceChannel)
			if err != nil {
				continue
			}
			ibcAckReceivedTx.DstChain, err = c.getChainName(tmpReceivedAck.Packet.DestinationChannel)
			if err != nil {
				continue
			}
			ibcAckReceivedTx.Sequence = tmpReceivedAck.Packet.Sequence
			ibcAckReceivedTx.TxHash = hash
			ibcAckReceivedTx.TimeStamp = txTime
			ibcAckReceivedTx.Height = height
			recivedAckTxs = append(recivedAckTxs, ibcAckReceivedTx)
		}
	}
	packets.Packets = bizPacketTxs
	packets.AckPackets = ackPacketTxs
	packets.RecivedAcks = recivedAckTxs
	return &packets, nil
}

func (c *TendermintClient) GetLightClientHeight(chainName string) (uint64, error) {
	// TODO
	return 0, nil
}
func (c *TendermintClient) GetNativeDecimal() (uint8, error) {
	return c.tokenDecimals, nil
}
func (c *TendermintClient) GetGasPrice() (*big.Int, error) {
	// TODO
	return nil, nil
}
func (c *TendermintClient) ChainName() string {
	return c.chainName
}
func (c *TendermintClient) GetLatestHeight() (uint64, error) {
	// TODO
	res, err := c.tendermintSDK.TMServiceQuery.GetLatestBlock(context.Background(), &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return uint64(res.Block.Header.Height), nil
}
func (c *TendermintClient) GetFrequency() int {
	return c.frequency
}
func (c *TendermintClient) NativeToken() string {
	return c.nativeToken
}
func (c *TendermintClient) GetBatchNumber() uint64 {
	return c.batchNumber
}
func (c *TendermintClient) StartHeight() uint64 {
	return c.startHeight
}
func (c *TendermintClient) RevisedHeight() uint64 {
	return c.revisedHeight
}

func (c *TendermintClient) GetPacketFee(srcChain, dstChain string, sequence int) (*PacketFee, error) {
	return nil, nil
}

type IBCAck struct {
	Packet channeltypes.Packet
	Ack    []byte
}

type IBCPacketTx struct {
	Packet    channeltypes.Packet
	TxHash    string
	TimeStamp time.Time
	Height    uint64
	Signer    common.Address
	Gas       uint64
	GasPrice  float64
	MultiId   string
}

type IBCAckTx struct {
	Ack       IBCAck
	TxHash    string
	TimeStamp time.Time
	Height    uint64
	Gas       uint64
	GasPrice  float64
	MultiId   string
}

type IBCBlockPackets struct {
	BizPackets  []IBCPacketTx
	AckPackets  []IBCAckTx
	RecivedAcks []IBCAckTx
}

func (c *TendermintClient) getPackets(stringEvents sdk.StringEvents) ([]channeltypes.Packet, error) {
	//var
	packetType := channeltypes.EventTypeSendPacket
	sequences := getValues(stringEvents, packetType, channeltypes.AttributeKeySequence)
	srcChannel := getValues(stringEvents, packetType, channeltypes.AttributeKeySrcChannel)
	srcPort := getValues(stringEvents, packetType, channeltypes.AttributeKeySrcPort)
	dstChannel := getValues(stringEvents, packetType, channeltypes.AttributeKeyDstChannel)
	dstPort := getValues(stringEvents, packetType, channeltypes.AttributeKeyDstPort)
	packetDatas := getValues(stringEvents, packetType, channeltypes.AttributeKeyData)
	if len(sequences) == 0 {
		return nil, nil
	}
	if !(len(sequences) == len(srcChannel) && len(sequences) == len(srcPort) && len(sequences) == len(dstChannel) && len(sequences) == len(dstPort) && len(sequences) == len(packetDatas)) {
		return nil, fmt.Errorf("invalid getAckPackets")
	}
	var packets []channeltypes.Packet
	for i := 0; i < len(sequences); i++ {
		sequenceStr := sequences[i]
		sequence, err := strconv.Atoi(sequenceStr)
		if err != nil {
			return nil, err
		}
		tmpPack := channeltypes.Packet{
			Sequence:           uint64(sequence),
			SourceChannel:      srcChannel[i],
			DestinationChannel: dstChannel[i],
			SourcePort:         srcPort[i],
			DestinationPort:    dstPort[i],
			Data:               []byte(packetDatas[i]),
		}
		packets = append(packets, tmpPack)
	}
	return packets, nil
}

func (c *TendermintClient) getAckPackets(stringEvents sdk.StringEvents) ([]IBCAck, error) {
	ackType := channeltypes.EventTypeWriteAck
	sequences := getValues(stringEvents, ackType, channeltypes.AttributeKeySequence)
	srcChannel := getValues(stringEvents, ackType, channeltypes.AttributeKeySrcChannel)
	srcPort := getValues(stringEvents, ackType, channeltypes.AttributeKeySrcPort)
	dstChannel := getValues(stringEvents, ackType, channeltypes.AttributeKeyDstChannel)
	dstPort := getValues(stringEvents, ackType, channeltypes.AttributeKeyDstPort)
	packetDatas := getValues(stringEvents, ackType, channeltypes.AttributeKeyData)
	acknowledgements := getValues(stringEvents, ackType, channeltypes.AttributeKeyAck)
	if len(sequences) == 0 {
		return nil, nil
	}
	if !(len(sequences) == len(srcChannel) && len(sequences) == len(srcPort) && len(sequences) == len(dstChannel) && len(sequences) == len(dstPort) && len(sequences) == len(packetDatas) && len(sequences) == len(acknowledgements)) {
		return nil, fmt.Errorf("invalid getAckPackets")
	}
	var ibcAcks []IBCAck
	for i := 0; i < len(sequences); i++ {
		sequenceStr := sequences[i]
		sequence, err := strconv.Atoi(sequenceStr)
		if err != nil {
			return nil, err
		}
		tmpPack := channeltypes.Packet{
			Sequence:           uint64(sequence),
			SourceChannel:      srcChannel[i],
			DestinationChannel: dstChannel[i],
			SourcePort:         srcPort[i],
			DestinationPort:    dstPort[i],
			Data:               []byte(packetDatas[i]),
		}
		var ibcAck IBCAck
		ibcAck.Packet = tmpPack
		ibcAck.Ack = []byte(acknowledgements[i])
		ibcAcks = append(ibcAcks, ibcAck)
	}
	return ibcAcks, nil
}

func (c *TendermintClient) getRecievedAcks(stringEvents sdk.StringEvents) ([]IBCAck, error) {
	ackReceivedType := channeltypes.EventTypeAcknowledgePacket
	sequences := getValues(stringEvents, ackReceivedType, channeltypes.AttributeKeySequence)
	srcChannel := getValues(stringEvents, ackReceivedType, channeltypes.AttributeKeySrcChannel)
	srcPort := getValues(stringEvents, ackReceivedType, channeltypes.AttributeKeySrcPort)
	dstChannel := getValues(stringEvents, ackReceivedType, channeltypes.AttributeKeyDstChannel)
	dstPort := getValues(stringEvents, ackReceivedType, channeltypes.AttributeKeyDstPort)
	if len(sequences) == 0 {
		return nil, nil
	}
	if !(len(sequences) == len(srcChannel) && len(sequences) == len(srcPort) && len(sequences) == len(dstChannel) && len(sequences) == len(dstPort)) {
		return nil, fmt.Errorf("invalid getRecievedAcks")
	}
	var ibcAcks []IBCAck
	for i := 0; i < len(sequences); i++ {
		sequenceStr := sequences[i]
		sequence, err := strconv.Atoi(sequenceStr)
		if err != nil {
			return nil, err
		}
		tmpPack := channeltypes.Packet{
			Sequence:           uint64(sequence),
			SourceChannel:      srcChannel[i],
			DestinationChannel: dstChannel[i],
			SourcePort:         srcPort[i],
			DestinationPort:    dstPort[i],
		}
		var ibcAck IBCAck
		ibcAck.Packet = tmpPack
		ibcAcks = append(ibcAcks, ibcAck)
	}
	return ibcAcks, nil
}

func getValues(se sdk.StringEvents, typ, key string) (res []string) {
	for _, e := range se {
		if e.Type == typ {
			for _, attr := range e.Attributes {
				if attr.Key == key {
					res = append(res, attr.Value)
				}
			}
		}
	}
	return res
}

func (c *TendermintClient) getChainName(channel string) (string, error) {
	if channel == "" {
		return "", fmt.Errorf("invalid channel")
	}
	if c.channelToChainName[channel] == "" {
		return "", fmt.Errorf("chainName not found,channel = %v", channel)
	}
	return c.channelToChainName[channel], nil
}

func (c *TendermintClient) getChannelID(chainName string) (string, error) {
	if chainName == "" {
		return "", fmt.Errorf("invalid chainName")
	}
	if c.chainNameToChannel[chainName] == "" {
		return "", fmt.Errorf("channel not found,chainName = %v", chainName)
	}
	return c.chainNameToChannel[chainName], nil
}

func (c *TendermintClient) getStatus(ack packettypes.Acknowledgement) (int8, string) {
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
