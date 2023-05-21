package zlm

import (
	"io/fs"
	"syscall"
	"time"

	"github.com/qq51529210/log"
)

func getTime(info fs.FileInfo) (time.Time, time.Time) {
	t := info.Sys().(*syscall.Stat_t)
	return time.Unix(t.Ctim.Sec, 0),
		time.Unix(t.Mtim.Sec, 0)

}
