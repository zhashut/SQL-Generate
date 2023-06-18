package models

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 15:32
 * Description: 举报模型
 */

type Report struct {
	ID             int64                 `gorm:"column:id;primaryKey;autoIncrement;not null comment 'id'" json:"id"`
	ReportedID     int64                 `gorm:"column:reportedId;not null comment '被举报对象id'" json:"reportedId"`
	ReportedUserId int64                 `gorm:"column:reportedUserId;not null comment '被举报用户id'" json:"reportedUserId"`
	Content        string                `gorm:"column:content;type:text comment '内容'" json:"content"`
	Type           int                   `gorm:"column:type; comment '举报实体类型（0-词库）'" json:"type"`
	Status         int                   `gorm:"column:status;type:int;not null comment '状态(0-未处理 1-已处理)'" json:"reviewStatus"`
	UserId         int64                 `gorm:"column:userId;type:bigint;not null comment '创建用户id'" json:"userId"`
	IsDelete       soft_delete.DeletedAt `gorm:"column:isDelete;softDelete:flag" json:"isDelete"`
	CreateTime     time.Time             `gorm:"column:createTime;autoCreateTime" json:"createTime"`
	UpdateTime     time.Time             `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
}

type ReportAddRequest struct {
	Content    string `gorm:"column:content;type:text comment '内容'" json:"content"`
	Type       int    `gorm:"column:type; comment '举报实体类型（0-词库）'" json:"type"`
	ReportedID int64  `gorm:"column:reportedId;not null comment '被举报对象id'" json:"reportedId"`
}

type ReportQueryRequest struct {
	Content        string `gorm:"column:content;type:text comment '内容'" json:"content"`
	Type           int    `gorm:"column:type; comment '举报实体类型（0-词库）'" json:"type"`
	ReportedID     int64  `gorm:"column:reportedId;not null comment '被举报对象id'" json:"reportedId"`
	ReportedUserId int64  `gorm:"column:reportedUserId;not null comment '被举报用户id'" json:"reportedUserId"`
	Status         int    `gorm:"column:status;type:int;not null comment '状态(0-未处理 1-已处理)'" json:"reviewStatus"`
	UserId         int64  `gorm:"column:userId;type:bigint;not null comment '创建用户id'" json:"userId"`
	PageRequest
}

type ReportUpdateRequest struct {
	ID     int64 `gorm:"column:id;primaryKey;autoIncrement;not null comment 'id'" json:"id"`
	Type   int   `gorm:"column:type; comment '举报实体类型（0-词库）'" json:"type"`
	Status int   `gorm:"column:status;type:int;not null comment '状态(0-未处理 1-已处理)'" json:"reviewStatus"`
}
