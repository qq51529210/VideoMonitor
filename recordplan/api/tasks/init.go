package tasks

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/tasks")
	//
	router.GET("/:id", get)
	router.PUT("/", put)
	router.DELETE("/", delete)
}

// type tasks struct {
// 	// 流的唯一标识
// 	ID string `json:"id" binding:"required"`
// 	// 开始录像回调
// 	StartCallback *string `json:"startCallback" binding:"required,max=255"`
// 	// 停止录像回调
// 	StopCallback *string `json:"stopCallback" binding:"required,max=255"`
// }
