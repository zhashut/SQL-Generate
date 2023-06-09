package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sql_generate/consts"
	"sql_generate/dao"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 21:35
 * Description: No Description
 */

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) UserRegister(ctx context.Context, u *models.UserRegister) (*int64, error) {
	if u.UserName == "" || u.UserAccount == "" || u.Password == "" || u.CheckPassword == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	if len(u.UserName) > 16 {
		return nil, fmt.Errorf("用户名过长")
	}
	if len(u.UserAccount) < 4 {
		return nil, fmt.Errorf("用户账户过短")
	}
	if len(u.Password) < 8 || len(u.CheckPassword) < 8 {
		return nil, fmt.Errorf("用户密码果断")
	}
	// 密码和校验密码相同
	if u.Password != u.CheckPassword {
		return nil, fmt.Errorf("两次输入的密码不一致")
	}
	password := encryptPassword(u.Password)
	uid, err := dao.CreateUser(ctx, u.UserName, u.UserAccount, password)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}

// 密码加密
func encryptPassword(originPassword string) string {
	hash := md5.New()
	hash.Write([]byte(consts.SALT))
	hash.Write([]byte(originPassword))
	encryptString := hex.EncodeToString(hash.Sum(nil))
	return encryptString
}
