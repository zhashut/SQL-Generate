package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 9:10
 * Description: No Description
 */

func InitGenerateSQL(r *gin.RouterGroup) {
	sqlRouter := r.Group("/sql")
	{
		sqlRouter.POST("/generate/schema", api.GenerateSQL)
	}
}
