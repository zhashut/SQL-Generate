package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"sql_generate/consts"
	"sql_generate/models"
	"sql_generate/respository/ses"
	"sql_generate/server"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 22:59
 * Description: 用户相关api
 */

var (
	userService = server.NewUserService()
)

func userToVo(user *models.User) *models.UserVo {
	return &models.UserVo{
		ID:          user.ID,
		UserName:    user.UserName,
		UserAccount: user.UserAccount,
		UserAvatar:  user.UserAvatar,
		Gender:      user.Gender,
		UserRole:    user.UserRole,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

// Register 用户注册
func Register(c *gin.Context) {
	var request models.UserRegister
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	resp, err := userService.Register(c, &request)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	ResponseSuccess(c, resp)
}

// Login 用户登录
func Login(c *gin.Context) {
	var request models.UserLogin
	if err := c.ShouldBind(&request); err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}
	s := sessions.Default(c)
	user, err := userService.Login(c, request.UserAccount, request.Password, s)
	if err != nil {
		ResponseFailed(c, ErrorInvalidParams)
		return
	}

	ResponseSuccess(c, userToVo(user))
}

// LoginUser 获取当前登录用户
func LoginUser(c *gin.Context) {
	// 获取登录用户信息
	user := ses.GetSession(c, consts.USER_LOGIN_STATE)
	if user == nil {
		return
	}

	ResponseSuccess(c, userToVo(user))
}

// Logout 用户退出
func Logout(c *gin.Context) {
	ses.DeleteSession(c, consts.USER_LOGIN_STATE)
	ResponseSuccess(c, "ok")
}
