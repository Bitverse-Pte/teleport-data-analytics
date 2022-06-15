package datas

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoadStateServer(t *testing.T) {
	route := gin.New()
	route.Static("/remoteData", "./remoteData")
	route.Run(":8000")
}
