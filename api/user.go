package api

import (
	"github.com/gin-gonic/gin"
	"sql_generate/models"
	"sql_generate/server"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 22:59
 * Description: 用户相关api
 */

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	var request models.UserRegister
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewUserService()
	resp, err := s.UserRegister(c, &request)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	ResponseSuccess(c, resp)
}
