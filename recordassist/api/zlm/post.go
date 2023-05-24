package zlm

import (
	"net/http"
	"recordassist/api/internal"
	"recordassist/db"

	"github.com/gin-gonic/gin"
	"github.com/qq51529210/log"
)

// "mediaServerId" : "your_server_id",
// "app" : "live",
// "file_name" : "15-53-02.mp4",
// "file_path" : "/root/zlmediakit/httpRoot/__defaultVhost__/record/live/obs/2019-09-20/15-53-02.mp4",
// "file_size" : 1913597,
// "folder" : "/root/zlmediakit/httpRoot/__defaultVhost__/record/live/obs/",
// "start_time" : 1568965982,
// "stream" : "obs",
// "time_len" : 11.0,
// "url" : "record/live/obs/2019-09-20/15-53-02.mp4",
// "vhost" : "__defaultVhost__"
type postReq struct {
	App       string  `json:"app"`
	Stream    string  `json:"stream"`
	FilePath  string  `json:"file_path"`
	FileSize  int64   `json:"file_size"`
	StartTime int64   `json:"start_time"`
	TimeLen   float64 `json:"time_len"`
	// VHost         string  `json:"vhost"`
	// FileName      string  `json:"file_name"`
	// Folder        string  `json:"folder"`
	// URL           string  `json:"url"`
	// MediaServerID string  `json:"mediaServerId"`
}

type postRes struct {
	Code int `json:"code"`
}

func post(ctx *gin.Context) {
	// 参数
	var req postReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}
	// 数据库
	var model db.Record
	model.App = req.App
	model.Stream = req.Stream
	model.Path = req.FilePath
	model.Size = req.FileSize
	model.CreateTime = req.StartTime
	model.Duration = req.TimeLen
	model.Status = db.RecordStatusReady
	_, err = db.Add(&model)
	if err != nil {
		log.Errorf("add file %s to db error:%s", req.FilePath, err.Error())
		ctx.JSON(http.StatusOK, &postRes{Code: 500})
		return
	}
	ctx.JSON(http.StatusOK, &postRes{Code: 0})
}
