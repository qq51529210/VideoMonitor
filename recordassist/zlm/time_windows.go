package zlm

import (
	"io/fs"
	"syscall"
	"time"
)

func getTime(info fs.FileInfo) (time.Time, time.Time) {
	t := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(t.CreationTime.Nanoseconds()/int64(time.Second), 0),
		time.Unix(t.LastWriteTime.Nanoseconds()/int64(time.Second), 0)
}
