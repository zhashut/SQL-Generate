package global

import (
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
)
