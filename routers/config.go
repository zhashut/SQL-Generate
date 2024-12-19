package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: zhashut
 * Date: 2024/12/14
 * Time: 16:38
 * Description: 配置相关
 */

func InitExposureConfig(r *gin.RouterGroup) {
	dictRouter := r.Group("/config")
	{
		dictRouter.GET("/address", api.GetAddress)
	}
}
