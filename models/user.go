package models

import (
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/9
 * Time: 20:31
 * Description: No Description
 */

type User struct {
	ID           int64     `gorm:"column:id;column:id;primaryKey;autoIncrement;not null comment 'id'" json:"id"`
	UserName     string    `gorm:"column:userName;type:varchar(256);comment '用户昵称'" json:"userName"`
	UserAccount  string    `gorm:"column:userAccount;type:varchar(256);not null comment '账号'" json:"userAccount"`
	UserPassword string    `gorm:"column:userPassword;type:varchar(256);not null comment '密码'" json:"userPassword"`
	UserAvatar   string    `gorm:"column:userAvatar;type:varchar(1024); comment '用户头像'" json:"userAvatar"`
	Gender       bool      `gorm:"column:gender" json:"gender"`
	UserRole     string    `gorm:"column:userRole;type:varchar(256);default:user;not null comment '用户角色: user/admin'" json:"userRole"`
	IsDelete     bool      `gorm:"column:isDelete" json:"isDelete"`
	CreatedAt    time.Time `gorm:"column:createTime" json:"createTime"`
	UpdatedAt    time.Time `gorm:"column:updateTime" json:"updateTime"`
}

type UserRegister struct {
	UserName      string `form:"userName" json:"userName" binding:"required"`
	UserAccount   string `form:"userAccount" json:"userAccount" binding:"required"`
	Password      string `form:"userPassword" json:"userPassword" binding:"required"`
	CheckPassword string `form:"checkPassword" json:"checkPassword" binding:"required"`
}
