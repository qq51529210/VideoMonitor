package zlm

import (
	"recordassist/api/internal"

	"github.com/gin-gonic/gin"
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
	VHost     string `json:"vhost"`
	App       string `json:"app"`
	Stream    string `json:"stream"`
	FileName  string `json:"file_name"`
	FilePath  string `json:"file_path"`
	FileSize  string `json:"file_size"`
	Folder    string `json:"folder"`
	StartTime string `json:"start_time"`
	TimeLen   string `json:"time_len"`
	URL       string `json:"url"`
}

func post(ctx *gin.Context) {
	// 参数
	var req postReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		internal.Handle400(ctx, err)
		return
	}

}
