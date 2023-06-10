package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 23:05
 * Description: No Description
 */

func InitUser(r *gin.RouterGroup) {
	sqlRouter := r.Group("/user")
	{
		sqlRouter.POST("/register", api.Register)
		sqlRouter.POST("/login", api.Login)
		sqlRouter.GET("/get/login", api.LoginUser)
	}
}
