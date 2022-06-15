package database

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"github.com/teleport-network/teleport-data-analytics/model"
)

func TestInitDB(t *testing.T) {
	myDB := InitDB("root:123456@tcp(127.0.0.1:3306)/backend?charset=utf8mb4&parseTime=True&loc=Local")
	var data model.CrossChainTransaction
	if err := myDB.Model(&model.CrossChainTransaction{}).Where("src_chain = ? and dest_chain = ? and sequence = ? ", "src", "dest", 2).Find(&data).Error; err != nil {
		fmt.Println(err == gorm.ErrRecordNotFound)
	}
}
