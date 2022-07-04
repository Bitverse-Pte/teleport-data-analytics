package chains

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	typestx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	xibcbsc "github.com/teleport-network/teleport/x/xibc/clients/light-clients/bsc/types"
	xibceth "github.com/teleport-network/teleport/x/xibc/clients/light-clients/eth/types"
	xibctendermint "github.com/teleport-network/teleport/x/xibc/clients/light-clients/tendermint/types"
	clienttypes "github.com/teleport-network/teleport/x/xibc/core/client/types"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
)

var protoCodec *codec.ProtoCodec

func makeCodec() *codec.ProtoCodec {
	if protoCodec != nil {
		return protoCodec
	}
	ir := codectypes.NewInterfaceRegistry()
	clienttypes.RegisterInterfaces(ir)
	govtypes.RegisterInterfaces(ir)
	xibcbsc.RegisterInterfaces(ir)
	xibctendermint.RegisterInterfaces(ir)
	xibceth.RegisterInterfaces(ir)
	packettypes.RegisterInterfaces(ir)
	ir.RegisterInterface("cosmos.v1beta1.Msg", (*sdk.Msg)(nil))
	typestx.RegisterInterfaces(ir)
	cryptocodec.RegisterInterfaces(ir)
	// ibc
	channeltypes.RegisterInterfaces(ir)
	protoCodec = codec.NewProtoCodec(ir)
	return protoCodec
}
