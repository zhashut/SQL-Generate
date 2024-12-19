package main

import (
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
	"sql_generate/global"
	"sql_generate/initialize"
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

	// 定义监听器
	var listener net.Listener
	// 获取 address
	address := fmt.Sprintf("%s:%d", global.ServerConfig.Host, global.ServerConfig.Port)
	// 检查是否有父进程传递的监听器（用于零停机重启）
	if fd, err := initialize.GetListenerFromParent(); err != nil {
		zap.S().Info("没有父进程传递的监听器，正常启动")
		listener, err = net.Listen("tcp", address)
		if err != nil {
			zap.S().Fatal("监听端口失败：", err.Error())
		}
	} else {
		listener = fd
		initialize.KillParentProcess()
	}

	// 启动 HTTP 服务
	go func() {
		if err := r.RunListener(listener); err != nil {
			zap.S().Fatal("服务器启动失败：", err.Error())
		}
	}()
	zap.S().Infof("服务器启动成功，监听端口: %d, PID: %d", global.ServerConfig.Port, os.Getpid())

	// 启动信号监听，处理优雅退出和零停机重启
	initialize.HandleSignals(listener)
}
