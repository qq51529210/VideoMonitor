package devices

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/devices")
	//
	// router.GET("/", list)
	// router.GET("/:id", get)
	// router.POST("/", post)
	// router.PATCH("/:id", patch)
	// router.DELETE("/:id", delete)
	// router.DELETE("/", batchDelete)
}
