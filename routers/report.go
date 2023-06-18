package routers

import (
	"github.com/gin-gonic/gin"
	"sql_generate/api"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 17:25
 * Description: 举报相关路由
 */

func InitReport(r *gin.RouterGroup) {
	reportRouter := r.Group("/report")
	{
		reportRouter.POST("/add", api.AddReport)
		reportRouter.GET("/get", api.GetReportByID)
		reportRouter.POST("/update", api.UpdateReport)
		reportRouter.POST("/delete", api.DeletedReport)
		reportRouter.GET("/list/page", api.GetReportListPage)
		reportRouter.GET("/list", api.GetReportList)
	}
}
