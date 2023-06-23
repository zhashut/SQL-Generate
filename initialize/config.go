package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"sql_generate/global"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 10:32
 * Description: 初始化配置
 */

func InitConfig() {
	configPrefix := "config"
	configFileName := fmt.Sprintf("%s.yaml", configPrefix)
	v := viper.New()
	// 文件路径设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)
	// viper的功能-动态监控变化
	_ = v.WriteConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		zap.S().Infof("配置文件发生变化: %s", e.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(&global.ServerConfig)
		zap.S().Infof("配置信息：%v", global.ServerConfig)
	})
	zap.S().Infof("%v", global.ServerConfig)
}
