package api

import (
	"github.com/gin-contrib/sessions"
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

// Register 用户注册
func Register(c *gin.Context) {
	var request models.UserRegister
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewUserService()
	resp, err := s.Register(c, &request)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	ResponseSuccess(c, resp)
}

// Login 用户登录
func Login(c *gin.Context) {
	var request models.UserLogin
	session := sessions.Default(c)
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := server.NewUserService()
	user, err := s.Login(c, request.UserAccount, request.Password, session)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	ResponseSuccess(c, user)
}

// LoginUser 获取当前登录用户
func LoginUser(c *gin.Context) {
	session := sessions.Default(c)
	s := server.NewUserService()
	user, err := s.GetLoginUser(c, session)
	if err != nil {
		ResponseErrorWithMsg(c, ErrorNotFound, err.Error())
		return
	}
	ResponseSuccess(c, user)
}
