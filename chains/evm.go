package chains

import (
	"context"
	_ "embed" // embed compiled smart contract
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/teleport-network/teleport-data-analytics/config"

	"github.com/teleport-network/teleport-data-analytics/chains/contracts"
)

var (
	//go:embed token.json
	tokenJson []byte
	//go:embed teleport-endpoint.json
	EndpointJson []byte
	//go:embed evm-endpoint.json
	evmEndpointJson []byte
)

type Evm struct {
	chainName string
	ethClient *ethclient.Client
	packet    *contracts.Packet
	//proxy
	tokenContract    evmtypes.CompiledContract
	endpointContract *ethbind.BoundContract
	tokenDecimals    uint8
	packetAddr       string
	packetTopic      string
	ackTopic         string
	receivedAckTopic string
	agentAddr        string
	EndPointAddr     string
	frequency        int
	batchNumber      uint64
	nativeToken      string
	revisedHeight    uint64
	startHeight      uint64
	lastPriceQuery   time.Time
	lastGasPrice     *big.Int
}

func NewEvmCli(evmCfg config.EvmConfig) (*Evm, error) {
	chainCfg := evmCfg
	chainName := chainCfg.ChainName
	rpcClient, err := rpc.DialContext(context.Background(), chainCfg.EvmUrl)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	packAddr := common.HexToAddress(chainCfg.PacketContract)
	packetFilter, err := contracts.NewPacket(packAddr, ethClient)
	if err != nil {
		return nil, err
	}
	if evmCfg.EndPointAddr == "" {
		logrus.Fatal("empty transfer contract addr")
	}
	var transferContract evmtypes.CompiledContract
	if chainName == "teleport" {
		transferContract = getEndpointContract()
	} else {
		transferContract = getEvmEndpointContract()
	}
	EndPointAddress := common.HexToAddress(evmCfg.EndPointAddr)
	return &Evm{
		chainName:        chainName,
		ethClient:        ethClient,
		packet:           packetFilter,
		packetAddr:       chainCfg.PacketContract,
		packetTopic:      chainCfg.PacketTopic,
		ackTopic:         chainCfg.AckTopic,
		receivedAckTopic: chainCfg.ReceivedAckTopic,
		frequency:        chainCfg.Frequency,
		batchNumber:      chainCfg.BatchNumber,
		tokenContract:    getTokenContract(), // TODO fix
		endpointContract: ethbind.NewBoundContract(EndPointAddress, transferContract.ABI, ethClient, ethClient, ethClient),
		agentAddr:        chainCfg.AgentAddr,
		EndPointAddr:     evmCfg.EndPointAddr,
		nativeToken:      evmCfg.NativeToken,
		tokenDecimals:    evmCfg.NativeDecimals,
		startHeight:      evmCfg.StartHeight,
		revisedHeight:    evmCfg.RevisedHeight,
	}, nil
}

func getEvmEndpointContract() evmtypes.CompiledContract {
	var endpointContract evmtypes.CompiledContract
	if err := json.Unmarshal(evmEndpointJson, &endpointContract); err != nil {
		panic(fmt.Errorf("getEvmEndpointContract json.Unmarshal :%+v", err))
	}
	return endpointContract
}

// teleport transfer
func getEndpointContract() evmtypes.CompiledContract {
	var endpointContract evmtypes.CompiledContract
	if err := json.Unmarshal(EndpointJson, &endpointContract); err != nil {
		panic(fmt.Errorf("getEndpointContract json.Unmarshal :%+v", err))
	}
	return endpointContract
}

func (eth *Evm) GetBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	addr := common.HexToAddress(addressHex)
	var amount TokenAmount
	b, err := eth.ethClient.BalanceAt(context.Background(), addr, blockNumber)
	if err != nil {
		return amount, err
	}
	amount.Amount = b
	amount.Decimals = eth.tokenDecimals
	return amount, nil
}

func (eth *Evm) GetNativeDecimal() (uint8, error) {
	return eth.tokenDecimals, nil
}

func (eth *Evm) GetTokenLimit(addr common.Address, blockNumber *big.Int) (TokenLimit, error) {
	var (
		out        []interface{}
		tokenLimit TokenLimit
	)
	opts := new(ethbind.CallOpts)
	opts.BlockNumber = blockNumber
	if err := eth.endpointContract.Call(opts, &out, "limits", addr); err != nil {
		return TokenLimit{}, err
	}
	if len(out) < 7 {
		return TokenLimit{}, fmt.Errorf("invalid token limit length,length = %v", len(out))
	}
	var ok bool
	tokenLimit.Enable, ok = out[0].(bool)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid Enable type")
	}
	tokenLimit.TimePeriod, ok = out[1].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid TimePeriod type")
	}
	tokenLimit.TimeBasedLimit = out[2].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid TimeBasedLimit type")
	}
	tokenLimit.MaxAmount = out[3].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid MaxAmount type")
	}
	tokenLimit.MinAmount = out[4].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid MinAmount type")
	}
	tokenLimit.PreviousTime = out[5].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid PreviousTime type")
	}
	tokenLimit.CurrentSupply = out[6].(*big.Int)
	if !ok {
		return tokenLimit, fmt.Errorf("invalid CurrentSupply type")
	}
	return tokenLimit, nil
}

