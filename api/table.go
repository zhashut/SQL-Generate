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

// AddTableInfo TODO 保存表，userId是硬编码
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
	// TODO 获取当前用户id
	//userService := server.NewUserService()
	//user, err := userService.GetLoginUser(c, session)
	//if err != nil {
	//	ResponseFailed(c, ErrorNotLogin)
	//	return
	//}
	tableInfo.UserId = 6
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
