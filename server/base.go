package server

import (
	"context"
	"github.com/gin-contrib/sessions"
	"math/rand"
	"sql_generate/core/schema"
	"sql_generate/models"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/12
 * Time: 18:44
 * Description: 通用方法/接口
 */

const (
	LIST_PAGE     = "listPage"
	ADD_LIST_PAGE = "addListPage:"
	MY_LIST_PAGE  = "myListPage:"
	MY_LIST       = "myList:"
	LIST          = "list"
)

type UserResolver interface {
	GetLoginUser(ctx context.Context, session sessions.Session) (*models.User, error)
	IsAdmin(ctx context.Context, session sessions.Session) (bool, error)
}

type GenerateResolver interface {
	GenerateAll(tableSchema *schema.TableSchema) (*models.Generate, error)
	ValidField(field *schema.Field) error
	ValidSchema(tableSchema *schema.TableSchema) error
}

type BuilderResolver interface {
	BuildCreateFieldSQL(field *schema.Field) (string, error)
	BuildCreateTableSql(tableSchema *schema.TableSchema) (string, error)
}

func randomTime(expireHour int) time.Duration {
	rand.Seed(time.Now().UnixMilli())
	return time.Duration(rand.Intn(11)+(expireHour*3600)) * time.Second
}
