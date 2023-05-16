package streams

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/streams")
	//
	router.GET("/:stream", get)
}
