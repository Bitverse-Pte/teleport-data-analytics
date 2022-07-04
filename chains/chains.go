package chains

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type BlockChain interface {
	GetPackets(fromBlock, toBlock uint64) ([]*BaseBlockPackets, error)
	GetPacketsByHash(txHash string) ([]BasePacketTx, error)
	GetLightClientHeight(chainName string) (uint64, error)
	GetTokenLimit(addr common.Address, blockNumber *big.Int) (TokenLimit, error)
	GetNativeDecimal() (uint8, error)
	GetGasPrice() (*big.Int, error)
	ChainName() string
	GetLatestHeight() (uint64, error)
	GetFrequency() int
	NativeToken() string
	GetBatchNumber() uint64
	StartHeight() uint64
	RevisedHeight() uint64
	GetPacketFee(srcChain, dstChain string, sequence int) (*PacketFee, error)
}

type TokenQuery interface {
	GetDecimals() (uint8, error)
	GetErc20Balance(addressHex string, blockNumber *big.Int) (TokenAmount, error)
	GetErc20TotalSupply(blockNumber *big.Int) (TokenAmount, error)
	//GetBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error)
	GetTokenBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error)
	GetTransferBalance(blockNumber *big.Int) (TokenAmount, error)
	// Only teleport evm
	GetInTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error)
	// Only teleport evm
	GetOutTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error)
	ChainName() string
	NativeToken() string
	Address() string
	TokenName() string
}
