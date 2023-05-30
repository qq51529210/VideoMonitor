package streams

import (
	weekplans "recordplan/api/streams/week_plans"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/streams/:stream")
	//
	weekplans.Init(router)
}
