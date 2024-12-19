package initialize

import (
	"go.uber.org/zap"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

/**
* Created with GoLand 2022.2.3.
* @author: zhashut
* Date: 2024/12/14
* Time: 14:11
* Description: No Description
 */

// 发送 SIGUSR2 信号以触发重启
func sendRestartSignal() {
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	if err != nil {
		zap.S().Error("无法找到当前进程：", err.Error())
		return
	}

	if err := process.Signal(syscall.SIGUSR2); err != nil {
		zap.S().Error("发送重启信号失败：", err.Error())
	}
}

func HandleSignals(listener net.Listener) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	for s := range sig {
		switch s {
		case syscall.SIGUSR2: // 触发零停机重启
			zap.S().Info("收到 SIGUSR2 信号，准备重启...")
			if err := SpawnChildProcess(listener); err != nil {
				zap.S().Error("重启失败：", err.Error())
			}
		case syscall.SIGINT, syscall.SIGTERM: // 优雅退出
			zap.S().Info("收到退出信号，正在优雅退出...")
			listener.Close()
			os.Exit(0)
		}
	}
}

// 获取父进程传递的监听器
func GetListenerFromParent() (net.Listener, error) {
	fd, err := strconv.Atoi(os.Getenv("APP_FD"))
	if err != nil {
		return nil, err
	}

	file := os.NewFile(uintptr(fd), "listener")
	listener, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// 杀死父进程
func KillParentProcess() {
	ppid, err := strconv.Atoi(os.Getenv("APP_PPID"))
	if err != nil {
		return
	}
	if p, err := os.FindProcess(ppid); err == nil {
		p.Signal(syscall.SIGQUIT)
	}
}

// 启动子进程
func SpawnChildProcess(listener net.Listener) error {
	argv0, err := os.Executable()
	if err != nil {
		return err
	}

	fd, _ := listener.(*net.TCPListener).File()
	files := []*os.File{os.Stdin, os.Stdout, os.Stderr, fd}

	// 设置环境变量
	os.Setenv("APP_FD", strconv.Itoa(int(fd.Fd())))
	os.Setenv("APP_PPID", strconv.Itoa(os.Getpid()))

	_, err = os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Files: files,
		Env:   os.Environ(),
	})
	return err
}
