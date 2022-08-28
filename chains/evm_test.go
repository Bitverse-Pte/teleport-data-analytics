package chains

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/teleport-network/teleport-data-analytics/tools"
	packettypes "github.com/teleport-network/teleport/x/xibc/core/packet/types"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	"github.com/teleport-network/teleport-data-analytics/config"
)

var (
	proxyAddr   = "0xe505e4f31527eb196e61156dfb96725b40394d23"
	proxyTopic  = "SendEvent(bytes,string,string,uint256)"
	agentAddr   = "0x0000000000000000000000000000000040000001"
	agentTopic  = "SendEvent(bytes,string,string,uint256)"
	packetTopic = "PacketSent(bytes)"
	ackTopic    = "AckWritten((string,string,uint64,string,bytes,bytes,string,uint64),bytes)"
)

func TestEth_GetMultiInfo2(t *testing.T) {
	cfg5 := config.EvmConfig{
		EvmUrl:         "https://evm-rpc.testnet.teleport.network",
		ChainName:      "teleport",
		ChainID:        "8001",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0xc7b95",
		AgentTopic:     agentTopic,
		AgentAddr:      agentAddr,
	}
	ethcli5, err := NewEvmCli(cfg5)
	require.NoError(t, err)
	height5, err := ethcli5.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println(height5)
	// "0xbc53cbfa10726c4536b9740de614b791e337ba167992529ec33fbbce1359b9ef
	packetTx, err := ethcli5.GetPacketsByHash("0xbc53cbfa10726c4536b9740de614b791e337ba167992529ec33fbbce1359b9ef")
	require.NoError(t, err)
	fmt.Println(packetTx)

}

func TestEth_GetLatest(t *testing.T) {
	cfg := config.EvmConfig{
		EvmUrl:         "https://arb-rinkeby.g.alchemy.com/v2/FEMPS7Mn-uYhyziwXOT6-EvmVAOWCeXl",
		ChainName:      "arbitrum",
		ChainID:        "421611",
		PacketContract: "0xc7b95",
		EndPointAddr:   "0xc7b95",
	}
	cfg2 := config.EvmConfig{
		EvmUrl:         "https://rinkeby.arbitrum.io/rpc",
		ChainName:      "arbitrum",
		ChainID:        "421611",
		PacketContract: "0xc7b95",
		EndPointAddr:   "0xc7b95",
	}
	cfg3 := config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0xc7b95",
	}

	//cfg4 := config.EvmConfig{
	//	EvmUrl:         "https://rinkeby.davionlabs.com",
	//	ChainName:      "rinkeby",
	//	ChainID:        "4",
	//	PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
	//	EndPointAddr: "0xc7b95",
	//}
	ethcli, err := NewEvmCli(cfg)
	require.NoError(t, err)
	height, err := ethcli.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println(height)

	ethcli2, err := NewEvmCli(cfg2)
	require.NoError(t, err)
	height2, err := ethcli2.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println(height2)

	ethcli3, err := NewEvmCli(cfg3)
	require.NoError(t, err)
	height3, err := ethcli3.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println(height3)

	//ethcli4, err := NewEvmCli(cfg4)
	//require.NoError(t, err)
	//height4,err:=  ethcli4.GetLatestHeight()
	//require.NoError(t, err)
	//fmt.Println(height4)

	cfg5 := config.EvmConfig{
		EvmUrl:         "https://evm-rpc.testnet.teleport.network",
		ChainName:      "teleport",
		ChainID:        "8001",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0xc7b95",
		AgentTopic:     agentTopic,
		AgentAddr:      agentAddr,
	}

	ethcli5, err := NewEvmCli(cfg5)
	require.NoError(t, err)
	height5, err := ethcli5.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println(height5)

	cfg6 := config.EvmConfig{
		EvmUrl:         "https://bsc.davionlabs.com",
		ChainName:      "teleport",
		ChainID:        "8001",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0xc7b95",
		AgentTopic:     agentTopic,
		AgentAddr:      agentAddr,
	}

	ethcli6, err := NewEvmCli(cfg6)
	require.NoError(t, err)
	height6, err := ethcli6.GetLatestHeight()
	require.NoError(t, err)
	fmt.Println("bsctest Height", height6)

}

func TestTeleport_NewToken(t *testing.T) {
	cfg := config.EvmConfig{
		EvmUrl:         "http://abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0x",
		AgentAddr:      agentAddr,
		AgentTopic:     agentTopic,
	}
	cli, err := NewEvmCli(cfg)
	require.NoError(t, err)
	token := NewErc20TokenQuery(cli,"0xc3296bd022783b8c851853b123df065778169a3f","usdt")
	d,err := token.GetDecimals()
	t.Logf("token decimals:%d",d)

}

