package api

import (
	"github.com/gin-gonic/gin"
	"sql_generate/consts"
	"sql_generate/models"
	"sql_generate/respository/ses"
	"sql_generate/server"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 17:21
 * Description: 举报相关api
 */

var (
	reportService = server.NewReportService()
)

// AddReport 添加举报信息
func AddReport(c *gin.Context) {
	var req models.ReportAddRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	id, err := reportService.AddReport(c, &req, user.ID)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, id)
}

// GetReportByID 根据id获取举报信息
func GetReportByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	report, err := reportService.GetReportByID(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, report)
}

// DeletedReport 删除举报信息
func DeletedReport(c *gin.Context) {
	var req models.OnlyIDRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	b, err := reportService.DeleteReport(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, b)
}

// UpdateReport 更新举报信息
func UpdateReport(c *gin.Context) {
	var req models.ReportUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	id, err := reportService.UpdateReport(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, id)
}

// GetReportListPage 分页获取列表
func GetReportListPage(c *gin.Context) {
	var req models.ReportQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	reportList, err := reportService.GetReportListPage(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  reportList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(reportList)),
	}
	ResponseSuccess(c, resp)
}

// GetReportList 获取列表
func GetReportList(c *gin.Context) {
	var req models.ReportQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	reportList, err := reportService.GetReportList(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  reportList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(reportList)),
	}
	ResponseSuccess(c, resp)
}
