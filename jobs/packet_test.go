package jobs

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/teleport-network/teleport-data-analytics/config"
	"github.com/teleport-network/teleport-data-analytics/jobs/datas"
	"github.com/teleport-network/teleport-data-analytics/model"
	"gorm.io/gorm"
	"testing"
	"time"
)

func Test_UpdateMetrics(t *testing.T) {
	var (
		srcChain         = "newScChain" + time.Now().Format("20060102150405")
		destChain        = "newDestChain" + time.Now().Format("20060102150405")
		amount           = "100"
		packetFee        = float64(100)
		packetFeeStr     = "100"
		tokenName        = "newTokenName" + time.Now().Format("20060102150405")
		existedTokenName = "usdt"
	)
	// just test new bridge
	scheduler := gocron.NewScheduler(time.UTC)
	pktSvc := NewPacketService(scheduler, config.LoadConfigs())
	// new bridge must not exist
	metrics := &model.SingleDirectionBridgeMetrics{}
	if err := pktSvc.PktPool.DB.Where("src_chain = ? AND dest_chain = ?", srcChain, destChain).First(metrics).Error; err == nil {
		if err != gorm.ErrRecordNotFound {
			t.Errorf("new bridge must not exist, but got error:%+v", err)
			return
		}
	}
	if metrics.ID != 0 {
		t.Errorf("new bridge must not exist, but got id:%+v", metrics.ID)
		return
	}

	tx := model.CrossChainTransaction{SrcChain: srcChain, DestChain: destChain, TokenName: existedTokenName, Amount: "0", PacketFee: 0}
	if err := pktSvc.PktPool.UpdateMetrics([]model.CrossChainTransaction{tx}); err != nil {
		t.Errorf("updateMetrics error:%+v", err)
		return
	}
	// validate new bridge
	if err := pktSvc.PktPool.DB.Where("src_chain = ? AND dest_chain = ?", srcChain, destChain).First(metrics).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			t.Errorf("new bridge must exist, but not found")
			return
		} else {
			t.Errorf("new bridge must exist, but got error:%+v", err)
		}
	}
	t.Logf("new bridge:%+v", metrics)
	// delete test bridge after test
	if err := pktSvc.PktPool.DB.Delete(metrics).Error; err != nil {
		t.Errorf("delete new bridge error:%+v", err)
		return
	}

	// new bridge & new token
	// prepare test data
	tx = model.CrossChainTransaction{SrcChain: srcChain, DestChain: destChain, TokenName: tokenName, Amount: amount, PacketFee: packetFee}
	if err := pktSvc.PktPool.UpdateMetrics([]model.CrossChainTransaction{tx}); err != nil {
		t.Errorf("updateMetrics error:%+v", err)
		return
	}

	// query new bridge
	metrics = &model.SingleDirectionBridgeMetrics{}
	if err := pktSvc.PktPool.DB.Where("src_chain = ? AND dest_chain = ?", srcChain, destChain).First(metrics).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			t.Errorf("new bridge must exist, but not found")
			return
		} else {
			t.Errorf("new bridge must exist, but got error:%+v", err)
		}
	}
	// query token table check packet fee
	tokenfee := &model.Token{}
	if err := pktSvc.PktPool.DB.Where("bridge_id=? and  token_name = ? and token_type=? ", metrics.ID, tx.TokenName, model.PacketFee).First(tokenfee).Error; err != nil {
		t.Errorf("query token error:%+v", err)
		return
	}
	if tokenfee.Amt == packetFeeStr {
		t.Log("pkt fee amt is correct")
	}
	// query token table check packet value
	tokenVal := &model.Token{}
	if err := pktSvc.PktPool.DB.Where("bridge_id=? and  token_name = ? and token_type=? ", metrics.ID, tx.TokenName, model.PacketValue).First(tokenVal).Error; err != nil {
		t.Errorf("query token error:%+v", err)
		return
	}
	if tokenVal.Amt == amount {
		t.Log("pkt amt is correct")
	}
	// delete test bridge and tokens after test
	if err := pktSvc.PktPool.DB.Delete(metrics).Error; err != nil {
		t.Errorf("delete new bridge error:%+v", err)
		return
	}
	// delete tokenfee after test
	if err := pktSvc.PktPool.DB.Delete(tokenfee, tokenfee).Error; err != nil {
		t.Errorf("delete new token error:%+v", err)
		return
	}
	// delete token amt after test
	if err := pktSvc.PktPool.DB.Delete(tokenVal).Error; err != nil {
		t.Errorf("delete new token error:%+v", err)
		return
	}

}

