package internal

import (
	"fmt"
	"time"

	"github.com/qq51529210/video-monitor/recordplan/db"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterStructValidation(timePeroid, db.TimePeroid{})
	}
}

func timePeroid(fn validator.StructLevel) {
	var start, end time.Time
	var err error
	ok := true
	//
	model := fn.Current().Interface().(db.TimePeroid)
	if model.Start == "" {
		fn.ReportError(model.Start, "start", "Start", "timePeroidEmpty", "")
		ok = false
	} else {
		start, err = time.Parse(db.TimePeroidFormat, model.Start)
		if err != nil {
			fn.ReportError(model.Start, "start", "Start", "timePeroidTimeFormat", "")
			ok = false
		}
	}
	if model.End == "" {
		fn.ReportError(model.End, "end", "End", "timePeroidEmpty", "")
		ok = false
	} else {
		end, err = time.Parse(db.TimePeroidFormat, model.End)
		if err != nil {
			fn.ReportError(model.End, "end", "End", "timePeroidTimeFormat", "")
			ok = false
		}
	}
	if ok && start.After(end) {
		fn.ReportError(model.End, "end", "TimePeroid", "timePeroidError", "")
	}
}

func FormatError(err error) any {
	switch v := err.(type) {
	case validator.ValidationErrors:
		var errs []string
		for _, es := range v {
			switch es.Tag() {
			case "timePeroidEmpty":
				errs = append(errs, fmt.Sprintf("[%s] must not be empty", es.Field()))
			case "timePeroidTimeFormat":
				errs = append(errs, fmt.Sprintf("[%s] error format", es.Field()))
			case "timePeroidError":
				errs = append(errs, "[start] must less than [end]")
			default:
				errs = append(errs, es.Error())
			}
		}
		return errs
	case *validator.InvalidValidationError:
		return v.Error()
	default:
		return err.Error()
	}
}
