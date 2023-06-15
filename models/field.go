package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/15
 * Time: 15:45
 * Description: 字段模型
 */

type FieldInfo struct {
	ID            int64                 `gorm:"column:id;primaryKey;autoIncrement;not null comment 'id'" json:"id"`
	Name          string                `gorm:"column:name;type:varchar(512) comment '名称'" json:"name"`
	FieldName     string                `gorm:"column:fieldName;type:varchar(512) comment '字段名称'" json:"fieldName"`
	Content       string                `gorm:"column:content;type:text comment '表信息（json）'" json:"content"`
	ReviewStatus  int                   `gorm:"column:reviewStatus;type:int;not null comment '状态(0-待审核,1-通过,2-拒绝)'" json:"reviewStatus"`
	ReviewMessage string                `gorm:"column:reviewMessage;type:varchar(512) comment '审核信息'" json:"reviewMessage"`
	UserId        int64                 `gorm:"column:userId;type:bigint;not null comment '创建用户id'" json:"userId"`
	IsDelete      soft_delete.DeletedAt `gorm:"column:isDelete;softDelete:flag" json:"isDelete"`
	CreateTime    time.Time             `gorm:"column:createTime;autoCreateTime" json:"createTime"`
	UpdateTime    time.Time             `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
}

type FieldInfoAddRequest struct {
	Name    string `form:"name" json:"name"`
	Content string `form:"content" json:"content"`
}

type FieldInfoQueryRequest struct {
	Name         string `form:"name" json:"name"`
	SearchName   string `form:"searchName" json:"searchName"`
	FieldName    string `form:"fieldName" json:"fieldName"`
	Content      string `form:"content" json:"content"`
	ReviewStatus int    `form:"reviewStatus" json:"reviewStatus"`
	UserID       int64  `form:"userId" json:"userId"`
	PageRequest
}