func Test_InitSingleDirectionBridgeMetrics(t *testing.T) {

}

// Query cross_chain_tx table and calc all the metrics and
// compare with the result in SingleDirectionBridgeMetrics and tokens table
func Test_ConsistentMetrics(t *testing.T) {
	scheduler := gocron.NewScheduler(time.UTC)
	pktSvc := NewPacketService(scheduler, config.LoadConfigs())
	// new bridge must not exist
	//metrics := &model.SingleDirectionBridgeMetrics{}
	bridges := datas.Bridges
	fmt.Println()
	for src, dests := range datas.BridgeNameAdjList {
		for _, dest := range dests {
			bridgeKey := pktSvc.PktPool.ChainMap[src] + "/" + pktSvc.PktPool.ChainMap[dest]
			tokens := bridges[bridgeKey].Tokens
			for _, token := range tokens {
				t.Logf("Validate %s->%s : %s", src, dest, token.Name)
				// query single_direction_bridge_metrics
				metrics := &model.SingleDirectionBridgeMetrics{}
				if err := pktSvc.PktPool.DB.Where("src_chain = ? AND dest_chain = ?", src, dest).First(metrics).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						t.Errorf("new bridge must exist, but not found")
						continue
					} else {
						t.Errorf("new bridge must exist, but got error:%+v", err)
					}
				}
				// query tokens table
				tks := &[]model.Token{}
				if err := pktSvc.PktPool.DB.Where("bridge_id=? and  token_name = ?", metrics.ID, token.Name).Find(tks).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						t.Log("token not found")
						continue
					} else {
						t.Errorf("query token error:%+v", err)
						return
					}
				}

				for _, tk := range *tks {
					if tk.TokenType == model.PacketFee {
						// query cross_chain_tx table.
						// sum the packet fee of all the cross_chain_tx of the bridge
						packetFee := new(string)
						if err := pktSvc.PktPool.DB.Model(&model.CrossChainTransaction{}).Select("coalesce(sum(packet_fee),0) as fee_Amount").Where("src_chain = ? AND dest_chain = ? AND token_name = ? ", src, dest, token.Name).Find(packetFee).Error; err != nil {
							t.Errorf("query cross_chain_tx error:%+v", err)
							return
						}
						if *packetFee == tk.Amt {
							t.Log("packet fee is correct")
						} else {
							t.Errorf("cross_chain_tx table packet fee is not equal that in token table :token table:%+v, cross_chain_tx table:%+v", tk.Amt, *packetFee)
							return
						}
					} else if tk.TokenType == model.PacketValue {
						// query cross_chain_tx table.
						// sum the packet value of all the cross_chain_tx of the bridge
						packetValue := new(string)
						if err := pktSvc.PktPool.DB.Model(&model.CrossChainTransaction{}).Select("coalesce(sum(amount),0) as value_Amount").Where("src_chain = ? AND dest_chain = ? AND token_name = ? ", src, dest, token.Name).Find(packetValue).Error; err != nil {
							t.Errorf("query cross_chain_tx error:%+v", err)
							return
						}
						if *packetValue == tk.Amt {
							t.Log("packet value is correct")
						} else {
							t.Errorf("cross_chain_tx table packet value is not equal that in token table :token table:%+v, cross_chain_tx table:%+v", tk.Amt, *packetValue)
							return
						}
					}
				}
			}
		}
	}
}