func TestTeleport_GetPacketFee(t *testing.T) {
	cfg := config.EvmConfig{
		EvmUrl:         "http://abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "0x0000000000000000000000000000000020000001",
		PacketTopic:    packetTopic,
		AckTopic:       ackTopic,
		EndPointAddr:   "0x0000000000000000000000000000000020000002",
		AgentAddr:      agentAddr,
		AgentTopic:     agentTopic,
	}
	ethcli, err := NewEvmCli(cfg)
	require.NoError(t, err)
	packetTxs, err := ethcli.getPackets(1149646, 1249646)
	t.Log(packetTxs)
	packetFee,_ := ethcli.GetPacketFee("teleport","rinkeby",75)
	t.Log(packetFee)
}

func TestTeleport_PacketStatus(t *testing.T) {
	cfg := config.EvmConfig{
		EvmUrl:         "http://abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "0x0000000000000000000000000000000020000001",
		PacketTopic:    packetTopic,
		AckTopic:       ackTopic,
		EndPointAddr:   "0x0000000000000000000000000000000020000002",
		AgentAddr:      agentAddr,
		AgentTopic:     agentTopic,
	}
	ethcli, err := NewEvmCli(cfg)
	require.NoError(t, err)
	packetTxs, err := ethcli.getPackets(1149646, 1249646)
	require.NoError(t, err)
	fmt.Println(packetTxs)
	packets:= []packettypes.Packet{{
		SrcChain:  packetTxs[0].SrcChain,
	}}
	b, _ := json.Marshal(packets)
	res,err := tools.Post("https://bridge.qa.davionlabs.com/bridge/status",bytes.NewBuffer(b),nil)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(res))
}

func TestTeleport_GetMultiInfo(t *testing.T) {
	cfg := config.EvmConfig{
		EvmUrl:         "http://abd46ec6e28754f0ab2aae29deaa0c11-1510914274.ap-southeast-1.elb.amazonaws.com:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "0x0000000000000000000000000000000020000001",
		PacketTopic:    packetTopic,
		AckTopic:       ackTopic,
		EndPointAddr:   "0x0000000000000000000000000000000020000002",
		AgentAddr:      agentAddr,
		AgentTopic:     agentTopic,
	}
	ethcli, err := NewEvmCli(cfg)
	require.NoError(t, err)
	//packets, err := ethcli.getAckPackets(1180875, 1180875)
	//fmt.Println(packets)
	tendermintCfg := config.TendermintConfig{
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
	}
	_ = NewTeleport(tendermintCfg, ethcli,nil)
	if err != nil {
		t.Log(err)
	}
	packets:= []packettypes.Packet{}
	t.Log(packets)
	b, _ := json.Marshal(packets)
	res,err := tools.Post("https://bridge.qa.davionlabs.com/bridge/status",bytes.NewBuffer(b),nil)
	if err != nil {
		fmt.Println(err)
	}
	t.Log(string(res))
}

func TestEth_GetHeightByHash(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
	})
	require.NoError(t, err)
	height, err := ethcli.GetHeightByHash("0x83a81f457c1b5823a73b13048f6c3ab60d58d2bfea9cf10b0f83d781584be733")
	require.NoError(t, err)

	txHash := common.HexToHash("0x83a81f457c1b5823a73b13048f6c3ab60d58d2bfea9cf10b0f83d781584be733")
	tx, _, err := ethcli.ethClient.TransactionByHash(context.Background(), txHash)
	require.NoError(t, err)
	require.Equal(t, "0x83a81f457c1b5823a73b13048f6c3ab60d58d2bfea9cf10b0f83d781584be733", tx.Hash().String())
	fmt.Printf("txHash:%v\n", tx.Hash().String())
	fmt.Printf("gas:%v\n", tx.Gas())
	fmt.Printf("gasPrice:%v\n", tx.GasPrice())
	fmt.Printf("gas * gasPrice:%v\n", tx.GasPrice().Uint64()*tx.Gas())
	scanGasUsed := 0.000441609006624135 / 0.000000001000000015
	fmt.Printf("scanUsed:%v\n", scanGasUsed)
	fmt.Printf("block height:%v\n", height)
}

