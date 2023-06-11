package server

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

type UserService struct{}

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
	uid, err := db.CreateUser(ctx, u.UserName, u.UserAccount, password)
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
	user, err := db.GetUserWithPassword(ctx, account, md5Pwd)
	if err != nil {
		return nil, fmt.Errorf("账号或密码错误")
	}
	// 保存登录态 TODO 这里的 session 可能不是全局的，考虑做成全局 global
	if res := session.Get(USER_LOGIN_STATE); res == nil {
		session.Set(USER_LOGIN_STATE, user)
		_ = session.Save()
	}
	return user, nil
}

// GetLoginUser 获取当前登录用户 TODO 这里获取不到登录用户
func (s *UserService) GetLoginUser(ctx context.Context, session sessions.Session) (*models.User, error) {
	currentUser := session.Get(USER_LOGIN_STATE).(*models.User)
	if currentUser == nil {
		return nil, fmt.Errorf("未登录")
	}
	user, err := db.GetUserByID(ctx, currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败")
	}
	return user, nil
}

// Logout 用户退出
func (s *UserService) Logout(ctx context.Context, session sessions.Session) (bool, error) {
	currentUser := session.Get(USER_LOGIN_STATE).(*models.User)
	if currentUser == nil {
		return false, fmt.Errorf("未登录")
	}
	// 移除登录态
	session.Delete(USER_LOGIN_STATE)
	return true, nil
}

// IsAdmin 是否为管理员
func (s *UserService) IsAdmin(ctx context.Context, session sessions.Session) (bool, error) {
	currentUser := session.Get(USER_LOGIN_STATE).(*models.User)
	return currentUser != nil && currentUser.UserRole == ADMIN_ROLE, nil
}

// 密码加密
func encryptPassword(originPassword string) string {
	hash := md5.New()
	hash.Write([]byte(SALT))
	hash.Write([]byte(originPassword))
	encryptString := hex.EncodeToString(hash.Sum(nil))
	return encryptString
}
