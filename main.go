package main

// @title bluebell
// @version 1.0
// @description 这是一个贴吧

// @contact.name jack bao
// @license.name Apache 2.0
// @host localhost:8080
// @BasePath /api/v1
import (
	"bluebell/controllers"
	"bluebell/dao/mysql"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// gin-swagger middleware
// swagger embed files
//GO WEB 开发通用脚手架模板
func main() {
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success")
	//3.初始化mysql连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	fmt.Println()
	defer mysql.Close()
	//4.初始化Redis连接
	//if err := redis.Init(settings.Conf.RedisConfig); err != nil {
	//	fmt.Printf("init redis failed, err:%v\n", err)
	//	return
	//}
	//defer redis.Close()
	//生成雪花算法id
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake faield, err:%v\n", err)
		return
	}
	//初始化gin框架内置翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans faield, err:%v\n", err)
		return
	}
	//5.注册路由 (开发模式)
	r := routes.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}

	//6.启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting")

}
