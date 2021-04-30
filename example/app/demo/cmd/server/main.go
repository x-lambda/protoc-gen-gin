package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/x-lambda/protoc-gen-gin-example/util/middleware/timeout"
	"github.com/x-lambda/protoc-gen-gin/example/util/conf"

	"github.com/gin-gonic/gin"
)

func main() {
	reload := make(chan struct{}, 1)
	stop := make(chan os.Signal, 1)

	// 监听配置文件变更
	conf.OnConfigChange(func() { reload <- struct{}{} })
	conf.WatchConfig()
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	fmt.Println("start server")
	startServer()

	for {
		select {
		case <-reload:
			os.Exit(0)
		case sg := <-stop:
			fmt.Println("exit ....")
			if sg == syscall.SIGINT {
				os.Exit(0)
			} else {
				os.Exit(0)
			}
		default:
			os.Exit(0)
		}
	}
}

func startServer() {
	// TODO ctx 处理
	router := gin.New()

	// middleware
	router.Use(timeout.Timeout(time.Millisecond * 500))

	register(router, internal)
	router.Run("127.0.0.1:8080")
}

func stopServer() {}
