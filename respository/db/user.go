package db

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

type UserDao struct{}

// CreateUser 新建用户
func (dao *UserDao) CreateUser(ctx context.Context, name, account, password string) (int64, error) {
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

// GetUserWithPassword 根据账号密码查询用户
func (dao *UserDao) GetUserWithPassword(ctx context.Context, account, password string) (*models.User, error) {
	var user models.User
	if res := global.DB.Where("userAccount = ? and userPassword = ?", account, password).First(&user); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

// GetUserByID 根据ID查询用户
func (dao *UserDao) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	if res := global.DB.Where(&models.User{ID: id}).First(&user); res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
