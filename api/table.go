package api

import (
	"github.com/gin-gonic/gin"
	"sql_generate/models"
	"sql_generate/server"
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
