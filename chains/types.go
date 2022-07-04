package chains

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

var AgentContract = "0x0000000000000000000000000000000040000001"

var chainMap ChainMap

const (
	DefaultTolerate float64 = 10
	TeleportChain           = "teleport"
	DefaultToken            = "tele"
	ZeroAddress             = "0x0000000000000000000000000000000000000000"
)

const (
	Pending int8 = iota + 1
	Success
	Fail
	Refund
)

type PacketTx struct {
	Packet    packettypes.Packet
	TxHash    string
	TimeStamp time.Time
	Height    uint64
	Signer    common.Address
	Gas       uint64
	GasPrice  float64
	MultiId   string
}

type BasePacketTx struct {
	SrcChain    string
	DstChain    string
	Sequence    uint64
	SrcChainId  string
	DestChainId string
	Sender      string
	Receiver    string
	Amount      string
	Token       string
	OriToken    string
	TxHash      string
	TimeStamp   time.Time
	Height      uint64
	Signer      common.Address
	Gas         uint64
	GasPrice    float64
	MultiId     string
	Code        int8
	ErrMsg      string
}

type BaseBlockPackets struct {
	Packets     []BasePacketTx
	AckPackets  []BasePacketTx
	RecivedAcks []BasePacketTx
}

type AckTx struct {
	Ack       AckPacket
	TxHash    string
	TimeStamp time.Time
	Height    uint64
	Gas       uint64
	GasPrice  float64
	MultiId   string
}

type BlockPackets struct {
	BizPackets  []PacketTx
	AckPackets  []AckTx
	RecivedAcks []AckTx
}

type AckPacket struct {
	Packet          packettypes.Packet
	Acknowledgement []byte
}

type MultiInfo struct {
	Id []byte `json:"id"`
	// with a later sequence number.
	Sequence *big.Int `json:"sequence"`
	// identifies the chain id of the sending chain.
	SrcChain string `json:"srcChain"`
	// identifies the chain id of the receiving chain.
	DestChain string `json:"destChain,omitempty" `
}

type AgentInfo struct {
	Id []byte `json:"id"`
	// with a later sequence number.
	Sequence *big.Int `json:"sequence"`
	// identifies the chain id of the receiving chain.
	DestChain string `json:"destChain,omitempty" `
}

type OutTokenAmount struct {
}

type TokenAmount struct {
	Amount   *big.Int
	Decimals uint8
	Scale    uint8
}

func (ta *TokenAmount) Add(tokenAmount *TokenAmount) {
	if ta.Decimals != tokenAmount.Decimals {
		panic("not support different decimals")
	}
	ta.Amount.Add(ta.Amount, tokenAmount.Amount)
}

func (ta *TokenAmount) Mul(tokenAmount *TokenAmount) {
	if ta.Decimals != tokenAmount.Decimals {
		panic("not support different decimals")
	}
	ta.Amount.Mul(ta.Amount, tokenAmount.Amount)
}

func (ta *TokenAmount) Sub(tokenAmount *TokenAmount) {
	if ta.Decimals != tokenAmount.Decimals {
		panic("not support different decimals")
	}
	ta.Amount.Sub(ta.Amount, tokenAmount.Amount)
}

func (ta *TokenAmount) Float64() (float64, error) {
	amountStr := ta.Amount.String()
	if len(amountStr) == 0 {
		return 0, nil
	}
	balanceStr := amountStr
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi error:%+v", err)
	}
	return balance, nil
}

type BigInt struct {
	*big.Int
}

func (Bi BigInt) Float64() (float64, error) {
	amountStr := Bi.String()
	if len(amountStr) == 0 {
		return 0, nil
	}
	balanceStr := amountStr
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi error:%+v", err)
	}
	return balance, nil
}

type Data struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

type EvmInToken struct {
	OriToken string
	Amount   *big.Int
	Bool     bool
}

type PacketFee struct {
	TokenAddress string   `json:"token_address"`
	FeeAmount    *big.Int `json:"fee_amount"`
}

type TokenLimit struct {
	Enable         bool     `json:"enable"`
	TimePeriod     *big.Int `json:"timePeriod"`
	TimeBasedLimit *big.Int `json:"timeBasedLimit"`
	MaxAmount      *big.Int `json:"maxAmount"`
	MinAmount      *big.Int `json:"minAmount"`
	PreviousTime   *big.Int `json:"previousTime"`
	CurrentSupply  *big.Int `json:"currentSupply"`
}

type ChainMap map[string]string

func NewChainMap() ChainMap {
	chainMap = make(map[string]string)
	return chainMap
}

func (cm ChainMap) GetIBCChainKey(chainName string) string {
	if chainName == TeleportChain {
		return fmt.Sprintf("%s-%s", chainName, "ibc")
	}
	return chainName
}

func (cm ChainMap) GetXIBCChainKey(chainName string) string {
	return chainName
}
