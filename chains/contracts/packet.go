// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// HeightData is an auto generated low-level Go binding around an user-defined struct.
type HeightData struct {
	RevisionNumber uint64
	RevisionHeight uint64
}

// PacketTypesFee is an auto generated low-level Go binding around an user-defined struct.
type PacketTypesFee struct {
	TokenAddress common.Address
	Amount       *big.Int
}

// PacketTypesPacket is an auto generated low-level Go binding around an user-defined struct.
type PacketTypesPacket struct {
	SrcChain        string
	DstChain        string
	Sequence        uint64
	Sender          string
	TransferData    []byte
	CallData        []byte
	CallbackAddress string
	FeeOption       uint64
}

// PacketMetaData contains all meta data concerning the Packet contract.
var PacketMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structPacketTypes.Packet\",\"name\":\"packet\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ack\",\"type\":\"bytes\"}],\"name\":\"AckPacket\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structPacketTypes.Packet\",\"name\":\"packet\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"ack\",\"type\":\"bytes\"}],\"name\":\"AckWritten\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"indexed\":false,\"internalType\":\"structPacketTypes.Packet\",\"name\":\"packet\",\"type\":\"tuple\"}],\"name\":\"PacketReceived\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"packetBytes\",\"type\":\"bytes\"}],\"name\":\"PacketSent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"FEE_MANAGER\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PAUSER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAYER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accessManager\",\"outputs\":[{\"internalType\":\"contractIAccessManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"ackStatus\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"packetBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"acknowledgement\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proofAcked\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"revision_number\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"revision_height\",\"type\":\"uint64\"}],\"internalType\":\"structHeight.Data\",\"name\":\"height\",\"type\":\"tuple\"}],\"name\":\"acknowledgePacket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"acks\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"code\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"relayer\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"addPacketFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"claim2HopsFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"clientManager\",\"outputs\":[{\"internalType\":\"contractIClientManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"commitments\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"endpoint\",\"outputs\":[{\"internalType\":\"contractIEndpoint\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"execute\",\"outputs\":[{\"internalType\":\"contractIExecute\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"fee2HopsRemaining\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"}],\"name\":\"getAckStatus\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestPacket\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"internalType\":\"structPacketTypes.Packet\",\"name\":\"packet\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"}],\"name\":\"getNextSequenceSend\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_endpointContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_executeContract\",\"type\":\"address\"}],\"name\":\"initEndpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_chainName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_relayChainName\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_clientManagerContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_accessManagerContract\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestPacket\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"packetFees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"receipts\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"packetBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint64\",\"name\":\"revision_number\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"revision_height\",\"type\":\"uint64\"}],\"internalType\":\"structHeight.Data\",\"name\":\"height\",\"type\":\"tuple\"}],\"name\":\"recvPacket\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"relayChainName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"srcChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"sequence\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"sender\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"transferData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"callData\",\"type\":\"bytes\"},{\"internalType\":\"string\",\"name\":\"callbackAddress\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"feeOption\",\"type\":\"uint64\"}],\"internalType\":\"structPacketTypes.Packet\",\"name\":\"packet\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structPacketTypes.Fee\",\"name\":\"fee\",\"type\":\"tuple\"}],\"name\":\"sendPacket\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"sequences\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PacketABI is the input ABI used to generate the binding from.
// Deprecated: Use PacketMetaData.ABI instead.
var PacketABI = PacketMetaData.ABI

// Packet is an auto generated Go binding around an Ethereum contract.
type Packet struct {
	PacketCaller     // Read-only binding to the contract
	PacketTransactor // Write-only binding to the contract
	PacketFilterer   // Log filterer for contract events
}

