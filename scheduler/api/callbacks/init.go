package callbacks

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/callback")

	router.GET("/start", getStart)
	router.GET("/stop", getStop)
}
