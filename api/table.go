package api

import (
	"github.com/gin-gonic/gin"
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

var (
	tableService = server.NewTableService()
)

// AddTableInfo 添加表
func AddTableInfo(c *gin.Context) {
	var req models.TableInfoAddRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	id, err := tableService.AddTableInfo(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, id)
}

// GetTableInfoByID 根据id获取表
func GetTableInfoByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dict, err := tableService.GetTableInfoByID(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, dict)
}

// DeletedTableInfo 删除表
func DeletedTableInfo(c *gin.Context) {
	var req models.OnlyIDRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	b, err := tableService.DeleteTableInfo(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, b)
}

// GetMyAddTableInfoListPage 分页获取当前用户创建的资源列表
func GetMyAddTableInfoListPage(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	fieldList, err := tableService.GetMyAddTableInfoListPage(c, &req)
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

// GetMyTableInfoListPage 分页获取当前用户可选的资源列表
func GetMyTableInfoListPage(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	fieldList, err := tableService.GetMyTableInfoListPage(c, &req)
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

// GetMyTableInfoList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func GetMyTableInfoList(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	fieldList, err := tableService.GetMyTableInfoList(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, fieldList)
}

// GetTableInfoListPage 分页获取列表
func GetTableInfoListPage(c *gin.Context) {
	var req models.TableInfoQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	fieldList, err := tableService.GetTableInfoListPage(c, &req)
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

// GenerateTableInfoCreateSql 生成创建表的 SQL
func GenerateTableInfoCreateSql(c *gin.Context) {
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
	resp, err := tableService.GenerateCreateSQL(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, resp)
}
