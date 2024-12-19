package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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
	configPrefix := "config-pro"
	configFileName := fmt.Sprintf("./deploy/%s.yaml", configPrefix)
	v := viper.New()
	// 文件路径设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("nacos 配置信息：%v", global.NacosConfig)
	// serverConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}
	// clientConifg
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	// 获取配置：GetConfig
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		panic(err)
	}
	// 将配置解析给 ServerConfig
	err = json.Unmarshal([]byte(content), global.ServerConfig)
	if err != nil {
		panic(err)
	}
	// 监听配置变化
	configClient.ListenConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group,
		OnChange: func(_, _, _, data string) {
			err = json.Unmarshal([]byte(data), global.ServerConfig)
			if err != nil {
				panic(err)
			}
			zap.S().Infof("服务配置更新：%v", global.ServerConfig)
			sendRestartSignal()
		},
	})
	zap.S().Infof("服务配置信息：%v", global.ServerConfig)
}
