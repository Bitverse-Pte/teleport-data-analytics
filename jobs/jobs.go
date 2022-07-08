package jobs

import (
	"fmt"

	"github.com/go-co-op/gocron"

	"github.com/sirupsen/logrus"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/config"
	"github.com/teleport-network/teleport-data-analytics/jobs/bridges"
	"github.com/teleport-network/teleport-data-analytics/jobs/datas"
	"github.com/teleport-network/teleport-data-analytics/jobs/monitoring"
	"github.com/teleport-network/teleport-data-analytics/jobs/packet"
	"github.com/teleport-network/teleport-data-analytics/metrics"
	"github.com/teleport-network/teleport-data-analytics/repo/database"
)

type PacketService struct {
	PktPool *packet.PacketPool
}

func NewPacketService(scheduler *gocron.Scheduler, cfg *config.Config) *PacketService {
	if err := datas.LoadState(cfg.Network); err != nil {
		panic(err.Error())
	}
	datas.ReloadState(cfg.Network, cfg.ReloadPeriod)
	db := database.InitDB(cfg.MysqlAddr)
	log := logrus.New()
	cs := make(map[string]chains.BlockChain)
	var balanceMonitorings []*monitoring.BalanceMonitoring
	chainMap := chains.NewChainMap()
	chainCliMap := make(map[string]chains.BlockChain)
	for _, evmCfg := range cfg.EvmChains {
		evmChain, err := chains.NewEvmCli(evmCfg)
		if err != nil {
			panic(fmt.Errorf("chains.NewEvmCli error:%v\nchainName:%v,chainID:%v", err.Error(), evmCfg.ChainName, evmCfg.ChainID))
		} else {
			logrus.Printf("NewEvmCli %v success", evmCfg.ChainName)
		}
		chainMap[chainMap.GetXIBCChainKey(evmChain.ChainName())] = evmCfg.ChainID
		chainCliMap[evmCfg.ChainID] = evmChain
		cs[evmCfg.ChainName] = evmChain
		for _, token := range evmCfg.BalanceMonitorings {
			tokenQuery := chains.NewTokenQuery(evmChain, token.AddressHex, token.TokenName, 0)
			balanceMonitoring := monitoring.NewBalanceMonitoring(tokenQuery, token.Accounts)
			balanceMonitorings = append(balanceMonitorings, balanceMonitoring)
		}
	}
	teleEvm := cs[cfg.Teleport.ChainName]
	teleportEvm := teleEvm
	for _, tendermintCfg := range cfg.TendermintChains {
		tendermintCli, err := chains.NewTendermintClient(tendermintCfg)
		if err != nil {
			panic(fmt.Errorf("chains.NewEvmCli error:%v\nchainName:%v,chainID:%v", err.Error(), tendermintCfg.ChainName, tendermintCfg.ChainID))
		} else {
			logrus.Printf("NewTendermintClient %v success", tendermintCfg.ChainName)
		}
		chainMap[chainMap.GetIBCChainKey(tendermintCfg.ChainName)] = tendermintCfg.ChainID
		chainCliMap[tendermintCfg.ChainID] = tendermintCli
		cs[tendermintCfg.ChainName] = tendermintCli
		for _, token := range tendermintCfg.BalanceMonitorings {
			tokenQuery := chains.NewTokenQuery(tendermintCli, token.AddressHex, token.TokenName, 0)
			balanceMonitoring := monitoring.NewBalanceMonitoring(tokenQuery, token.Accounts)
			balanceMonitorings = append(balanceMonitorings, balanceMonitoring)
		}
	}
	teleportTendermint := cs[cfg.Teleport.ChainName]
	if teleportEvm == nil {
		log.Fatalln("teleport evm client not init")
	}

	if teleportTendermint == nil {
		log.Fatalln("teleport tendermint client not init")
	}
	teleChain := chains.NewTeleport(cfg.Teleport, teleportEvm,teleportTendermint)
	if cfg.Teleport.AgentAddr != "" {
		logrus.Infof("teleport agent addr:%voverwrite default：%v", cfg.Teleport.AgentAddr, chains.AgentContract)
		chains.AgentContract = cfg.Teleport.AgentAddr
	} else {
		logrus.Infof("teleport agent addr is empty,use default：%v", chains.AgentContract)
	}
	cs[cfg.Teleport.ChainName] = teleChain
	metricsManager := metrics.NewMetricManager()
	reconciliationCli := bridges.NewBridge(log, db, datas.Bridges, chainCliMap)
	pool := packet.NewPacketDBPool(db, log, cs, chainMap, reconciliationCli, cfg.ReconcileEnable, metricsManager)
	monitoringSrv := monitoring.NewMonitoring(log, metricsManager, balanceMonitorings, cs, db, chainMap)
	monitoringSrv.Monitoring(scheduler)
	return &PacketService{
		PktPool: pool,
	}
}
