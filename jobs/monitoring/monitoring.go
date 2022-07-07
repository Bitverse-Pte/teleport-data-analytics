package monitoring

import (
	"fmt"
	"github.com/teleport-network/teleport-data-analytics/jobs/datas"
	"math/big"
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
	chainMap           chains.ChainMap
}

func NewMonitoring(log *logrus.Logger, metricsManager *metrics.MetricManager, balanceMonitorings []*BalanceMonitoring, chains map[string]chains.BlockChain, db *gorm.DB, chainMap chains.ChainMap) *Monitoring {
	return &Monitoring{
		log:                log,
		balanceMonitorings: balanceMonitorings,
		metricsManager:     metricsManager,
		chains:             chains,
		db:                 db,
		chainMap:           chainMap,
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
func (m *Monitoring) bridgeMetrics() {
	// query from singleDirectionBridgeMetrics table
	// travels all the dests
	var (
		metric = &model.SingleDirectionBridgeMetrics{}
		err    error
	)

	for src, dests := range datas.BridgeNameAdjList {
		// travel all the dests in the adjList
		for _, cn := range dests {
			metric = &model.SingleDirectionBridgeMetrics{}
			// query the total count of packets from the bridge
			if err = m.db.Where("src_chain = ? and dest_chain = ?", src, cn).Last(metric).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					m.log.Warnf("singleDirectionBridgeMetrics query error:%+v", err)
					continue
				} else {
					m.log.Errorf("singleDirectionBridgeMetrics query error:%+v", err)
					continue
				}
			}
			// pkt amount on the bridge
			m.metricsManager.BridgeGauge.With("src_chain", src).With("dest_chain", cn).With("option", "pkt_amt").Set(float64(metric.PktAmt))
			// pkt failed amount on the bridge
			m.metricsManager.BridgeGauge.With("src_chain", src).With("dest_chain", cn).With("option", "failed_amt").Set(float64(metric.FailedAmt))
			// bridgeToken amount on the bridge

			// for all bridgeTokens under the bridge
			// get bridge key
			bridgeKey := m.chainMap.GetXIBCChainKey(src) + "/" + m.chainMap.GetXIBCChainKey(cn)
			bridgeTokens := datas.Bridges[bridgeKey].Tokens
			for _, bridgeToken := range bridgeTokens {
				// query packet fee
				tokens := &[]model.Token{}
				if err = m.db.Where("bridge_id= ? and token_name = ? ", metric.ID, bridgeToken.Name).Find(tokens).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						m.log.Warnf("singleDirectionBridgeMetrics query error:%+v", err)
						continue
					} else {
						m.log.Errorf("singleDirectionBridgeMetrics query error:%+v", err)
						continue
					}
				}

				for _, token := range *tokens {
					bigflt, ok := big.NewFloat(float64(0)).SetString(token.Amt)
					if !ok {
						m.log.Errorf("parse token Amt error")
						continue
					}
					flt, accuracy := bigflt.Float64()
					if accuracy != big.Exact {
						m.log.Errorf("amt overflow")
						continue
					}
					m.metricsManager.BridgeGauge.With("src_chain", src).With("dest_chain", cn).With("token_name", token.TokenName).With("token_type", token.TokenType.String()).Set(flt)
				}

				// query packet amt
			}
		}
	}

}