// PacketCaller is an auto generated read-only Go binding around an Ethereum contract.
type PacketCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PacketTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PacketTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PacketFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PacketFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PacketSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PacketSession struct {
	Contract     *Packet           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PacketCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PacketCallerSession struct {
	Contract *PacketCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PacketTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PacketTransactorSession struct {
	Contract     *PacketTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PacketRaw is an auto generated low-level Go binding around an Ethereum contract.
type PacketRaw struct {
	Contract *Packet // Generic contract binding to access the raw methods on
}

// PacketCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PacketCallerRaw struct {
	Contract *PacketCaller // Generic read-only contract binding to access the raw methods on
}

// PacketTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PacketTransactorRaw struct {
	Contract *PacketTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPacket creates a new instance of Packet, bound to a specific deployed contract.
func NewPacket(address common.Address, backend bind.ContractBackend) (*Packet, error) {
	contract, err := bindPacket(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Packet{PacketCaller: PacketCaller{contract: contract}, PacketTransactor: PacketTransactor{contract: contract}, PacketFilterer: PacketFilterer{contract: contract}}, nil
}

// NewPacketCaller creates a new read-only instance of Packet, bound to a specific deployed contract.
func NewPacketCaller(address common.Address, caller bind.ContractCaller) (*PacketCaller, error) {
	contract, err := bindPacket(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PacketCaller{contract: contract}, nil
}

// NewPacketTransactor creates a new write-only instance of Packet, bound to a specific deployed contract.
func NewPacketTransactor(address common.Address, transactor bind.ContractTransactor) (*PacketTransactor, error) {
	contract, err := bindPacket(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PacketTransactor{contract: contract}, nil
}

// NewPacketFilterer creates a new log filterer instance of Packet, bound to a specific deployed contract.
func NewPacketFilterer(address common.Address, filterer bind.ContractFilterer) (*PacketFilterer, error) {
	contract, err := bindPacket(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PacketFilterer{contract: contract}, nil
}

// bindPacket binds a generic wrapper to an already deployed contract.
func bindPacket(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PacketABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Packet *PacketRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Packet.Contract.PacketCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Packet *PacketRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Packet.Contract.PacketTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Packet *PacketRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Packet.Contract.PacketTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Packet *PacketCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Packet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Packet *PacketTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Packet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Packet *PacketTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Packet.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Packet *PacketCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Packet *PacketSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Packet.Contract.DEFAULTADMINROLE(&_Packet.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Packet *PacketCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Packet.Contract.DEFAULTADMINROLE(&_Packet.CallOpts)
}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(bytes32)
func (_Packet *PacketCaller) FEEMANAGER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "FEE_MANAGER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(bytes32)
func (_Packet *PacketSession) FEEMANAGER() ([32]byte, error) {
	return _Packet.Contract.FEEMANAGER(&_Packet.CallOpts)
}

// FEEMANAGER is a free data retrieval call binding the contract method 0xea26266c.
//
// Solidity: function FEE_MANAGER() view returns(bytes32)
func (_Packet *PacketCallerSession) FEEMANAGER() ([32]byte, error) {
	return _Packet.Contract.FEEMANAGER(&_Packet.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Packet *PacketCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Packet *PacketSession) PAUSERROLE() ([32]byte, error) {
	return _Packet.Contract.PAUSERROLE(&_Packet.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_Packet *PacketCallerSession) PAUSERROLE() ([32]byte, error) {
	return _Packet.Contract.PAUSERROLE(&_Packet.CallOpts)
}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Packet *PacketCaller) RELAYERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "RELAYER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Packet *PacketSession) RELAYERROLE() ([32]byte, error) {
	return _Packet.Contract.RELAYERROLE(&_Packet.CallOpts)
}

// RELAYERROLE is a free data retrieval call binding the contract method 0x926d7d7f.
//
// Solidity: function RELAYER_ROLE() view returns(bytes32)
func (_Packet *PacketCallerSession) RELAYERROLE() ([32]byte, error) {
	return _Packet.Contract.RELAYERROLE(&_Packet.CallOpts)
}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_Packet *PacketCaller) AccessManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "accessManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_Packet *PacketSession) AccessManager() (common.Address, error) {
	return _Packet.Contract.AccessManager(&_Packet.CallOpts)
}

// AccessManager is a free data retrieval call binding the contract method 0xfdcb6068.
//
// Solidity: function accessManager() view returns(address)
func (_Packet *PacketCallerSession) AccessManager() (common.Address, error) {
	return _Packet.Contract.AccessManager(&_Packet.CallOpts)
}

// AckStatus is a free data retrieval call binding the contract method 0xf6a6b4f2.
//
// Solidity: function ackStatus(bytes ) view returns(uint8)
func (_Packet *PacketCaller) AckStatus(opts *bind.CallOpts, arg0 []byte) (uint8, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "ackStatus", arg0)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// AckStatus is a free data retrieval call binding the contract method 0xf6a6b4f2.
//
// Solidity: function ackStatus(bytes ) view returns(uint8)
func (_Packet *PacketSession) AckStatus(arg0 []byte) (uint8, error) {
	return _Packet.Contract.AckStatus(&_Packet.CallOpts, arg0)
}

// AckStatus is a free data retrieval call binding the contract method 0xf6a6b4f2.
//
// Solidity: function ackStatus(bytes ) view returns(uint8)
func (_Packet *PacketCallerSession) AckStatus(arg0 []byte) (uint8, error) {
	return _Packet.Contract.AckStatus(&_Packet.CallOpts, arg0)
}

// Acks is a free data retrieval call binding the contract method 0xfcd5670f.
//
// Solidity: function acks(bytes ) view returns(uint64 code, bytes result, string message, string relayer, uint64 feeOption)
func (_Packet *PacketCaller) Acks(opts *bind.CallOpts, arg0 []byte) (struct {
	Code      uint64
	Result    []byte
	Message   string
	Relayer   string
	FeeOption uint64
}, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "acks", arg0)

	outstruct := new(struct {
		Code      uint64
		Result    []byte
		Message   string
		Relayer   string
		FeeOption uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Code = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Result = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	outstruct.Message = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Relayer = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.FeeOption = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// Acks is a free data retrieval call binding the contract method 0xfcd5670f.
//
// Solidity: function acks(bytes ) view returns(uint64 code, bytes result, string message, string relayer, uint64 feeOption)
func (_Packet *PacketSession) Acks(arg0 []byte) (struct {
	Code      uint64
	Result    []byte
	Message   string
	Relayer   string
	FeeOption uint64
}, error) {
	return _Packet.Contract.Acks(&_Packet.CallOpts, arg0)
}

// Acks is a free data retrieval call binding the contract method 0xfcd5670f.
//
// Solidity: function acks(bytes ) view returns(uint64 code, bytes result, string message, string relayer, uint64 feeOption)
func (_Packet *PacketCallerSession) Acks(arg0 []byte) (struct {
	Code      uint64
	Result    []byte
	Message   string
	Relayer   string
	FeeOption uint64
}, error) {
	return _Packet.Contract.Acks(&_Packet.CallOpts, arg0)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_Packet *PacketCaller) ChainName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "chainName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_Packet *PacketSession) ChainName() (string, error) {
	return _Packet.Contract.ChainName(&_Packet.CallOpts)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() view returns(string)
func (_Packet *PacketCallerSession) ChainName() (string, error) {
	return _Packet.Contract.ChainName(&_Packet.CallOpts)
}

// ClientManager is a free data retrieval call binding the contract method 0x79e8be1d.
//
// Solidity: function clientManager() view returns(address)
func (_Packet *PacketCaller) ClientManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "clientManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ClientManager is a free data retrieval call binding the contract method 0x79e8be1d.
//
// Solidity: function clientManager() view returns(address)
func (_Packet *PacketSession) ClientManager() (common.Address, error) {
	return _Packet.Contract.ClientManager(&_Packet.CallOpts)
}

// ClientManager is a free data retrieval call binding the contract method 0x79e8be1d.
//
// Solidity: function clientManager() view returns(address)
func (_Packet *PacketCallerSession) ClientManager() (common.Address, error) {
	return _Packet.Contract.ClientManager(&_Packet.CallOpts)
}

// Commitments is a free data retrieval call binding the contract method 0x7912b8e6.
//
// Solidity: function commitments(bytes ) view returns(bytes32)
func (_Packet *PacketCaller) Commitments(opts *bind.CallOpts, arg0 []byte) ([32]byte, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "commitments", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Commitments is a free data retrieval call binding the contract method 0x7912b8e6.
//
// Solidity: function commitments(bytes ) view returns(bytes32)
func (_Packet *PacketSession) Commitments(arg0 []byte) ([32]byte, error) {
	return _Packet.Contract.Commitments(&_Packet.CallOpts, arg0)
}

// Commitments is a free data retrieval call binding the contract method 0x7912b8e6.
//
// Solidity: function commitments(bytes ) view returns(bytes32)
func (_Packet *PacketCallerSession) Commitments(arg0 []byte) ([32]byte, error) {
	return _Packet.Contract.Commitments(&_Packet.CallOpts, arg0)
}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_Packet *PacketCaller) Endpoint(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "endpoint")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_Packet *PacketSession) Endpoint() (common.Address, error) {
	return _Packet.Contract.Endpoint(&_Packet.CallOpts)
}

// Endpoint is a free data retrieval call binding the contract method 0x5e280f11.
//
// Solidity: function endpoint() view returns(address)
func (_Packet *PacketCallerSession) Endpoint() (common.Address, error) {
	return _Packet.Contract.Endpoint(&_Packet.CallOpts)
}

// Execute is a free data retrieval call binding the contract method 0x61461954.
//
// Solidity: function execute() view returns(address)
func (_Packet *PacketCaller) Execute(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "execute")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Execute is a free data retrieval call binding the contract method 0x61461954.
//
// Solidity: function execute() view returns(address)
func (_Packet *PacketSession) Execute() (common.Address, error) {
	return _Packet.Contract.Execute(&_Packet.CallOpts)
}

// Execute is a free data retrieval call binding the contract method 0x61461954.
//
// Solidity: function execute() view returns(address)
func (_Packet *PacketCallerSession) Execute() (common.Address, error) {
	return _Packet.Contract.Execute(&_Packet.CallOpts)
}

// Fee2HopsRemaining is a free data retrieval call binding the contract method 0x2b275e6c.
//
// Solidity: function fee2HopsRemaining(address ) view returns(uint256)
func (_Packet *PacketCaller) Fee2HopsRemaining(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "fee2HopsRemaining", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee2HopsRemaining is a free data retrieval call binding the contract method 0x2b275e6c.
//
// Solidity: function fee2HopsRemaining(address ) view returns(uint256)
func (_Packet *PacketSession) Fee2HopsRemaining(arg0 common.Address) (*big.Int, error) {
	return _Packet.Contract.Fee2HopsRemaining(&_Packet.CallOpts, arg0)
}

// Fee2HopsRemaining is a free data retrieval call binding the contract method 0x2b275e6c.
//
// Solidity: function fee2HopsRemaining(address ) view returns(uint256)
func (_Packet *PacketCallerSession) Fee2HopsRemaining(arg0 common.Address) (*big.Int, error) {
	return _Packet.Contract.Fee2HopsRemaining(&_Packet.CallOpts, arg0)
}

// GetAckStatus is a free data retrieval call binding the contract method 0x10ca494a.
//
// Solidity: function getAckStatus(string dstChain, uint64 sequence) view returns(uint8)
func (_Packet *PacketCaller) GetAckStatus(opts *bind.CallOpts, dstChain string, sequence uint64) (uint8, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "getAckStatus", dstChain, sequence)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetAckStatus is a free data retrieval call binding the contract method 0x10ca494a.
//
// Solidity: function getAckStatus(string dstChain, uint64 sequence) view returns(uint8)
func (_Packet *PacketSession) GetAckStatus(dstChain string, sequence uint64) (uint8, error) {
	return _Packet.Contract.GetAckStatus(&_Packet.CallOpts, dstChain, sequence)
}

// GetAckStatus is a free data retrieval call binding the contract method 0x10ca494a.
//
// Solidity: function getAckStatus(string dstChain, uint64 sequence) view returns(uint8)
func (_Packet *PacketCallerSession) GetAckStatus(dstChain string, sequence uint64) (uint8, error) {
	return _Packet.Contract.GetAckStatus(&_Packet.CallOpts, dstChain, sequence)
}

// GetLatestPacket is a free data retrieval call binding the contract method 0x775acbce.
//
// Solidity: function getLatestPacket() view returns((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketCaller) GetLatestPacket(opts *bind.CallOpts) (PacketTypesPacket, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "getLatestPacket")

	if err != nil {
		return *new(PacketTypesPacket), err
	}

	out0 := *abi.ConvertType(out[0], new(PacketTypesPacket)).(*PacketTypesPacket)

	return out0, err

}

// GetLatestPacket is a free data retrieval call binding the contract method 0x775acbce.
//
// Solidity: function getLatestPacket() view returns((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketSession) GetLatestPacket() (PacketTypesPacket, error) {
	return _Packet.Contract.GetLatestPacket(&_Packet.CallOpts)
}

// GetLatestPacket is a free data retrieval call binding the contract method 0x775acbce.
//
// Solidity: function getLatestPacket() view returns((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketCallerSession) GetLatestPacket() (PacketTypesPacket, error) {
	return _Packet.Contract.GetLatestPacket(&_Packet.CallOpts)
}

// GetNextSequenceSend is a free data retrieval call binding the contract method 0x39151dc6.
//
// Solidity: function getNextSequenceSend(string dstChain) view returns(uint64)
func (_Packet *PacketCaller) GetNextSequenceSend(opts *bind.CallOpts, dstChain string) (uint64, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "getNextSequenceSend", dstChain)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetNextSequenceSend is a free data retrieval call binding the contract method 0x39151dc6.
//
// Solidity: function getNextSequenceSend(string dstChain) view returns(uint64)
func (_Packet *PacketSession) GetNextSequenceSend(dstChain string) (uint64, error) {
	return _Packet.Contract.GetNextSequenceSend(&_Packet.CallOpts, dstChain)
}

// GetNextSequenceSend is a free data retrieval call binding the contract method 0x39151dc6.
//
// Solidity: function getNextSequenceSend(string dstChain) view returns(uint64)
func (_Packet *PacketCallerSession) GetNextSequenceSend(dstChain string) (uint64, error) {
	return _Packet.Contract.GetNextSequenceSend(&_Packet.CallOpts, dstChain)
}

// LatestPacket is a free data retrieval call binding the contract method 0xe55027c6.
//
// Solidity: function latestPacket() view returns(string srcChain, string dstChain, uint64 sequence, string sender, bytes transferData, bytes callData, string callbackAddress, uint64 feeOption)
func (_Packet *PacketCaller) LatestPacket(opts *bind.CallOpts) (struct {
	SrcChain        string
	DstChain        string
	Sequence        uint64
	Sender          string
	TransferData    []byte
	CallData        []byte
	CallbackAddress string
	FeeOption       uint64
}, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "latestPacket")

	outstruct := new(struct {
		SrcChain        string
		DstChain        string
		Sequence        uint64
		Sender          string
		TransferData    []byte
		CallData        []byte
		CallbackAddress string
		FeeOption       uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SrcChain = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.DstChain = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Sequence = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Sender = *abi.ConvertType(out[3], new(string)).(*string)
	outstruct.TransferData = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.CallData = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	outstruct.CallbackAddress = *abi.ConvertType(out[6], new(string)).(*string)
	outstruct.FeeOption = *abi.ConvertType(out[7], new(uint64)).(*uint64)

	return *outstruct, err

}

// LatestPacket is a free data retrieval call binding the contract method 0xe55027c6.
//
// Solidity: function latestPacket() view returns(string srcChain, string dstChain, uint64 sequence, string sender, bytes transferData, bytes callData, string callbackAddress, uint64 feeOption)
func (_Packet *PacketSession) LatestPacket() (struct {
	SrcChain        string
	DstChain        string
	Sequence        uint64
	Sender          string
	TransferData    []byte
	CallData        []byte
	CallbackAddress string
	FeeOption       uint64
}, error) {
	return _Packet.Contract.LatestPacket(&_Packet.CallOpts)
}

// LatestPacket is a free data retrieval call binding the contract method 0xe55027c6.
//
// Solidity: function latestPacket() view returns(string srcChain, string dstChain, uint64 sequence, string sender, bytes transferData, bytes callData, string callbackAddress, uint64 feeOption)
func (_Packet *PacketCallerSession) LatestPacket() (struct {
	SrcChain        string
	DstChain        string
	Sequence        uint64
	Sender          string
	TransferData    []byte
	CallData        []byte
	CallbackAddress string
	FeeOption       uint64
}, error) {
	return _Packet.Contract.LatestPacket(&_Packet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Packet *PacketCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Packet *PacketSession) Owner() (common.Address, error) {
	return _Packet.Contract.Owner(&_Packet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Packet *PacketCallerSession) Owner() (common.Address, error) {
	return _Packet.Contract.Owner(&_Packet.CallOpts)
}

// PacketFees is a free data retrieval call binding the contract method 0x374f7f03.
//
// Solidity: function packetFees(bytes ) view returns(address tokenAddress, uint256 amount)
func (_Packet *PacketCaller) PacketFees(opts *bind.CallOpts, arg0 []byte) (struct {
	TokenAddress common.Address
	Amount       *big.Int
}, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "packetFees", arg0)

	outstruct := new(struct {
		TokenAddress common.Address
		Amount       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TokenAddress = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PacketFees is a free data retrieval call binding the contract method 0x374f7f03.
//
// Solidity: function packetFees(bytes ) view returns(address tokenAddress, uint256 amount)
func (_Packet *PacketSession) PacketFees(arg0 []byte) (struct {
	TokenAddress common.Address
	Amount       *big.Int
}, error) {
	return _Packet.Contract.PacketFees(&_Packet.CallOpts, arg0)
}

// PacketFees is a free data retrieval call binding the contract method 0x374f7f03.
//
// Solidity: function packetFees(bytes ) view returns(address tokenAddress, uint256 amount)
func (_Packet *PacketCallerSession) PacketFees(arg0 []byte) (struct {
	TokenAddress common.Address
	Amount       *big.Int
}, error) {
	return _Packet.Contract.PacketFees(&_Packet.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Packet *PacketCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Packet *PacketSession) Paused() (bool, error) {
	return _Packet.Contract.Paused(&_Packet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Packet *PacketCallerSession) Paused() (bool, error) {
	return _Packet.Contract.Paused(&_Packet.CallOpts)
}

// Receipts is a free data retrieval call binding the contract method 0xa6992b83.
//
// Solidity: function receipts(bytes ) view returns(bool)
func (_Packet *PacketCaller) Receipts(opts *bind.CallOpts, arg0 []byte) (bool, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "receipts", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Receipts is a free data retrieval call binding the contract method 0xa6992b83.
//
// Solidity: function receipts(bytes ) view returns(bool)
func (_Packet *PacketSession) Receipts(arg0 []byte) (bool, error) {
	return _Packet.Contract.Receipts(&_Packet.CallOpts, arg0)
}

// Receipts is a free data retrieval call binding the contract method 0xa6992b83.
//
// Solidity: function receipts(bytes ) view returns(bool)
func (_Packet *PacketCallerSession) Receipts(arg0 []byte) (bool, error) {
	return _Packet.Contract.Receipts(&_Packet.CallOpts, arg0)
}

// RelayChainName is a free data retrieval call binding the contract method 0xe7f4bc1e.
//
// Solidity: function relayChainName() view returns(string)
func (_Packet *PacketCaller) RelayChainName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "relayChainName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// RelayChainName is a free data retrieval call binding the contract method 0xe7f4bc1e.
//
// Solidity: function relayChainName() view returns(string)
func (_Packet *PacketSession) RelayChainName() (string, error) {
	return _Packet.Contract.RelayChainName(&_Packet.CallOpts)
}

// RelayChainName is a free data retrieval call binding the contract method 0xe7f4bc1e.
//
// Solidity: function relayChainName() view returns(string)
func (_Packet *PacketCallerSession) RelayChainName() (string, error) {
	return _Packet.Contract.RelayChainName(&_Packet.CallOpts)
}

// Sequences is a free data retrieval call binding the contract method 0xeeebb020.
//
// Solidity: function sequences(bytes ) view returns(uint64)
func (_Packet *PacketCaller) Sequences(opts *bind.CallOpts, arg0 []byte) (uint64, error) {
	var out []interface{}
	err := _Packet.contract.Call(opts, &out, "sequences", arg0)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Sequences is a free data retrieval call binding the contract method 0xeeebb020.
//
// Solidity: function sequences(bytes ) view returns(uint64)
func (_Packet *PacketSession) Sequences(arg0 []byte) (uint64, error) {
	return _Packet.Contract.Sequences(&_Packet.CallOpts, arg0)
}

// Sequences is a free data retrieval call binding the contract method 0xeeebb020.
//
// Solidity: function sequences(bytes ) view returns(uint64)
func (_Packet *PacketCallerSession) Sequences(arg0 []byte) (uint64, error) {
	return _Packet.Contract.Sequences(&_Packet.CallOpts, arg0)
}

// AcknowledgePacket is a paid mutator transaction binding the contract method 0x9140d507.
//
// Solidity: function acknowledgePacket(bytes packetBytes, bytes acknowledgement, bytes proofAcked, (uint64,uint64) height) returns()
func (_Packet *PacketTransactor) AcknowledgePacket(opts *bind.TransactOpts, packetBytes []byte, acknowledgement []byte, proofAcked []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "acknowledgePacket", packetBytes, acknowledgement, proofAcked, height)
}

// AcknowledgePacket is a paid mutator transaction binding the contract method 0x9140d507.
//
// Solidity: function acknowledgePacket(bytes packetBytes, bytes acknowledgement, bytes proofAcked, (uint64,uint64) height) returns()
func (_Packet *PacketSession) AcknowledgePacket(packetBytes []byte, acknowledgement []byte, proofAcked []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.Contract.AcknowledgePacket(&_Packet.TransactOpts, packetBytes, acknowledgement, proofAcked, height)
}

// AcknowledgePacket is a paid mutator transaction binding the contract method 0x9140d507.
//
// Solidity: function acknowledgePacket(bytes packetBytes, bytes acknowledgement, bytes proofAcked, (uint64,uint64) height) returns()
func (_Packet *PacketTransactorSession) AcknowledgePacket(packetBytes []byte, acknowledgement []byte, proofAcked []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.Contract.AcknowledgePacket(&_Packet.TransactOpts, packetBytes, acknowledgement, proofAcked, height)
}

// AddPacketFee is a paid mutator transaction binding the contract method 0xc9dfad20.
//
// Solidity: function addPacketFee(string dstChain, uint64 sequence, uint256 amount) payable returns()
func (_Packet *PacketTransactor) AddPacketFee(opts *bind.TransactOpts, dstChain string, sequence uint64, amount *big.Int) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "addPacketFee", dstChain, sequence, amount)
}

// AddPacketFee is a paid mutator transaction binding the contract method 0xc9dfad20.
//
// Solidity: function addPacketFee(string dstChain, uint64 sequence, uint256 amount) payable returns()
func (_Packet *PacketSession) AddPacketFee(dstChain string, sequence uint64, amount *big.Int) (*types.Transaction, error) {
	return _Packet.Contract.AddPacketFee(&_Packet.TransactOpts, dstChain, sequence, amount)
}

// AddPacketFee is a paid mutator transaction binding the contract method 0xc9dfad20.
//
// Solidity: function addPacketFee(string dstChain, uint64 sequence, uint256 amount) payable returns()
func (_Packet *PacketTransactorSession) AddPacketFee(dstChain string, sequence uint64, amount *big.Int) (*types.Transaction, error) {
	return _Packet.Contract.AddPacketFee(&_Packet.TransactOpts, dstChain, sequence, amount)
}

// Claim2HopsFee is a paid mutator transaction binding the contract method 0x730f9c36.
//
// Solidity: function claim2HopsFee(address[] tokens, address receiver) returns()
func (_Packet *PacketTransactor) Claim2HopsFee(opts *bind.TransactOpts, tokens []common.Address, receiver common.Address) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "claim2HopsFee", tokens, receiver)
}

// Claim2HopsFee is a paid mutator transaction binding the contract method 0x730f9c36.
//
// Solidity: function claim2HopsFee(address[] tokens, address receiver) returns()
func (_Packet *PacketSession) Claim2HopsFee(tokens []common.Address, receiver common.Address) (*types.Transaction, error) {
	return _Packet.Contract.Claim2HopsFee(&_Packet.TransactOpts, tokens, receiver)
}

// Claim2HopsFee is a paid mutator transaction binding the contract method 0x730f9c36.
//
// Solidity: function claim2HopsFee(address[] tokens, address receiver) returns()
func (_Packet *PacketTransactorSession) Claim2HopsFee(tokens []common.Address, receiver common.Address) (*types.Transaction, error) {
	return _Packet.Contract.Claim2HopsFee(&_Packet.TransactOpts, tokens, receiver)
}

// InitEndpoint is a paid mutator transaction binding the contract method 0x46d6ca73.
//
// Solidity: function initEndpoint(address _endpointContract, address _executeContract) returns()
func (_Packet *PacketTransactor) InitEndpoint(opts *bind.TransactOpts, _endpointContract common.Address, _executeContract common.Address) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "initEndpoint", _endpointContract, _executeContract)
}

// InitEndpoint is a paid mutator transaction binding the contract method 0x46d6ca73.
//
// Solidity: function initEndpoint(address _endpointContract, address _executeContract) returns()
func (_Packet *PacketSession) InitEndpoint(_endpointContract common.Address, _executeContract common.Address) (*types.Transaction, error) {
	return _Packet.Contract.InitEndpoint(&_Packet.TransactOpts, _endpointContract, _executeContract)
}

// InitEndpoint is a paid mutator transaction binding the contract method 0x46d6ca73.
//
// Solidity: function initEndpoint(address _endpointContract, address _executeContract) returns()
func (_Packet *PacketTransactorSession) InitEndpoint(_endpointContract common.Address, _executeContract common.Address) (*types.Transaction, error) {
	return _Packet.Contract.InitEndpoint(&_Packet.TransactOpts, _endpointContract, _executeContract)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string _chainName, string _relayChainName, address _clientManagerContract, address _accessManagerContract) returns()
func (_Packet *PacketTransactor) Initialize(opts *bind.TransactOpts, _chainName string, _relayChainName string, _clientManagerContract common.Address, _accessManagerContract common.Address) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "initialize", _chainName, _relayChainName, _clientManagerContract, _accessManagerContract)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string _chainName, string _relayChainName, address _clientManagerContract, address _accessManagerContract) returns()
func (_Packet *PacketSession) Initialize(_chainName string, _relayChainName string, _clientManagerContract common.Address, _accessManagerContract common.Address) (*types.Transaction, error) {
	return _Packet.Contract.Initialize(&_Packet.TransactOpts, _chainName, _relayChainName, _clientManagerContract, _accessManagerContract)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string _chainName, string _relayChainName, address _clientManagerContract, address _accessManagerContract) returns()
func (_Packet *PacketTransactorSession) Initialize(_chainName string, _relayChainName string, _clientManagerContract common.Address, _accessManagerContract common.Address) (*types.Transaction, error) {
	return _Packet.Contract.Initialize(&_Packet.TransactOpts, _chainName, _relayChainName, _clientManagerContract, _accessManagerContract)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Packet *PacketTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Packet *PacketSession) Pause() (*types.Transaction, error) {
	return _Packet.Contract.Pause(&_Packet.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Packet *PacketTransactorSession) Pause() (*types.Transaction, error) {
	return _Packet.Contract.Pause(&_Packet.TransactOpts)
}

// RecvPacket is a paid mutator transaction binding the contract method 0xea3449b6.
//
// Solidity: function recvPacket(bytes packetBytes, bytes proof, (uint64,uint64) height) returns()
func (_Packet *PacketTransactor) RecvPacket(opts *bind.TransactOpts, packetBytes []byte, proof []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "recvPacket", packetBytes, proof, height)
}

// RecvPacket is a paid mutator transaction binding the contract method 0xea3449b6.
//
// Solidity: function recvPacket(bytes packetBytes, bytes proof, (uint64,uint64) height) returns()
func (_Packet *PacketSession) RecvPacket(packetBytes []byte, proof []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.Contract.RecvPacket(&_Packet.TransactOpts, packetBytes, proof, height)
}

// RecvPacket is a paid mutator transaction binding the contract method 0xea3449b6.
//
// Solidity: function recvPacket(bytes packetBytes, bytes proof, (uint64,uint64) height) returns()
func (_Packet *PacketTransactorSession) RecvPacket(packetBytes []byte, proof []byte, height HeightData) (*types.Transaction, error) {
	return _Packet.Contract.RecvPacket(&_Packet.TransactOpts, packetBytes, proof, height)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Packet *PacketTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Packet *PacketSession) RenounceOwnership() (*types.Transaction, error) {
	return _Packet.Contract.RenounceOwnership(&_Packet.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Packet *PacketTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Packet.Contract.RenounceOwnership(&_Packet.TransactOpts)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4072d026.
//
// Solidity: function sendPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, (address,uint256) fee) payable returns()
func (_Packet *PacketTransactor) SendPacket(opts *bind.TransactOpts, packet PacketTypesPacket, fee PacketTypesFee) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "sendPacket", packet, fee)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4072d026.
//
// Solidity: function sendPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, (address,uint256) fee) payable returns()
func (_Packet *PacketSession) SendPacket(packet PacketTypesPacket, fee PacketTypesFee) (*types.Transaction, error) {
	return _Packet.Contract.SendPacket(&_Packet.TransactOpts, packet, fee)
}

// SendPacket is a paid mutator transaction binding the contract method 0x4072d026.
//
// Solidity: function sendPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, (address,uint256) fee) payable returns()
func (_Packet *PacketTransactorSession) SendPacket(packet PacketTypesPacket, fee PacketTypesFee) (*types.Transaction, error) {
	return _Packet.Contract.SendPacket(&_Packet.TransactOpts, packet, fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Packet *PacketTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Packet *PacketSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Packet.Contract.TransferOwnership(&_Packet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Packet *PacketTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Packet.Contract.TransferOwnership(&_Packet.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Packet *PacketTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Packet.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Packet *PacketSession) Unpause() (*types.Transaction, error) {
	return _Packet.Contract.Unpause(&_Packet.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Packet *PacketTransactorSession) Unpause() (*types.Transaction, error) {
	return _Packet.Contract.Unpause(&_Packet.TransactOpts)
}

// PacketAckPacketIterator is returned from FilterAckPacket and is used to iterate over the raw logs and unpacked data for AckPacket events raised by the Packet contract.
type PacketAckPacketIterator struct {
	Event *PacketAckPacket // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketAckPacketIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketAckPacket)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketAckPacket)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketAckPacketIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketAckPacketIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketAckPacket represents a AckPacket event raised by the Packet contract.
type PacketAckPacket struct {
	Packet PacketTypesPacket
	Ack    []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAckPacket is a free log retrieval operation binding the contract event 0x0ff8d9c09209e78f3ccc3861e1d2ea9b5bb54872c2bea768a79e820fe169995f.
//
// Solidity: event AckPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) FilterAckPacket(opts *bind.FilterOpts) (*PacketAckPacketIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "AckPacket")
	if err != nil {
		return nil, err
	}
	return &PacketAckPacketIterator{contract: _Packet.contract, event: "AckPacket", logs: logs, sub: sub}, nil
}

// WatchAckPacket is a free log subscription operation binding the contract event 0x0ff8d9c09209e78f3ccc3861e1d2ea9b5bb54872c2bea768a79e820fe169995f.
//
// Solidity: event AckPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) WatchAckPacket(opts *bind.WatchOpts, sink chan<- *PacketAckPacket) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "AckPacket")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketAckPacket)
				if err := _Packet.contract.UnpackLog(event, "AckPacket", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAckPacket is a log parse operation binding the contract event 0x0ff8d9c09209e78f3ccc3861e1d2ea9b5bb54872c2bea768a79e820fe169995f.
//
// Solidity: event AckPacket((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) ParseAckPacket(log types.Log) (*PacketAckPacket, error) {
	event := new(PacketAckPacket)
	if err := _Packet.contract.UnpackLog(event, "AckPacket", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketAckWrittenIterator is returned from FilterAckWritten and is used to iterate over the raw logs and unpacked data for AckWritten events raised by the Packet contract.
type PacketAckWrittenIterator struct {
	Event *PacketAckWritten // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketAckWrittenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketAckWritten)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketAckWritten)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketAckWrittenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketAckWrittenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketAckWritten represents a AckWritten event raised by the Packet contract.
type PacketAckWritten struct {
	Packet PacketTypesPacket
	Ack    []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAckWritten is a free log retrieval operation binding the contract event 0x2061932a7cdcb993f721f99b41b0920bf92298de324c7725b5d88e5e3c301671.
//
// Solidity: event AckWritten((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) FilterAckWritten(opts *bind.FilterOpts) (*PacketAckWrittenIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "AckWritten")
	if err != nil {
		return nil, err
	}
	return &PacketAckWrittenIterator{contract: _Packet.contract, event: "AckWritten", logs: logs, sub: sub}, nil
}

// WatchAckWritten is a free log subscription operation binding the contract event 0x2061932a7cdcb993f721f99b41b0920bf92298de324c7725b5d88e5e3c301671.
//
// Solidity: event AckWritten((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) WatchAckWritten(opts *bind.WatchOpts, sink chan<- *PacketAckWritten) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "AckWritten")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketAckWritten)
				if err := _Packet.contract.UnpackLog(event, "AckWritten", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAckWritten is a log parse operation binding the contract event 0x2061932a7cdcb993f721f99b41b0920bf92298de324c7725b5d88e5e3c301671.
//
// Solidity: event AckWritten((string,string,uint64,string,bytes,bytes,string,uint64) packet, bytes ack)
func (_Packet *PacketFilterer) ParseAckWritten(log types.Log) (*PacketAckWritten, error) {
	event := new(PacketAckWritten)
	if err := _Packet.contract.UnpackLog(event, "AckWritten", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Packet contract.
type PacketInitializedIterator struct {
	Event *PacketInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketInitialized represents a Initialized event raised by the Packet contract.
type PacketInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Packet *PacketFilterer) FilterInitialized(opts *bind.FilterOpts) (*PacketInitializedIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PacketInitializedIterator{contract: _Packet.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Packet *PacketFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PacketInitialized) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketInitialized)
				if err := _Packet.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Packet *PacketFilterer) ParseInitialized(log types.Log) (*PacketInitialized, error) {
	event := new(PacketInitialized)
	if err := _Packet.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Packet contract.
type PacketOwnershipTransferredIterator struct {
	Event *PacketOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketOwnershipTransferred represents a OwnershipTransferred event raised by the Packet contract.
type PacketOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Packet *PacketFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PacketOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Packet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PacketOwnershipTransferredIterator{contract: _Packet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Packet *PacketFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PacketOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Packet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketOwnershipTransferred)
				if err := _Packet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Packet *PacketFilterer) ParseOwnershipTransferred(log types.Log) (*PacketOwnershipTransferred, error) {
	event := new(PacketOwnershipTransferred)
	if err := _Packet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketPacketReceivedIterator is returned from FilterPacketReceived and is used to iterate over the raw logs and unpacked data for PacketReceived events raised by the Packet contract.
type PacketPacketReceivedIterator struct {
	Event *PacketPacketReceived // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketPacketReceivedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketPacketReceived)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketPacketReceived)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketPacketReceivedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketPacketReceivedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketPacketReceived represents a PacketReceived event raised by the Packet contract.
type PacketPacketReceived struct {
	Packet PacketTypesPacket
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPacketReceived is a free log retrieval operation binding the contract event 0x1a78a2119315e61a73287bf641797ac20ee88a70a4244e54cf069dea3fc0af1e.
//
// Solidity: event PacketReceived((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketFilterer) FilterPacketReceived(opts *bind.FilterOpts) (*PacketPacketReceivedIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "PacketReceived")
	if err != nil {
		return nil, err
	}
	return &PacketPacketReceivedIterator{contract: _Packet.contract, event: "PacketReceived", logs: logs, sub: sub}, nil
}

// WatchPacketReceived is a free log subscription operation binding the contract event 0x1a78a2119315e61a73287bf641797ac20ee88a70a4244e54cf069dea3fc0af1e.
//
// Solidity: event PacketReceived((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketFilterer) WatchPacketReceived(opts *bind.WatchOpts, sink chan<- *PacketPacketReceived) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "PacketReceived")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketPacketReceived)
				if err := _Packet.contract.UnpackLog(event, "PacketReceived", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePacketReceived is a log parse operation binding the contract event 0x1a78a2119315e61a73287bf641797ac20ee88a70a4244e54cf069dea3fc0af1e.
//
// Solidity: event PacketReceived((string,string,uint64,string,bytes,bytes,string,uint64) packet)
func (_Packet *PacketFilterer) ParsePacketReceived(log types.Log) (*PacketPacketReceived, error) {
	event := new(PacketPacketReceived)
	if err := _Packet.contract.UnpackLog(event, "PacketReceived", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketPacketSentIterator is returned from FilterPacketSent and is used to iterate over the raw logs and unpacked data for PacketSent events raised by the Packet contract.
type PacketPacketSentIterator struct {
	Event *PacketPacketSent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketPacketSentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketPacketSent)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketPacketSent)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketPacketSentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketPacketSentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketPacketSent represents a PacketSent event raised by the Packet contract.
type PacketPacketSent struct {
	PacketBytes []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPacketSent is a free log retrieval operation binding the contract event 0x4cd1f79eb551dd9f49759ed1e187e8fdf4beee52dbc7a294f95aabd1a049c1fc.
//
// Solidity: event PacketSent(bytes packetBytes)
func (_Packet *PacketFilterer) FilterPacketSent(opts *bind.FilterOpts) (*PacketPacketSentIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "PacketSent")
	if err != nil {
		return nil, err
	}
	return &PacketPacketSentIterator{contract: _Packet.contract, event: "PacketSent", logs: logs, sub: sub}, nil
}

// WatchPacketSent is a free log subscription operation binding the contract event 0x4cd1f79eb551dd9f49759ed1e187e8fdf4beee52dbc7a294f95aabd1a049c1fc.
//
// Solidity: event PacketSent(bytes packetBytes)
func (_Packet *PacketFilterer) WatchPacketSent(opts *bind.WatchOpts, sink chan<- *PacketPacketSent) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "PacketSent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketPacketSent)
				if err := _Packet.contract.UnpackLog(event, "PacketSent", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePacketSent is a log parse operation binding the contract event 0x4cd1f79eb551dd9f49759ed1e187e8fdf4beee52dbc7a294f95aabd1a049c1fc.
//
// Solidity: event PacketSent(bytes packetBytes)
func (_Packet *PacketFilterer) ParsePacketSent(log types.Log) (*PacketPacketSent, error) {
	event := new(PacketPacketSent)
	if err := _Packet.contract.UnpackLog(event, "PacketSent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Packet contract.
type PacketPausedIterator struct {
	Event *PacketPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketPaused represents a Paused event raised by the Packet contract.
type PacketPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Packet *PacketFilterer) FilterPaused(opts *bind.FilterOpts) (*PacketPausedIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &PacketPausedIterator{contract: _Packet.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Packet *PacketFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *PacketPaused) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketPaused)
				if err := _Packet.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Packet *PacketFilterer) ParsePaused(log types.Log) (*PacketPaused, error) {
	event := new(PacketPaused)
	if err := _Packet.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PacketUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Packet contract.
type PacketUnpausedIterator struct {
	Event *PacketUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PacketUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PacketUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PacketUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PacketUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PacketUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PacketUnpaused represents a Unpaused event raised by the Packet contract.
type PacketUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Packet *PacketFilterer) FilterUnpaused(opts *bind.FilterOpts) (*PacketUnpausedIterator, error) {

	logs, sub, err := _Packet.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &PacketUnpausedIterator{contract: _Packet.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Packet *PacketFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *PacketUnpaused) (event.Subscription, error) {

	logs, sub, err := _Packet.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PacketUnpaused)
				if err := _Packet.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Packet *PacketFilterer) ParseUnpaused(log types.Log) (*PacketUnpaused, error) {
	event := new(PacketUnpaused)
	if err := _Packet.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
