package streams

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/tasks")
	//
	router.PUT("/", put)
	router.DELETE("/", delete)
}

type stream struct {
	// 流的唯一标识
	Stream string `json:"stream" binding:"required,max=128"`
	// 开始录像回调
	StartCallback *string `json:"startCallback" binding:"required,max=255"`
	// 停止录像回调
	StopCallback *string `json:"stopCallback" binding:"required,max=255"`
}
