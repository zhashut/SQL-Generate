package dao

import (
	"context"
	"fmt"
	"sql_generate/global"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 21:34
 * Description: user模块的数据处理
 */

func CreateUser(ctx context.Context, name, account, password string) (int64, error) {
	var user models.User
	if res := global.DB.Where(&models.User{UserAccount: account}).First(&user); res.RowsAffected > 0 {
		return user.ID, fmt.Errorf("用户已经存在")
	}
	user.UserName = name
	user.UserAccount = account
	user.UserPassword = password

	global.DB.Save(&user)
	return user.ID, nil
}
