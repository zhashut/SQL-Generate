package models

import (
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 18:12
 * Description: 表模型
 */

type TableInfo struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement;not null comment 'id'" json:"id"`
	Name          string    `gorm:"column:name;type:varchar(512) comment '名称'" json:"name"`
	Content       string    `gorm:"column:content;type:text comment '表信息（json）'"`
	ReviewStatus  int       `gorm:"column:reviewStatus;type:int;not null comment '状态(0-待审核,1-通过,2-拒绝)'"`
	ReviewMessage string    `gorm:"column:reviewMessage;type:varchar(512) comment '审核信息'"`
	UserId        int64     `gorm:"column:userId;type:bigint;not null comment '创建用户id'"`
	IsDelete      bool      `gorm:"column:isDelete" json:"isDelete"`
	CreateTime    time.Time `gorm:"column:createTime;autoCreateTime" json:"createTime"`
	UpdateTime    time.Time `gorm:"column:updateTime;autoUpdateTime" json:"updateTime"`
}

type TableInfoAddRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
