package model

import (
	"time"

	"gorm.io/gorm"
)

type PacketStatus int8

const (
	Pending PacketStatus = iota + 1
	Success
	Fail
	Refund
)

type SyncState struct {
	ChainName string `gorm:"primary_key"`
	Height    uint64 `gorm:"default:1"`
}

// single direction bridge metrics
// index idx_from_to
type SingleDirectionBridgeMetrics struct {
	gorm.Model
	SrcChain  string `gorm:"index:idx_from_to;comment:'src chain name'"` // chain name
	DestChain string `gorm:"index:idx_from_to"`                          // chain name
	PktAmt    int64  `gorm:"default:0;type:int"`                         // packet amount
	//UsdtFeeAmt string `gorm:"default:0;"`                 // usdt fee amount
	//TeleFeeAmt string `gorm:"default:0;"`                 // tele fee amount
	FailedAmt int64 `gorm:"default:0;"` // failed amount

	// bridge metrics contains token info
	//Token Token `gorm:"foreignkey:bridge_id;constraint:OnDelete:cascade;"` // token

}

// token table
type Token struct {
	gorm.Model
	BridgeID  uint
	TokenType int64 `gorm:"default:0;"` // token type:[0:packet_fee,1:packet_value,2:liquidity_fee,3:liquidity_value]
	TokenName string
	Amt       string `gorm:"default:0;"` // amount
}

// global bridge metrics
type GlobalMetrics struct {
	ID        int64     `gorm:"primary_key,auto_increment,unsigned"` // primary key
	UserAmt   int64     `gorm:"default:0;type:int"`                  // user amount
	CreatedAt time.Time `gorm:"type:timestamp"`                      // created time
}

type CrossPacket struct {
	ID            string
	PacketType    string
	SrcChain      string
	DestChain     string
	CommitmentKey string
	Height        uint64
	TxHash        string
	TxTime        time.Time
	Sender        string
	Receiver      string
	Token         string
	OriToken      string
	Amount        string
	Status        int8
	ErrMessage    string
	Sequence      uint64
	MultiID       string
}

type CrossChainTransaction struct {
	gorm.Model
	SrcChain             string    `gorm:"type:varchar(20);uniqueIndex:ck_t"`
	RelayChain           string    `gorm:"type:varchar(20);uniqueIndex:ck_t"`
	DestChain            string    `gorm:"type:varchar(20);uniqueIndex:ck_t"`
	Sequence             uint64    `gorm:"column:sequence;uniqueIndex:ck_t"`
	SrcChainId           string    `gorm:"column:src_chain_id"`
	DestChainId          string    `gorm:"column:dest_chain_id"`
	Amount               string    `gorm:"column:amount"` // TODO to float
	SrcAmountFloat       float64   `gorm:"column:src_amount_float"`
	AmountFloat          float64   `gorm:"column:amount_float"`
	AmountRaw            string    `gorm:"column:amount_raw"`
	Status               int8      `gorm:"tinyint;default:0"` // 1:pending,2:success,3:err,4:refund
	ErrMessage           string    `gorm:"type:varchar(100);default:''"`
	SrcHeight            uint64    `gorm:"column:src_height;index"`
	Sender               string    `gorm:"column:sender;index"`
	SendTxHash           string    `gorm:"column:send_tx_hash;index"`
	SendTxTime           time.Time `gorm:"column:send_tx_time;type:timestamp;index"`
	SendTokenAddress     string    `gorm:"column:send_token_address;index"`
	DestHeight           uint64    `gorm:"column:dest_height;index"`
	Receiver             string    `gorm:"column:receiver;index"`
	ReceiveTxHash        string    `gorm:"column:receive_tx_hash;index"`
	ReceiveTxTime        time.Time `gorm:"column:receive_tx_time;type:timestamp;index"`
	ReceiveTokenAddress  string    `gorm:"column:receive_token_address;index"`
	RefundHeight         uint64    `gorm:"column:refund_height;index"`
	RefundTxHash         string    `gorm:"column:refund_tx_hash;index"`
	RefundTxTime         time.Time `gorm:"column:refund_tx_time;type:timestamp;index"`
	MultiID1             string    `gorm:"column:multi_id1;index"`
	MultiID2             string    `gorm:"column:multi_id2;index"`
	TokenName            string    `gorm:"column:token_name;type:varchar(20);default:''"`
	SrcAmount            float64   `gorm:"column:src_amount"`
	SrcTrimAmount        float64   `gorm:"column:src_trim_amount"`
	DestAmount           float64   `gorm:"column:dest_amount"`
	DestTrimAmount       float64   `gorm:"column:dest_trim_amount"`
	DifferenceValue      float64   `gorm:"column:difference_value"`
	ReconcileResult      string    `gorm:"column:reconcile_result;default:''"`
	ReconcileMsg         string    `gorm:"column:reconcile_msg;default:''"`
	PacketGas            float64   `gorm:"column:packet_gas"`
	PacketFee            float64   `gorm:"column:packet_fee"`
	PacketGasPrice       float64   `gorm:"column:packet_gas_price"`
	AckGas               float64   `gorm:"column:ack_gas"`
	AckFee               float64   `gorm:"column:ack_fee"`
	AckGasPrice          float64   `gorm:"column:ack_gas_price"`
	PacketFeeDemand      float64   `gorm:"column:packet_fee_demand"`
	PacketFeePaid        float64   `gorm:"column:packet_fee_paid"`
	PacketFeeDemandToken string    `gorm:"column:packet_fee_demand_token"`
	AckFeeDemand         float64   `gorm:"column:ack_fee_demand"`
	AckFeeDemandToken    string    `gorm:"column:ack_fee_demand_token"`
	AckFeePaid           float64   `gorm:"column:ack_fee_paid"`
}

type PacketRelation struct {
	gorm.Model
	SrcChain  string `gorm:"column:src_chain;index"`
	DestChain string `gorm:"column:dest_chain;index"`
	Sequence  uint64 `gorm:"column:sequence;index"`
	PacketID  string `gorm:"column:packet_id;index"`
}

type Record struct {
	gorm.Model
	Sequence  string
	TokenName string
	ChainName string
	Amount    float64
	Result    bool
}

type BridgeReconcileResult struct {
	gorm.Model
	TokenName       string
	SrcChain        string
	SrcToken        string
	SrcAmount       float64
	SrcTrimAmount   float64
	DestChain       string
	DestToken       string
	DestAmount      float64
	DestTrimAmount  float64
	DifferenceValue float64
	Result          string
	ResultMsg       string
	ErrorMsg        string
}

type BlockReconcileResult struct {
	gorm.Model
	ChainName         string
	TokenName         string
	TokenAddress      string
	TokenAmount       float64
	Height            uint64
	LatestBlockAmount float64
	LastBlockAmount   float64
	DifferenceValue   float64
	Result            bool
	ResultMsg         string
	ErrorMsg          string
}
