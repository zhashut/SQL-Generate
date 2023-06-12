package api

import (
	"context"
	"sql_generate/global"
	"sql_generate/models"
	"sql_generate/server"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/13
 * Time: 6:46
 * Description: 通用
 */

func GetUserBySession(c context.Context) (*models.User, error) {
	// 获取当前用户id
	userService := server.NewUserService()
	return userService.GetLoginUser(c, global.Session)
}

type PageInfo struct {
	Records  interface{} `json:"records"`
	Pages    int64       `json:"current"`
	PageSize int64       `json:"size"`
	Total    int64       `json:"total"`
}
