package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"simple_file_server/config"
	"simple_file_server/router"
	"simple_file_server/zlogs"
)

func main() {
	zlogs.InitLogger()
	logger := zlogs.GetLogger()
	defer logger.Sync()

	zlogs.Info("服务正在启动...")
	if err := start(); err != nil {
		zlogs.Panic(fmt.Sprintf("服务启动失败:", err.Error()))
	}
}

func start() error {
	s := &http.Server{
		Addr:    ":8899",
		Handler: ginEngine(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctxFunc context.CancelFunc) {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT)
		for {
			select {
			case <-quit:
				ctxFunc()
				return
			}
		}
	}(cancel)

	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			zlogs.Panic(fmt.Sprintf("服务退出异常:", err.Error()))
			panic(500)
		}
	}()

	zlogs.Info("服务启动成功！")
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		zlogs.Debug("服务已正常关闭")
		return nil
	}
	return err
}

func ginEngine() *gin.Engine {
	gin.SetMode(func() string {
		if config.IsDev() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())

	e := gin.New()
	e.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Server error",
		})
	}))

	zlogs.Info("初始化路由配置...")
	router.InitRouter(e)
	return e
}
