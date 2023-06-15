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
 * Date: 2023/6/14
 * Time: 10:43
 * Description: No Description
 */

var (
	dictService = server.NewDictService()
)

// AddDict 添加词条
func AddDict(c *gin.Context) {
	var req models.DictAddRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	id, err := dictService.AddDict(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, id)
}

// GetDictByID 根据id获取词条
func GetDictByID(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dict, err := dictService.GetDictByID(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, dict)
}

// DeletedDict 删除词条
func DeletedDict(c *gin.Context) {
	var req models.OnlyIDRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	b, err := dictService.DeleteDict(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, b)
}

// GetMyAddDictListPage 分页获取当前用户创建的资源列表
func GetMyAddDictListPage(c *gin.Context) {
	var req models.DictQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dictList, err := dictService.GetMyAddDictListPage(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  dictList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(dictList)),
	}
	ResponseSuccess(c, resp)
}

// GetMyDictListPage 分页获取当前用户可选的资源列表
func GetMyDictListPage(c *gin.Context) {
	var req models.DictQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dictList, err := dictService.GetMyDictListPage(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  dictList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(dictList)),
	}
	ResponseSuccess(c, resp)
}

// GetMyDictList 获取当前用户可选的全部资源列表（只返回 id 和名称）
func GetMyDictList(c *gin.Context) {
	var req models.DictQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dictList, err := dictService.GetMyDictList(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, dictList)
}

// GetDictListPage 分页获取列表
func GetDictListPage(c *gin.Context) {
	var req models.DictQueryRequest
	if err := c.ShouldBind(&req); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	dictList, err := dictService.GetDictListPage(c, &req)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	resp := &PageInfo{
		Records:  dictList,
		Pages:    req.Pages,
		PageSize: req.PageSize,
		Total:    int64(len(dictList)),
	}
	ResponseSuccess(c, resp)
}

// GenerateDictCreateSql 生成创建表的 SQL
func GenerateDictCreateSql(c *gin.Context) {
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
	resp, err := dictService.GenerateCreateSQL(c, id)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorPERATION, err.Error())
		return
	}
	ResponseSuccess(c, resp)
}
