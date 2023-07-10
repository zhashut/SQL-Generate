package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	. "sql_generate/consts"
	"sql_generate/models"
	"sql_generate/respository/db"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 21:35
 * Description: 用户模块的业务层
 */

type UserService struct {
	DB *db.UserDao
}

func NewUserService() *UserService {
	return &UserService{}
}

// Register 用户注册
func (s *UserService) Register(ctx context.Context, u *models.UserRegister) (*int64, error) {
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
	uid, err := s.DB.CreateUser(ctx, u.UserName, u.UserAccount, password)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, account, password string, session sessions.Session) (*models.User, error) {
	if account == "" || password == "" {
		return nil, fmt.Errorf("参数不能为空")
	}
	if len(account) < 4 {
		return nil, fmt.Errorf("账号不得少于4位")
	}
	if len(password) < 8 {
		return nil, fmt.Errorf("密码不得少于8位")
	}
	md5Pwd := encryptPassword(password)
	user, err := s.DB.GetUserWithPassword(ctx, account, md5Pwd)
	if err != nil {
		return nil, fmt.Errorf("账号或密码错误")
	}
	// 保存登录态
	if res := session.Get(USER_LOGIN_STATE); res == nil {
		marshal, _ := json.Marshal(&user)
		session.Set(USER_LOGIN_STATE, marshal)
		_ = session.Save()
	}
	return user, nil
}

// 密码加密
func encryptPassword(originPassword string) string {
	hash := md5.New()
	hash.Write([]byte(SALT))
	hash.Write([]byte(originPassword))
	encryptString := hex.EncodeToString(hash.Sum(nil))
	return encryptString
}
