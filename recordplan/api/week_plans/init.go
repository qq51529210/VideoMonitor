package weekplans

import (
	"encoding/json"
	"recordplan/api/week_plans/streams"
	"recordplan/db"
	"strings"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(router gin.IRouter) {
	router = router.Group("/week_plans")
	//
	router.GET("/", list)
	router.GET("/:id", get)
	router.POST("/", post)
	router.PATCH("/:id", patch)
	router.DELETE("/:id", delete)
	router.DELETE("/", batchDelete)
	//
	streams.Init(router.Group("/:id"))
}

func jsonTimePeroid(peroids [][]*db.TimePeroid) *string {
	var buf strings.Builder
	json.NewEncoder(&buf).Encode(peroids)
	data := buf.String()
	return &data
}