func (eth *Evm) GetInTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
	var (
		out    []interface{}
		amount TokenAmount
	)
	if addressHex == "" {
		return amount, fmt.Errorf("empty addressHex")
	}
	opts := new(ethbind.CallOpts)
	opts.BlockNumber = blockNumber
	var address interface{}
	if eth.ChainName() == "teleport" {
		if srcChain == "" {
			return amount, fmt.Errorf("empty srcChain")
		}
		address = fmt.Sprintf("%v/%v", addressHex, srcChain)
	} else {
		address = addressHex
	}
	if err := eth.endpointContract.Call(opts, &out, "getBindings", address); err != nil {
		return TokenAmount{}, err
	}
	if len(out) < 5 {
		return amount, fmt.Errorf("invalid bindings get,chainName:%v,len(out) = %v,data:%v", eth.ChainName(), len(out), out)
	}
	b, ok := out[4].(bool)
	if !ok {
		return amount, fmt.Errorf("invalid assert,chainName:%v,data:%v", eth.ChainName(), out[4])
	}
	if !b {
		return amount, fmt.Errorf("unbindings token")
	}
	scale, ok := out[3].(uint8)
	if !ok {
		return amount, fmt.Errorf("invalid assert,chainName:%v,data:%v", eth.ChainName(), out[3])
	}
	amount.Scale = scale
	a, ok := out[2].(*big.Int)
	if !ok {
		return amount, fmt.Errorf("invalid assert,chainName:%v,data:%v", eth.ChainName(), out[2])
	}
	amount.Amount = a
	return amount, nil
}

func (eth *Evm) GetOutTokenAmount(addressHex string, destChain string, blockNumber *big.Int) (TokenAmount, error) {
	var (
		out    []interface{}
		amount TokenAmount
	)
	opts := new(ethbind.CallOpts)
	opts.BlockNumber = blockNumber
	address := common.HexToAddress(addressHex)
	//if err := eth.transferContract.Call(opts, &out, "outTokens", fmt.Sprintf("%v/%v", addressHex, srcChain)); err != nil {
	if err := eth.endpointContract.Call(opts, &out, "outTokens", address, destChain); err != nil {
		return TokenAmount{}, err
	}
	if len(out) == 0 {
		return amount, fmt.Errorf("invalid bindings get,len(out) = 0")
	}
	tokenAmount, ok := out[0].(*big.Int)
	if !ok {
		return amount, fmt.Errorf("invalid type,data:%v", out[0])
	}
	amount.Amount = tokenAmount
	return amount, nil
}

func (eth *Evm) GetFrequency() int {
	return eth.frequency
}

func (eth *Evm) NativeToken() string {
	return eth.nativeToken
}

func (eth *Evm) GetBatchNumber() uint64 {
	return eth.batchNumber
}

func (eth *Evm) GetHeightByHash(hash string) (uint64, error) {
	hashByte := common.HexToHash(hash)
	//block, err := eth.ethClient.BlockByHash(context.Background(), hashByte)
	recept, err := eth.ethClient.TransactionReceipt(context.Background(), hashByte)
	if err != nil {
		return 0, err
	}
	return recept.BlockNumber.Uint64(), nil
}

