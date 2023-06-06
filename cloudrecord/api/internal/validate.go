package internal

// import (
// 	"github.com/qq51529210/video-monitor/cloudrecord/db"
// 	"fmt"
// 	"time"

// 	"github.com/gin-gonic/gin/binding"
// 	"github.com/go-playground/validator/v10"
// )

// func init() {
// 	val, ok := binding.Validator.Engine().(*validator.Validate)
// 	if ok {
// 		val.RegisterStructValidation(timePeroid, db.TimePeroid{})
// 	}
// }

// func timePeroid(fn validator.StructLevel) {
// 	var start, end time.Time
// 	var err error
// 	ok := true
// 	//
// 	model := fn.Current().Interface().(db.TimePeroid)
// 	if model.Start == "" {
// 		fn.ReportError(model.Start, "start", "Start", "timePeroidEmpty", "")
// 		ok = false
// 	} else {
// 		start, err = time.Parse(db.TimePeroidFormat, model.Start)
// 		if err != nil {
// 			fn.ReportError(model.Start, "start", "Start", "timePeroidTimeFormat", "")
// 			ok = false
// 		}
// 	}
// 	if model.End == "" {
// 		fn.ReportError(model.End, "end", "End", "timePeroidEmpty", "")
// 		ok = false
// 	} else {
// 		end, err = time.Parse(db.TimePeroidFormat, model.End)
// 		if err != nil {
// 			fn.ReportError(model.End, "end", "End", "timePeroidTimeFormat", "")
// 			ok = false
// 		}
// 	}
// 	if ok && start.After(end) {
// 		fn.ReportError(model.End, "end", "TimePeroid", "timePeroidError", "")
// 	}
// }
