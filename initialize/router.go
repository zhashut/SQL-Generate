package initialize

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"sql_generate/middlewares"
	"sql_generate/routers"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:10
 * Description: 统一全局路由初始化
 */

func Router() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(middlewares.Cors(), sessions.Sessions("mysession", store))
	apiGroup := r.Group("/api")
	routers.InitGenerateSQL(apiGroup)
	routers.InitUser(apiGroup)
	routers.InitTable(apiGroup)
	return r
}
