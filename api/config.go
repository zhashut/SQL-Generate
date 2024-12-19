package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sql_generate/server"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: zhashut
 * Date: 2024/12/14
 * Time: 16:45
 * Description: No Description
 */

func GetAddress(c *gin.Context) {
	address, _ := server.GetAddress(c)
	zap.S().Infof("api GetAddress address: %v", address)
	ResponseSuccess(c, address)
}
