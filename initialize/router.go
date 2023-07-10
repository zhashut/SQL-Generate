package initialize

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"sql_generate/global"
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
	reds := global.ServerConfig.RedisConfig
	store, _ := redis.NewStore(10, "tcp",
		fmt.Sprintf("%s:%d", reds.Host, reds.Port), reds.Password, []byte("secret"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   reds.ExpireHour * 60 * 60,
		HttpOnly: true,
	})
	r.Use(middlewares.Cors(), sessions.Sessions("zhashut", store))
	apiGroup := r.Group("/api")
	routers.InitGenerateSQL(apiGroup)
	routers.InitUser(apiGroup)
	routers.InitTable(apiGroup)
	routers.InitDict(apiGroup)
	routers.InitField(apiGroup)
	routers.InitReport(apiGroup)
	return r
}
