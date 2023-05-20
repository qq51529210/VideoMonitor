package callbacks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/log"
)

func getStart(ctx *gin.Context) {
	log.Debug("call back start")
	ctx.Status(http.StatusOK)
}

func getStop(ctx *gin.Context) {
	log.Debug("call back stop")
	ctx.Status(http.StatusOK)
}
