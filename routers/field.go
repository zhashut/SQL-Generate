package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/15
 * Time: 17:25
 * Description: 字段相关路由
 */

func InitField(r *gin.RouterGroup) {
	fieldRouter := r.Group("/field_info")
	{
		fieldRouter.POST("/add", api.AddField)
		fieldRouter.GET("/my/add/list/page", api.GetMyAddFieldListPage)
		fieldRouter.GET("/my/list/page", api.GetMyFieldListPage)
		fieldRouter.GET("/get", api.GetFieldByID)
		fieldRouter.POST("/generate/sql", api.GenerateFieldCreateSql)
		fieldRouter.POST("/delete", api.DeletedField)
		fieldRouter.GET("/list/page", api.GetFieldListPage)
		fieldRouter.GET("/my/list", api.GetMyFieldList)
	}
}
