package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/14
 * Time: 10:36
 * Description: 词库相关路由
 */

func InitDict(r *gin.RouterGroup) {
	dictRouter := r.Group("/dict")
	{
		dictRouter.POST("/add", api.AddDict)
		dictRouter.GET("/my/add/list/page", api.GetMyAddDictListPage)
		dictRouter.GET("/my/list/page", api.GetMyDictListPage)
		dictRouter.GET("/get", api.GetDictByID)
		dictRouter.POST("/generate/sql", api.GenerateDictCreateSql)
		dictRouter.POST("/delete", api.DeletedDict)
		dictRouter.GET("/list/page", api.GetDictListPage)
		dictRouter.GET("/my/list", api.GetMyDictList)
	}
}
