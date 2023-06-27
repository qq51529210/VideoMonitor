package task

import (
	"time"

	"github.com/qq51529210/log"
	"github.com/qq51529210/video-monitor/onvif/db"
	"github.com/qq51529210/video-monitor/onvif/ovf"
)

// updateDeviceKeepaliveRoutine 启动一个协程
// 检查并更新数据库中所有的设备的心跳
func updateDeviceKeepaliveRoutine(dur time.Duration) {
	timer := time.NewTimer(0)
	defer func() {
		// 抓异常
		log.Recover(recover())
		// 计时器
		timer.Stop()
	}()
	for {
		// 数据库
		models, err := db.DeviceDA.All()
		if err != nil {
			log.Errorf("[keepalive] db error: %s", err.Error())
			timer.Reset(time.Second)
			<-timer.C
			continue
		}
		// dur 平分到每一个的时间，最小 1 毫秒
		d := time.Duration(float64(dur) / float64(len(models)))
		if d < time.Millisecond {
			d = time.Millisecond
		}
		// 循环检查
		for _, model := range models {
			timer.Reset(time.Millisecond)
			<-timer.C
			updateDeviceKeepalive(model)
		}
	}
}

// updateDeviceKeepalive 查询
func updateDeviceKeepalive(model *db.Device) {
	updateDB := false
	online := *model.Online
	// 调用 api
	_, err := ovf.GetDeviceTime(model)
	if err != nil {
		log.Errorf("[keepalive] query device %d time error: %s", model.ID, err.Error())
		// 在线 -> 离线
		if online == db.True {
			online = db.False
			updateDB = true
		}
	} else {
		// 离线 -> 在线
		if online == db.False {
			online = db.True
			updateDB = true
		}
	}
	// 更新
	if updateDB {
		err = db.UpdateDeviceOnline(model.ID, online)
		if err != nil {
			log.Errorf("[keepalive] update db device %d online error: %s", model.ID, err.Error())
		}
	}
}
