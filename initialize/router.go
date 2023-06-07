package initialize

import (
	"github.com/gin-gonic/gin"
	"sql_generate/middlewares"
	"sql_generate/routers"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:10
 * Description: No Description
 */

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.Cors())
	apiGroup := r.Group("/api")
	routers.InitGenerateSQL(apiGroup)
	return r
}