func (eth *Evm) GetPackets(fromBlock, toBlock uint64) ([]*BaseBlockPackets, error) {
	Packets, err := eth.getPackets(fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	ackPackets, err := eth.getAckPackets(fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	receivedAcks, err := eth.getReceivedAcks(fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	packets := &BaseBlockPackets{
		Packets:     Packets,
		AckPackets:  ackPackets,
		RecivedAcks: receivedAcks,
	}
	return []*BaseBlockPackets{packets}, nil
}

func (eth *Evm) GetLightClientHeight(chainName string) (uint64, error) {
	return 0, nil
}

func (eth *Evm) GetPacketsByHash(txHash string) ([]BasePacketTx, error) {
	tx, err := eth.ethClient.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	return eth.handlePacketLog(tx.Logs)
}

func (eth *Evm) GetLatestHeight() (uint64, error) {
	return eth.ethClient.BlockNumber(context.Background())
}

func (eth *Evm) ChainName() string {
	return eth.chainName
}

func (eth *Evm) StartHeight() uint64 {
	return eth.startHeight
}

func (eth *Evm) RevisedHeight() uint64 {
	return eth.revisedHeight
}

func (eth *Evm) GetGasPrice() (*big.Int, error) {
	if time.Now().Sub(eth.lastPriceQuery) < 5*time.Second {
		return eth.lastGasPrice, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	lastGasPrice, err := eth.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	eth.lastGasPrice = lastGasPrice
	eth.lastPriceQuery = time.Now()
	return lastGasPrice, nil
}

func (eth *Evm) GetPacketFee(srcChain, dstChain string, sequence int) (*PacketFee, error) {
	key := fmt.Sprintf("%s/%s", dstChain, strconv.Itoa(sequence))
	packetFee, err := eth.packet.PacketFees(nil, []byte(key))
	if err != nil {
		log.Err(err).Str("key", key).Msg("fail to query PacketFees from contract ")
		return nil, err
	}
	log.Info().Str("key", key).Msgf("packet fee result: %v", packetFee)
	if packetFee.Amount == nil || packetFee.Amount.Sign() == 0 { // no fee
		return &PacketFee{FeeAmount: big.NewInt(0)}, nil
	}
	return &PacketFee{
		TokenAddress: packetFee.TokenAddress.Hex(),
		FeeAmount:    packetFee.Amount,
	}, nil
}

// get packets from block
func (eth *Evm) getPackets(fromBlock, toBlock uint64) ([]BasePacketTx, error) {
	address := common.HexToAddress(eth.packetAddr)
	topic := eth.packetTopic
	logs, err := eth.getLogs(address, topic, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}

	var Packets []BasePacketTx
	for _, log := range logs {
		packSent, err := eth.packet.ParsePacketSent(log)
		if err != nil {
			return nil, err
		}

		tx, _, err := eth.ethClient.TransactionByHash(context.Background(), log.TxHash)
		if err != nil {
			return nil, err
		}

		signer := ethtypes.NewLondonSigner(tx.ChainId())
		msg, err := tx.AsMessage(signer, tx.GasFeeCap())
		if err != nil {
			return nil, err
		}
		block, err := eth.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
		if err != nil {
			return nil, err
		}
		gasPrice, err := BigInt{tx.GasPrice()}.Float64()
		if err != nil {
			return nil, err
		}
		var packet packettypes.Packet
		if err := packet.ABIDecode(packSent.PacketBytes); err != nil {
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
			// transfer data. keep empty if not used.
			//TransferData: transferData.,
			Receiver: transferData.Receiver,
			Amount:   amount.String(),
			Token:    transferData.Token,
			OriToken: transferData.OriToken,

			TxHash:    log.TxHash.String(),
			Height:    log.BlockNumber,
			Signer:    msg.From(),
			TimeStamp: time.Unix(int64(block.Time()), 0).Local(),
			Gas:       tx.Gas(),
			GasPrice:  gasPrice,
		}
		tmpPack.SrcChainId = chainMap[chainMap.GetXIBCChainKey(tmpPack.SrcChain)]
		tmpPack.DestChainId = chainMap[chainMap.GetXIBCChainKey(tmpPack.DstChain)]

		Packets = append(Packets, tmpPack)
	}
	return Packets, nil
}

func (eth *Evm) handlePacketLog(logs []*ethtypes.Log) ([]BasePacketTx, error) {
	var Packets []BasePacketTx
	for _, log := range logs {
		packSent, err := eth.packet.ParsePacketSent(*log)
		if err != nil {
			continue
		}
		tx, _, err := eth.ethClient.TransactionByHash(context.Background(), log.TxHash)
		if err != nil {
			return nil, err
		}
		signer := ethtypes.NewLondonSigner(tx.ChainId())
		msg, err := tx.AsMessage(signer, tx.GasFeeCap())
		if err != nil {
			return nil, err
		}
		block, err := eth.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
		if err != nil {
			return nil, err
		}
		gasPrice, err := BigInt{tx.GasPrice()}.Float64()
		if err != nil {
			return nil, err
		}
		//var packet types.Packet
		var packet packettypes.Packet
		if err := packet.ABIDecode(packSent.PacketBytes); err != nil {
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
			Sequence:  packet.Sequence,
			SrcChain:  packet.SrcChain,
			DstChain:  packet.DstChain,
			Sender:    packet.Sender,
			Receiver:  transferData.Receiver,
			Amount:    amount.String(),
			Token:     transferData.Token,
			OriToken:  transferData.OriToken,
			TxHash:    log.TxHash.String(),
			Height:    log.BlockNumber,
			Signer:    msg.From(),
			TimeStamp: time.Unix(int64(block.Time()), 0).Local(),
			Gas:       tx.Gas(),
			GasPrice:  gasPrice,
		}

		Packets = append(Packets, tmpPack)
	}
	return Packets, nil
}

// get ack packets from block
func (eth *Evm) getAckPackets(fromBlock, toBlock uint64) ([]BasePacketTx, error) {
	address := common.HexToAddress(eth.packetAddr)
	topic := eth.ackTopic
	logs, err := eth.getLogs(address, topic, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}

	var ackPackets []BasePacketTx
	for _, log := range logs {
		ackWritten, err := eth.packet.ParseAckWritten(log)
		if err != nil {
			return nil, err
		}
		tx, _, err := eth.ethClient.TransactionByHash(context.Background(), log.TxHash)
		if err != nil {
			return nil, err
		}
		gasPrice, err := BigInt{tx.GasPrice()}.Float64()
		if err != nil {
			return nil, err
		}
		block, err := eth.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
		if err != nil {
			return nil, err
		}
		if ackWritten.Packet.TransferData == nil {
			continue
		}
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(ackWritten.Packet.TransferData); err != nil {
			return nil, err
		}
		var ack packettypes.Acknowledgement
		if err := ack.ABIDecode(ackWritten.Ack); err != nil {
			return nil, err
		}
		status, msg := eth.getStatus(ack)
		tmpAckPack := BasePacketTx{
			Sequence: ackWritten.Packet.Sequence,
			SrcChain: ackWritten.Packet.SrcChain,
			DstChain: ackWritten.Packet.DstChain,
			Code:     status,
			ErrMsg:   msg,
		}
		tmpAckPack.TxHash = log.TxHash.String()
		tmpAckPack.Height = log.BlockNumber
		tmpAckPack.TimeStamp = time.Unix(int64(block.Time()), 0).Local()
		tmpAckPack.Gas = tx.Gas()
		tmpAckPack.GasPrice = gasPrice
		ackPackets = append(ackPackets, tmpAckPack)
	}
	return ackPackets, nil
}

// get ack packets from block
func (eth *Evm) getReceivedAcks(fromBlock, toBlock uint64) ([]BasePacketTx, error) {
	address := common.HexToAddress(eth.packetAddr)
	topic := eth.ackTopic
	logs, err := eth.getLogs(address, topic, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	var ackPackets []BasePacketTx
	for _, log := range logs {
		packetAckPacket, err := eth.packet.ParseAckPacket(log)
		if err != nil {
			return nil, err
		}
		tx, _, err := eth.ethClient.TransactionByHash(context.Background(), log.TxHash)
		if err != nil {
			return nil, err
		}
		gasPrice, err := BigInt{tx.GasPrice()}.Float64()
		if err != nil {
			return nil, err
		}
		block, err := eth.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(log.BlockNumber)))
		if err != nil {
			return nil, err
		}
		if packetAckPacket.Packet.TransferData == nil {
			continue
		}
		var transferData packettypes.TransferData
		if err := transferData.ABIDecode(packetAckPacket.Packet.TransferData); err != nil {
			return nil, err
		}
		var ack packettypes.Acknowledgement
		if err := ack.ABIDecode(packetAckPacket.Ack); err != nil {
			return nil, err
		}
		status, _ := eth.getStatus(ack)
		if status != Fail {
			continue
		}
		tmpAckPack := BasePacketTx{
			Sequence: packetAckPacket.Packet.Sequence,
			SrcChain: packetAckPacket.Packet.SrcChain,
			DstChain: packetAckPacket.Packet.DstChain,
			Sender:   packetAckPacket.Packet.Sender,
			Receiver: transferData.Receiver,
			Code:     Refund,
		}
		tmpAckPack.TxHash = log.TxHash.String()
		tmpAckPack.Height = log.BlockNumber
		tmpAckPack.TimeStamp = time.Unix(int64(block.Time()), 0).Local()
		tmpAckPack.Gas = tx.Gas()
		tmpAckPack.GasPrice = gasPrice
		ackPackets = append(ackPackets, tmpAckPack)
	}
	return ackPackets, nil
}

func (eth *Evm) getLogs(address common.Address, topic string, fromBlock, toBlock uint64) ([]ethtypes.Log, error) {
	filter := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{ethcrypto.Keccak256Hash([]byte(topic))}},
	}
	return eth.ethClient.FilterLogs(context.Background(), filter)
}

func (eth *Evm) getAgentLogs(address common.Address, topic string, fromBlock, toBlock uint64) ([]ethtypes.Log, error) {
	filter1 := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
	}
	topics := [][]common.Hash{{ethcrypto.Keccak256Hash([]byte(topic))}}
	logs, err := eth.ethClient.FilterLogs(context.Background(), filter1)
	var ret []ethtypes.Log
	if err == nil {
		for _, log := range logs {
			if address != log.Address {
				continue
			}
			// If the to filtered topics is greater than the amount of topics in logs, skip.
			if len(topics) > len(log.Topics) {
				return nil, nil
			}
			ret = append(ret, log)
		}
	}
	return ret, err
}

func (eth *Evm) getStatus(ack packettypes.Acknowledgement) (int8, string) {
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
