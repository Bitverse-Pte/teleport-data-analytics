package chains

import (
	"encoding/json"
	"fmt"
	"math/big"

	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

const DefaultNativeToken = "0x0000000000000000000000000000000000000000"

type tokenQuery struct {
	*Evm
	tokenName  string
	addressHex string
	contract   *ethbind.BoundContract
	decimals   uint8
}

func NewTokenQuery(eth *Evm, addressHex, tokenName string) TokenQuery {
	token := getTokenContract()
	address := common.HexToAddress(addressHex)
	tQuery := &tokenQuery{
		contract:   ethbind.NewBoundContract(address, token.ABI, eth.ethClient, eth.ethClient, eth.ethClient),
		Evm:        eth,
		addressHex: addressHex,
		tokenName:  tokenName,
	}
	decimals, err := tQuery.GetDecimals()
	if err != nil {
		panic(fmt.Errorf("NewTokenQuery error:%+v\n,chainName:%v,addressHex:%v,tokenName:%v", err, eth.ChainName(), addressHex, tokenName))
	}
	fmt.Printf("NewTokenQuery success,chainName:%v,addressHex:%v,tokenName:%v", eth.ChainName(), addressHex, tokenName)
	tQuery.decimals = decimals
	return tQuery
}

func getTokenContract() evmtypes.CompiledContract {
	var tokenContract evmtypes.CompiledContract
	if err := json.Unmarshal(tokenJson, &tokenContract); err != nil {
		panic(fmt.Errorf("getTokenContract json.Unmarshal error:%+v", err))
	}
	return tokenContract
}

func (t *tokenQuery) AddressHex() string {
	return t.addressHex
}

func (t *tokenQuery) TokenName() string {
	return t.tokenName
}

func (t *tokenQuery) GetErc20TotalSupply(blockNumber *big.Int) (TokenAmount, error) {
	var out []interface{}
	var total TokenAmount
	opts := new(ethbind.CallOpts)
	opts.BlockNumber = blockNumber
	if err := t.contract.Call(opts, &out, "totalSupply"); err != nil {
		return total, err
	}
	if len(out) == 0 {
		return total, fmt.Errorf("invalid totalSupply,len = 0")
	}
	decimals, err := t.getErc20Decimals()
	if err != nil {
		return total, fmt.Errorf("get decimals error:%v", err)
	}
	total.Decimals = decimals
	if totalSupply, ok := out[0].(*big.Int); ok {
		total.Amount = totalSupply
		return total, nil
	}
	return TokenAmount{}, fmt.Errorf("invalid totalSupply type")
}

func (t *tokenQuery) QueryTotalBurnToken() (TokenAmount, error) {
	return TokenAmount{}, nil
}

func (t *tokenQuery) getErc20Decimals() (uint8, error) {
	if t.decimals != 0 {
		return t.decimals, nil
	}
	var out []interface{}
	if err := t.contract.Call(nil, &out, "decimals"); err != nil {
		return 0, err
	}
	if len(out) == 0 {
		return 0, fmt.Errorf("invalid decimals,len = 0")
	}
	if decimals, ok := out[0].(uint8); ok {
		t.decimals = decimals
		return decimals, nil
	}
	return 0, fmt.Errorf("invalid decimals type")
}
func (t *tokenQuery) GetDecimals() (uint8, error) {
	if t.NativeToken() == t.TokenName() {
		return t.Evm.GetNativeDecimal()
	}
	return t.getErc20Decimals()
}

func (t *tokenQuery) GetTokenBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	var balanceQuery func(addressHex string, blockNumber *big.Int) (TokenAmount, error)
	if t.addressHex == DefaultNativeToken {
		balanceQuery = t.Evm.GetBalance
	} else {
		balanceQuery = t.GetErc20Balance
	}
	return balanceQuery(addressHex, blockNumber)
}

func (t *tokenQuery) GetErc20Balance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	address := common.HexToAddress(addressHex)
	var (
		out    []interface{}
		amount TokenAmount
	)
	opts := new(ethbind.CallOpts)
	opts.BlockNumber = blockNumber
	if err := t.contract.Call(opts, &out, "balanceOf", address); err != nil {
		return TokenAmount{}, err
	}
	if len(out) == 0 {
		return amount, fmt.Errorf("invalid balance,len(out) = 0")
	}
	if balance, ok := out[0].(*big.Int); ok {
		decimals, err := t.getErc20Decimals()
		if err != nil {
			return amount, err
		}
		amount.Amount = balance
		amount.Decimals = decimals
		return amount, nil
	}
	return amount, fmt.Errorf("invalid balance type")
}

func (t *tokenQuery) GetTransferBalance(blockNumber *big.Int) (TokenAmount, error) {
	var (
		amount TokenAmount
		err    error
	)
	if t.tokenName == t.Evm.NativeToken() {
		amount, err = t.Evm.GetBalance(t.Evm.EndPointAddr, blockNumber)
	} else {
		amount, err = t.GetErc20Balance(t.Evm.EndPointAddr, blockNumber)
	}
	return amount, err
}

func (t *tokenQuery) GetInTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
	decimals, err := t.GetDecimals()
	if err != nil {
		return TokenAmount{}, err
	}
	tokenAmount, err := t.Evm.GetInTokenAmount(addressHex, srcChain, blockNumber)
	if err != nil {
		return TokenAmount{}, err
	}
	tokenAmount.Decimals = decimals
	return tokenAmount, nil
}

func (t *tokenQuery) GetOutTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
	decimals, err := t.GetDecimals()
	if err != nil {
		return TokenAmount{}, err
	}
	tokenAmount, err := t.Evm.GetOutTokenAmount(addressHex, srcChain, blockNumber)
	if err != nil {
		return TokenAmount{}, err
	}
	tokenAmount.Decimals = decimals
	return tokenAmount, nil
}
