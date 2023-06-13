package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"sql_generate/core/builder"
	"sql_generate/core/schema"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/server"
	"strconv"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 19:57
 * Description: No Description
 */

// AddTableInfo  保存表
func AddTableInfo(c *gin.Context) {
	var request models.TableInfoAddRequest
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewTableService()
	tableInfo := &models.TableInfo{
		Name:    request.Name,
		Content: request.Content,
	}
	// 检验
	if err := s.ValidAndHandleTableInfo(c, tableInfo, true); err != nil {
		ResponseErrorWithMsg(c, ErrorInvalidParams, err.Error())
		return
	}
	// 获取当前用户
	user, err := GetUserBySession(c)
	if err != nil {
		ResponseFailed(c, ErrorNotLogin)
		return
	}
	tableInfo.UserId = user.ID
	// 保存表
	result, err := s.AddTableInfo(c, tableInfo)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorInvalidParams, err.Error())
		return
	}
	if !result {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, result)
}

// GetMyTableInfoList 分页获取当前用户创建的资源列表
func GetMyTableInfoList(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	// 获取当前登录用户
	user, err := GetUserBySession(c)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorNotLogin, err.Error())
		return
	}
	req.UserID = user.ID
	s := server.NewTableService()
	list, err := s.GetMyTableInfoList(c, &req)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	resp := &PageInfo{
		Records:  list,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(list)),
	}
	ResponseSuccess(c, resp)
}

// GetTableInfoByID 根据ID获取TableInfo
func GetTableInfoByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewTableService()
	resp, err := s.GetTableInfoById(c, id)
	if err != nil {
		ResponseFailed(c, ErrorNotFound)
		return
	}
	ResponseSuccess(c, resp)
}

// GenerateCreateSql 生成创建表的 SQL
func GenerateCreateSql(c *gin.Context) {
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
	s := server.NewTableService()
	tableInfo, err := s.GetTableInfoById(c, id)
	if err != nil {
		ResponseFailed(c, ErrorNotFound)
		return
	}
	// 将 Content Unmarshal到 schema 对象中
	var tableSchema *schema.TableSchema
	if err := json.Unmarshal([]byte(tableInfo.Content), &tableSchema); err != nil {
		ResponseFailed(c, ErrorPERATION)
		return
	}
	// 构造建表SQL
	createTableSQL, err := builder.NewSQLBuilder().BuildCreateTableSql(tableSchema)
	if err != nil {
		ResponseFailed(c, ErrorPERATION)
		return
	}
	ResponseSuccess(c, createTableSQL)
}

func DeletedTableInfo(c *gin.Context) {
	var req models.OnlyIDRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	us := server.NewUserService()
	user, err := us.GetLoginUser(c, global.Session)
	if err != nil {
		ResponseFailed(c, ErrorNotLogin)
		return
	}
	// 判断是否存在
	s := server.NewTableService()
	tableInfo, err := s.GetTableInfoById(c, req.ID)
	if err != nil {
		ResponseFailed(c, ErrorNotFound)
		return
	}
	// 仅本人或管理员可以删除
	admin, err := us.IsAdmin(c, global.Session)
	if err != nil {
		ResponseFailed(c, ErrorNoAuth)
		return
	}
	if tableInfo.UserId != user.ID && !admin {
		ResponseFailed(c, ErrorNoAuth)
		return
	}
	b, err := s.DeletedTableInfoByID(c, req.ID)
	if err != nil {
		ResponseFailed(c, ErrorPERATION)
		return
	}
	ResponseSuccess(c, b)
}

// GetTableInfoList 分页获取列表
func GetTableInfoList(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewTableService()
	list, err := s.GetTableInfoList(c, &req)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	resp := &PageInfo{
		Records:  list,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(list)),
	}
	ResponseSuccess(c, resp)
}
