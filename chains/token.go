package chains

import (
	"encoding/json"
	"fmt"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"math/big"
	"strings"
)

const DefaultNativeToken = "0x0000000000000000000000000000000000000000"

func NewTokenQuery(chain BlockChain, tokenAddress, tokenName string, decimals uint8) TokenQuery {
	if strings.Contains(tokenAddress, "0x") {
		evmcli, ok := chain.(*Evm)
		if !ok {
			panic(fmt.Sprintf("invalid erc20 token,chain:%s,tokenName:%s,tokenAddress:%s", chain.ChainName(), tokenName, tokenAddress))
		}
		return NewErc20TokenQuery(evmcli, tokenAddress, tokenName)
	}
	tendermintCli, ok := chain.(*TendermintClient)
	if !ok {
		panic(fmt.Sprintf("invalid token,chain:%s,tokenName:%s,tokenAddress:%s", chain.ChainName(), tokenName, tokenAddress))
	}
	return NewTendermintTokenQuery(tendermintCli, tokenAddress, tokenName, decimals)
}


type tokenErc20Query struct {
	*Evm
	tokenName  string
	addressHex string
	contract   *ethbind.BoundContract
	decimals   uint8
}



func NewErc20TokenQuery(eth *Evm, addressHex, tokenName string) TokenQuery {
	token := getTokenContract()
	address := common.HexToAddress(addressHex)
	tQuery := &tokenErc20Query{
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

func (t *tokenErc20Query) Address() string {
	return t.addressHex
}

func (t *tokenErc20Query) TokenName() string {
	return t.tokenName
}

func (t *tokenErc20Query) GetErc20TotalSupply(blockNumber *big.Int) (TokenAmount, error) {
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

func (t *tokenErc20Query) QueryTotalBurnToken() (TokenAmount, error) {
	return TokenAmount{}, nil
}

func (t *tokenErc20Query) getErc20Decimals() (uint8, error) {
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
func (t *tokenErc20Query) GetDecimals() (uint8, error) {
	if t.NativeToken() == t.TokenName() {
		return t.Evm.GetNativeDecimal()
	}
	return t.getErc20Decimals()
}

func (t *tokenErc20Query) GetTokenBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	var balanceQuery func(addressHex string, blockNumber *big.Int) (TokenAmount, error)
	if t.addressHex == DefaultNativeToken {
		balanceQuery = t.Evm.GetBalance
	} else {
		balanceQuery = t.GetErc20Balance
	}
	return balanceQuery(addressHex, blockNumber)
}

func (t *tokenErc20Query) GetErc20Balance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
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

func (t *tokenErc20Query) GetTransferBalance(blockNumber *big.Int) (TokenAmount, error) {
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

func (t *tokenErc20Query) GetInTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
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

func (t *tokenErc20Query) GetOutTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
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

type tendermintTokenQuery struct {
	*TendermintClient
	tokenName string
	denom     string
	decimals  uint8
}

func NewTendermintTokenQuery(tendermintCli *TendermintClient, denom, tokenName string, decimals uint8) TokenQuery {
	tQuery := tendermintTokenQuery{
		TendermintClient: tendermintCli,
		tokenName:        tokenName,
		denom:            denom,
		decimals:         decimals,
	}
	return &tQuery
}

func (t *tendermintTokenQuery) Address() string {
	return t.denom
}

func (t *tendermintTokenQuery) TokenName() string {
	return t.tokenName
}

func (t *tendermintTokenQuery) GetErc20TotalSupply(blockNumber *big.Int) (TokenAmount, error) {
	// TODO
	return TokenAmount{}, fmt.Errorf("invalid totalSupply type")
}

func (t *tendermintTokenQuery) getErc20Decimals() (uint8, error) {
	// TODO
	return 0, fmt.Errorf("invalid decimals type")
}
func (t *tendermintTokenQuery) GetDecimals() (uint8, error) {
	// TODO decimals
	return t.decimals, nil
}

func (t *tendermintTokenQuery) GetTokenBalance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	return t.TendermintClient.GetTokenBalance(t.denom, addressHex)
}

func (t *tendermintTokenQuery) GetErc20Balance(addressHex string, blockNumber *big.Int) (TokenAmount, error) {
	return TokenAmount{}, fmt.Errorf("invalid balance type")
}

func (t *tendermintTokenQuery) GetTransferBalance(blockNumber *big.Int) (TokenAmount, error) {
	return TokenAmount{}, nil
}

func (t *tendermintTokenQuery) GetInTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
	return TokenAmount{}, nil
}

func (t *tendermintTokenQuery) GetOutTokenAmount(addressHex string, srcChain string, blockNumber *big.Int) (TokenAmount, error) {
	return TokenAmount{}, nil
}
