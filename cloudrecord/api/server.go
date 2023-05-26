package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/qq51529210/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"cloudrecord/api/internal/docs"
	"cloudrecord/api/records"

	"github.com/gin-gonic/gin"
)

// 服务实例
var (
	ser   server
	Serve = ser.Serve
)

// server 表示 http api 服务
type server struct {
	gin *gin.Engine
	ser http.Server
}

// Serve 初始化后开始服务
func (s *server) Serve(port int) {
	// 路由
	s.initRouter()
	// http 服务
	s.ser.Addr = fmt.Sprintf(":%d", port)
	err := s.ser.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// initRouter 初始化路由
func (s *server) initRouter() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	s.gin = gin.New()
	s.ser.Handler = s.gin
	// 文档
	docs.SwaggerInfo.BasePath = "/"
	s.gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.gin.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"文档地址": "/docs/swagger/index.html"})
	})
	// 中间件
	s.gin.Use(global)
	// 路由
	records.Init(s.gin)
}

// 全局第一个中间件
func global(ctx *gin.Context) {
	now := time.Now()
	var cost time.Duration
	// 清理
	defer func() {
		re := recover()
		if re == nil {
			log.Debugf("%s %s %s cost %v", ctx.Request.RemoteAddr, ctx.Request.Method,
				ctx.Request.RequestURI, cost)
		} else {
			log.Debugf("%s %s %s panic %v", ctx.Request.RemoteAddr, ctx.Request.Method,
				ctx.Request.RequestURI, re)
		}
	}()
	// 执行
	ctx.Next()
	cost = time.Since(now)
}