func TestEth_GetPacket(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0xf7268301384fb751e49fafdacd02c693eabb142c",
		PacketTopic:    "PacketSent(bytes)",
		AckTopic:       "AckWritten((string,string,uint64,string,bytes,bytes,string,uint64),bytes)",
		ReceivedAckTopic: "AckPacket((string,string,uint64,string,bytes,bytes,string,uint64),bytes)",
		EndPointAddr:   "0xe4916fd50499601dfe4fd2b40ee6d93a8035fcab",
	})
	if err != nil {
		panic(err)
	}
	packets, err := ethcli.getReceivedAcks(10963774, 10964774)
	require.NoError(t, err)
	t.Log(packets)
}

func TestBsc_GetPacket(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://data-seed-prebsc-1-s1.binance.org:8545",
		ChainName:      "bsctest",
		ChainID:        "97",
		PacketContract: "0x8e84ef5d13a129183b838d833e4ac14eb0c5ceab",
		PacketTopic:    "PacketSent(bytes)",
		AckTopic:       "AckWritten((string,string,uint64,string,bytes,bytes,string,uint64),bytes)",
		ReceivedAckTopic: "AckPacket((string,string,uint64,string,bytes,bytes,string,uint64),bytes)",
		EndPointAddr:   "0xfe30de51bdb9b9784f1a05d5531d221bf66eaf70",
	})
	if err != nil {
		panic(err)
	}
	packets, err := ethcli.GetPacketsByHash("0x6be2e3c7341b8ce966e39bb44b0be9404d6a38235f037eb651c2b7b12d1f07e0")
	require.NoError(t, err)
	t.Logf("%+v",packets)
	packets, err = ethcli.GetPacketsByHash("0x05bf57b48a65ed7daa20d736bae4665f5b6f106790ca907243b7e1c59cd40fb0")
	require.NoError(t, err)
	t.Logf("%+v",packets)
	packets, err = ethcli.GetPacketsByHash("0xb00a45a6385b962271509b55d44a5fc49d04115a1f753e2f7df819bac8f7af02")
	require.NoError(t, err)
	t.Logf("%+v",packets)
}

func TestEvm_GetOutTokenAmount(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0xc7b952b46a1115a51b7395ca3b7f473c6b37befe",
		EndPointAddr:   "0x41baacc9cf251b1046d72610bbc96af69e03ed0d",
	})
	if err != nil {
		panic(err)
	}
	tokenAmount, err := ethcli.GetOutTokenAmount("0xd81d7ba3926fae537c05eb39ec597c891cbf8827", "teleport", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(tokenAmount)

}

func TestEvm_GetOutTokenAmount2(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "http://10.41.20.10:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "",
		EndPointAddr:   "0x0000000000000000000000000000000030000001",
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}
	// eth token
	amount, err := ethcli.GetOutTokenAmount("0x0000000000000000000000000000000000000000", "rinkeby", nil)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Println(amount)
}

func TestEvm_GetInTokenAmount(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0xea0f10eea3ab1ea7a9628e5d21346744a05f2bae",
		PacketTopic:    packetTopic,
		EndPointAddr:   "0x01b1ed143cf7c80f1658c10baab3faf900d64f98",
	})
	if err != nil {
		panic(err)
	}
	packetTx, err := ethcli.GetPacketsByHash("0x5a2308d0eef0396e9750da9bf693cc8b0d9ea8e412c3483f6d2d32e14e9f4c")
	require.NoError(t, err)
	fmt.Println(packetTx)
	amount, err := ethcli.GetInTokenAmount("0x2b2454ad0c2142bd02ff38d8728c022a4a90feb7", "teleport", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(amount)
}

func TestEvm_GetTokenLimit(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
		ChainName:      "rinkeby",
		ChainID:        "4",
		PacketContract: "0x6c034eb404cb960cd51bf2570f8a71d61f1f39e7",
		EndPointAddr:   "0x466a9a8978f01197f22072efbb1709769fa78059",
	})
	if err != nil {
		panic(err)
	}
	addr := common.HexToAddress("0x2b2454ad0c2142bd02ff38d8728c022a4a90feb7")
	amountLimit, err := ethcli.GetTokenLimit(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(amountLimit)
}

func TestEvm_GetInTokenAmount2(t *testing.T) {
	ethcli, err := NewEvmCli(config.EvmConfig{
		EvmUrl:         "http://10.41.20.10:8545",
		ChainName:      "teleport",
		ChainID:        "7001",
		PacketContract: "",
		EndPointAddr:   "0x0000000000000000000000000000000030000001",
	})
	if err != nil {
		fmt.Printf("%+v", err)
	}
	// eth token
	amount, err := ethcli.GetInTokenAmount("0x0000000000000000000000000000000000000000", "rinkeby", nil)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	fmt.Println(amount)
}
