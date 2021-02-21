package routes

import (
	"bluebell/controllers"
	_ "bluebell/docs" // swagger embed files
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //设置成发布模式
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.LoadHTMLFiles("templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// 接口文档生成
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controllers.SignupHandler)
	//登录业务路由
	v1.POST("/login", controllers.LoginHandler)

	v1.GET("/community", controllers.CommunityHandler)
	v1.GET("/community/:id", controllers.CommunityDetailHandler)
	//帖子详情
	v1.GET("/post/:id", controllers.PostedDetailHandler)
	v1.GET("/posts", controllers.GetPostListHandler)
	//根据时间或分数获取帖子列表
	v1.GET("/posts2", controllers.GetPostListHandler2)

	v1.Use(middlewares.JWTAuthMiddleware()) //应用JWT中间件
	{
		//发帖
		v1.POST("/post", controllers.PostedHandler)
		//投票
		v1.POST("/vote", controllers.PostVoteController)
	}
	// pprof 性能分析
	pprof.Register(r)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
