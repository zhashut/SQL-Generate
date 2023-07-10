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
 * Date: 2023/6/15
 * Time: 17:22
 * Description: 字段相关api
 */

var (
	fieldService = server.NewFieldService()
)

// AddField 添加字段
func AddField(c *gin.Context) {
	var req models.FieldInfoAddRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	id, err := fieldService.AddField(c, &req, user.ID)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, id)
}

// GetFieldByID 根据id获取字段
func GetFieldByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dict, err := fieldService.GetFieldByID(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, dict)
}

// DeletedField 删除字段
func DeletedField(c *gin.Context) {
	var req models.OnlyIDRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	b, err := fieldService.DeleteField(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, b)
}

// GetMyAddFieldListPage 分页获取当前用户创建的资源列表
func GetMyAddFieldListPage(c *gin.Context) {
	var req models.FieldInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	fieldList, err := fieldService.GetMyAddFieldListPage(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  fieldList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(fieldList)),
	}
	ResponseSuccess(c, resp)
}

// GetMyFieldListPage 分页获取当前用户可选的资源列表
func GetMyFieldListPage(c *gin.Context) {
	var req models.FieldInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	fieldList, err := fieldService.GetMyFieldListPage(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  fieldList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(fieldList)),
	}
	ResponseSuccess(c, resp)
}

// GetMyFieldList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func GetMyFieldList(c *gin.Context) {
	var req models.FieldInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	fieldList, err := fieldService.GetMyFieldList(c, &req, user)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, fieldList)
}

// GetFieldListPage 分页获取列表
func GetFieldListPage(c *gin.Context) {
	var req models.FieldInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	fieldList, err := fieldService.GetFieldListPage(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  fieldList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(fieldList)),
	}
	ResponseSuccess(c, resp)
}

// GenerateFieldCreateSql 生成创建表的 SQL
func GenerateFieldCreateSql(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	idStr := string(data)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	resp, err := fieldService.GenerateCreateSQL(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, resp)
}
