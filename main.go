package main

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"sql_generate/global"
	"sql_generate/initialize"
	"syscall"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/5/14
 * Time: 22:10
 * Description: No Description
 */

func main() {
	// 初始化配置
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitCache()
	r := initialize.Router()

	// 启动监听端口
	zap.S().Debugf("启动服务器，端口：%d", global.ServerConfig.Port)
	go func() {
		if err := r.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Fatal("启动失败", err.Error())
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Infof("服务退出成功")
}
