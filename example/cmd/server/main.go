package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	reload := make(chan int, 1)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	startServer()

	for {
		select {
		case <-reload:
		case sg := <-stop:
			stopServer()
			// 仿 nginx 使用 HUP 信号重载配置
			if sg == syscall.SIGHUP {
				startServer()
			} else {
				return
			}
		}
	}
}

func startServer() {
	router := gin.New()
	register(router, internal)
	router.Run()
}

func stopServer() {}
