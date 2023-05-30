package log

import (
	"os"
	"runtime"
	"time"
)

// FormatTime format is "2006-01-02 15:04:05.000000"
func FormatTime(log *Log) {
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// Date
	log.WriteIntLeftAlign(year, 4)
	log.b = append(log.b, '-')
	log.WriteIntRightAlign(int(month), 2)
	log.b = append(log.b, '-')
	log.WriteIntRightAlign(day, 2)
	log.b = append(log.b, ' ')
	// Time
	log.WriteIntRightAlign(hour, 2)
	log.b = append(log.b, ':')
	log.WriteIntRightAlign(minute, 2)
	log.b = append(log.b, ':')
	log.WriteIntRightAlign(second, 2)
	// Nanosecond
	log.b = append(log.b, '.')
	log.WriteInt(t.Nanosecond())
}

// Header 用于格式化日志头
type Header interface {
	Time(log *Log)
	Stack(log *Log, depth int)
}

// DefaultHeader 实现 Header 接口
// 格式 2006-01-02 15:04:05.000000
type DefaultHeader struct {
}

func (th *DefaultHeader) Time(log *Log) {
	FormatTime(log)
}

func (th *DefaultHeader) Stack(log *Log, depth int) {
}

// FileNameHeader 实现 Header 接口
// 格式 2006-01-02 15:04:05.000000 [fileName:fileLine]
type FileNameHeader struct {
}

func (th *FileNameHeader) Time(log *Log) {
	FormatTime(log)
}

func (th *FileNameHeader) Stack(log *Log, depth int) {
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if os.IsPathSeparator(path[i]) {
				path = path[i+1:]
				break
			}
		}
	}
	log.b = append(log.b, ' ')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.WriteInt(line)
}

// FilePathHeader 实现 Header 接口
// 格式 2006-01-02 15:04:05.000000 [filePath:fileLine]
type FilePathHeader struct {
}

func (th *FilePathHeader) Time(log *Log) {
	FormatTime(log)
}

func (th *FilePathHeader) Stack(log *Log, depth int) {
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	log.b = append(log.b, ' ')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.WriteInt(line)
}
