package jobs

import (
	"github.com/go-co-op/gocron"
	"github.com/teleport-network/teleport-data-analytics/config"
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
		packetFeeToken   = "newPacketFeeToken" + time.Now().Format("20060102150405")
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

	tx := model.CrossChainTransaction{SrcChain: srcChain, DestChain: destChain, TokenName: existedTokenName, Amount: "0", PacketFeeDemandToken: existedTokenName, PacketFee: 0}
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
	tx = model.CrossChainTransaction{SrcChain: srcChain, DestChain: destChain, TokenName: tokenName, Amount: amount, PacketFeeDemandToken: packetFeeToken, PacketFee: packetFee}
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
	// query token table check fee value
	tokenfee := &model.Token{}
	if err := pktSvc.PktPool.DB.Where("bridge_id=? and  token_name = ? and token_type=? ", metrics.ID, tx.PacketFeeDemandToken, 0).First(tokenfee).Error; err != nil {
		t.Errorf("query token error:%+v", err)
		return
	}
	if tokenfee.Amt == packetFeeStr {
		t.Log("pkt fee amt is correct")
	}
	// query token table check value
	tokenVal := &model.Token{}
	if err := pktSvc.PktPool.DB.Where("bridge_id=? and  token_name = ? and token_type=? ", metrics.ID, tx.TokenName, 1).First(tokenVal).Error; err != nil {
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
