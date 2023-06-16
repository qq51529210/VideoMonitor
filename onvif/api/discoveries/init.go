package discoveries

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/discoveries")
	//
	// router.GET("/", list)
}
