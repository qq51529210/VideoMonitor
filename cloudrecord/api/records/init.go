package records

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/records")
	//
	router.GET("/", list)
	router.GET("/:id", get)
	router.POST("/", post)
	router.DELETE("/:id", delete)
	router.DELETE("/", batchDelete)
}
