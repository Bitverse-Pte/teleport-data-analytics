package monitoring

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/teleport-network/teleport-data-analytics/chains"
	"github.com/teleport-network/teleport-data-analytics/metrics"
	"github.com/teleport-network/teleport-data-analytics/model"
)

type Monitoring struct {
	balanceMonitorings []*BalanceMonitoring
	metricsManager     *metrics.MetricManager
	chains             map[string]chains.BlockChain
	db                 *gorm.DB
	log                *logrus.Logger
	lightClients       []string //chainName
}

func NewMonitoring(log *logrus.Logger, metricsManager *metrics.MetricManager, balanceMonitorings []*BalanceMonitoring, chains map[string]chains.BlockChain, db *gorm.DB) *Monitoring {
	return &Monitoring{
		log:                log,
		balanceMonitorings: balanceMonitorings,
		metricsManager:     metricsManager,
		chains:             chains,
		db:                 db,
	}
}

func (m *Monitoring) Monitoring(scheduler *gocron.Scheduler) {
	_, err := scheduler.Every(30).Seconds().Do(func() {
		m.balanceMonitoring()
	})
	if err != nil {
		panic(fmt.Errorf("balanceMonitoring scheduler.Every exec error:%+v", err))
	}
	_, err = scheduler.Every(30).Seconds().Do(func() {
		m.lightClientMonitoring()
	})
	if err != nil {
		panic(fmt.Errorf("lightClientMonitoring scheduler.Every exec error:%+v", err))
	}
	_, err = scheduler.Every(30).Seconds().Do(func() {
		m.pendingPacketMonitoring()
	})

	// every day exposed the bridge metrics
	//_, err = scheduler.Every(24).Hours().Do(func() {
	//	m.bridgeMetrics()
	//})

	if err != nil {
		panic(fmt.Errorf("pendingPacketMonitoring scheduler.Every exec error:%+v", err))
	}
}

func (m *Monitoring) balanceMonitoring() {
	for _, balanceMonitoring := range m.balanceMonitorings {
		tokenQuery := balanceMonitoring.TokensQuery
		for _, account := range balanceMonitoring.Accounts {
			balance, err := tokenQuery.GetTokenBalance(account.AddressHex, nil)
			if err != nil {
				m.log.Errorf("tokenQuery.GetErc20Balance error:%+v", err)
			}
			balanceFloat, err := balance.Float64()
			if err != nil {
				m.log.Errorf("balance.Float64 error:%+v", err)
			}
			m.metricsManager.Gauge.With("chain_name", tokenQuery.ChainName()).
				With("option", fmt.Sprintf("%v_%v-balance", account.Name, tokenQuery.TokenName())).Set(balanceFloat)
		}
	}
}

func (m *Monitoring) lightClientMonitoring() {
	for _, cn := range m.lightClients {
		latestHeight, err := m.chains[cn].GetLatestHeight()
		if err != nil {
			m.log.Errorf("lightClientMonitoring GetLatestHeight error:%+v", err)
		}
		height, err := m.chains[chains.TeleportChain].GetLightClientHeight(cn)
		if err != nil {
			m.log.Errorf("lightClientMonitoring GetLightClientHeight error:%+v", err)
		}
		m.metricsManager.Gauge.With("chain_name", cn).With("option", "light-client_low-height").Set(float64(int(latestHeight) - int(height)))
	}
}

func (m *Monitoring) pendingPacketMonitoring() {
	t := time.Now().AddDate(0, 0, -1)
	for cn := range m.chains {
		var totalCount int64
		if err := m.db.Model(&model.CrossChainTransaction{}).Where("dest_chain = ? and status = ? and send_tx_time > ?", cn, model.Pending, t).Count(&totalCount).Error; err != nil {
			m.log.Errorf("pending packets query count error:%+v", err)
			return
		}
		m.metricsManager.Gauge.With("chain_name", cn).With("option", "pending_count").Set(float64(totalCount))
	}
}

// bridgeMetrics exposed the bridge metrics
//func (m *Monitoring) bridgeMetrics() {
//	// query from singleDirectionBridgeMetrics table
//	// travels all the chainNames
//	var (
//		metric = &model.SingleDirectionBridgeMetrics{}
//		err    error
//	)
//
//	for chainName, chainNames := range datas.BridgeNameAdjList {
//		// travel all the chainNames in the adjList
//		for _, cn := range chainNames {
//			// query the total count of packets from the bridge
//			if err = m.db.Where("src_chain = ? and dest_chain = ?", chainName, cn).Last(metric).Error; err != nil {
//				if err == gorm.ErrRecordNotFound {
//					m.log.Warnf("singleDirectionBridgeMetrics query error:%+v", err)
//					continue
//				} else {
//					m.log.Errorf("singleDirectionBridgeMetrics query error:%+v", err)
//					continue
//				}
//			}
//			// pkt amount on the bridge
//			m.metricsManager.BridgeGauge.With("src_chain", chainName).With("dest_chain", cn).With("option", "pkt_amt").Set(float64(metric.PktAmt))
//			// pkt failed amount on the bridge
//			m.metricsManager.BridgeGauge.With("src_chain", chainName).With("dest_chain", cn).With("option", "failed_amt").Set(float64(metric.FailedAmt))
//			// usdt amount on the bridge
//			// TODO: precision issue
//			tmpV, _ := strconv.ParseFloat(metric.UAmt, 64)
//			m.metricsManager.BridgeGauge.With("src_chain", chainName).With("dest_chain", cn).With("option", "usdt_amt").Set(tmpV)
//			// tele amount on the bridge
//			tmpV, _ = strconv.ParseFloat(metric.TeleAmt, 64)
//			m.metricsManager.BridgeGauge.With("src_chain", chainName).With("dest_chain", cn).With("option", "tele_amt").Set(tmpV)
//			// TODO: fee amount distinguish between different tokens
//			tmpV, _ = strconv.ParseFloat(metric.FeeAmt, 64)
//			m.metricsManager.BridgeGauge.With("src_chain", chainName).With("dest_chain", cn).With("option", "fee_amt").Set(tmpV)
//		}
//	}
//
//}
