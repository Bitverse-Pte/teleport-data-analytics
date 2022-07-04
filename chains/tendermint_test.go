package chains

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/teleport-network/teleport-data-analytics/config"
	"testing"
)

func TestIBCClient_GetPackcetsByEvent(t *testing.T) {
	ibcClient,err:= NewTendermintClient(config.TendermintConfig{
		Url:       "abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:9090",
		ChainName: "teleport",
		ChainID:   "teleport_7001_1",
		AgentAddr: "0x0000000000000000000000000000000040000001",
		Frequency: 1,
		// Refer to block generation speed
		BatchNumber:        1,
		RevisedHeight:      1,
		StartHeight:        1,
		BalanceMonitorings: nil,
	})
	if err != nil {
		fmt.Println(err)
	}
	_,_= ibcClient.GetPackcetsByEvent()
}

func TestTeleport_GetBlockPackets(t *testing.T) {
	ibcClient,err:= NewTendermintClient(config.TendermintConfig{
		Url:       "abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:9090",
		ChainName: "teleport",
		ChainID:   "teleport_7001_1",
		AgentAddr: "0x0000000000000000000000000000000040000001",
		Frequency: 1,
		// Refer to block generation speed
		BatchNumber:        1,
		RevisedHeight:      1,
		StartHeight:        1,
		BalanceMonitorings: nil,
		ChannelInfo: []config.ChannelInfo{{
			ChainName:"cosmos-hub",
			ChannelID:"channel-475",
		},{
			ChainName:"teleport",
			ChannelID:"channel-1",
		}},
	})
	require.NoError(t, err)
	res,err:= ibcClient.GetBlockPackets(1315407)
	require.NoError(t, err)
	t.Log("res:",res)
}

func TestCosmos_GetBlockPackets(t *testing.T) {
	ibcClient,err:= NewTendermintClient(config.TendermintConfig{
		Url:       "rpc.sentry-02.theta-testnet.polypore.xyz:9090",
		ChainName: "cosmos-hub",
		ChainID:   "theta-testnet-001",
		AgentAddr: "0x0000000000000000000000000000000040000001",
		Frequency: 1,
		// Refer to block generation speed
		BatchNumber:        1,
		RevisedHeight:      1,
		StartHeight:        1,
		BalanceMonitorings: nil,
		ChannelInfo: []config.ChannelInfo{{
			ChainName:"cosmos-hub",
			ChannelID:"channel-475",
		},{
			ChainName:"teleport",
			ChannelID:"channel-1",
		}},
	})

	require.NoError(t, err)
	res,err:= ibcClient.GetBlockPackets(10854653)
	require.NoError(t, err)
	t.Logf("res:%+v",res)
}

func TestNewTendermintClient(t *testing.T) {
	s1 := []int{1,2,3}
	s2 := []int{4,5,6}
	s3 := append(s1,s2...)
	fmt.Println(s3)
}



