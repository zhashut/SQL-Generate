package initialize

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sql_generate/global"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/22
 * Time: 18:28
 * Description: 初始化缓存配置
 */

func InitCache() {
	redisCfg := global.ServerConfig.RedisConfig
	global.CaChe = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password, // 密码
		DB:       redisCfg.Db,       // 数据库
		PoolSize: redisCfg.PoolSize, // 连接池大小
	})

	_, err := global.CaChe.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
