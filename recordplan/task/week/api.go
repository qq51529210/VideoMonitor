package week

import "github.com/qq51529210/video-monitor/recordplan/db"

// Run 启动后台检查
func Run(checkInterval, concurrency, apiTimeout int) error {
	return _checker.init(checkInterval, concurrency, apiTimeout)
}

// Reload 重新加载，在添加和修改数据库后调用
func Reload(id string) {
	_checker.Add(1)
	go _checker.loadRoutine(id)
}

// DeleteStream 移除所有计划的 stream ，再删除 stream 后调用
func DeleteStream(stream string) {
	_checker.removeStream(stream)
}

// IsRecording 返回指定 id 的录像是否在录像时间段。
// 再查询出关联的 stream ，就可以知道 stream 是否需要录像。
// 如果 id 不存在，返回 ErrNotFound
func IsRecording(id string) (bool, error) {
	_checker.RLock()
	plan := _checker.weekplan[id]
	_checker.RUnlock()
	if plan == nil {
		return false, db.ErrNotFound
	}
	return plan.isRecording, nil
}
