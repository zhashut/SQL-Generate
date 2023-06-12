package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:51
 * Description: No Description
 */

func InitTable(r *gin.RouterGroup) {
	tableRouter := r.Group("/table_info")
	{
		tableRouter.POST("/add", api.AddTableInfo)
		tableRouter.GET("/my/add/list/page", api.GetMyTableInfoList)
	}
}
