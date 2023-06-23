package global

import (
	"github.com/gin-contrib/sessions"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sql_generate/config"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 10:31
 * Description: 全局变量
 */

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	DB           *gorm.DB
	Session      sessions.Session
	CaChe        *redis.Client
)
