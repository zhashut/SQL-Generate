package server

import (
	"context"
	"github.com/gin-contrib/sessions"
	"gorm.io/gorm"
	"sql_generate/core/schema"
	"sql_generate/models"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/12
 * Time: 18:44
 * Description: 通用方法/接口
 */

type UserResolver interface {
	GetLoginUser(ctx context.Context, session sessions.Session) (*models.User, error)
	IsAdmin(ctx context.Context, session sessions.Session) (bool, error)
}

type GenerateResolver interface {
	GenerateAll(tableSchema *schema.TableSchema) (*models.Generate, error)
	ValidField(field schema.Field) error
}

// Paginate 封装分页方法
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 20:
			pageSize = 20
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
