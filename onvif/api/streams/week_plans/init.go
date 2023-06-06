package weekplans

import (
	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/week_plans")
	//
	router.GET("/", list)
}
