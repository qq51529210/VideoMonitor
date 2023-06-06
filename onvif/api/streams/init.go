package streams

import (
	weekplans "recordplan/api/streams/week_plans"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/streams")
	//
	router.DELETE("/", batchDelete)
	router.DELETE("/:stream", delete)
	//
	weekplans.Init(router.Group("/:stream"))
}
