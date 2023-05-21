package zlm

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/week_plans")
	//
	router.POST("/on_record_mp4", post)
}
