package internal

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func init() {
	// val, ok := binding.Validator.Engine().(*validator.Validate)
	// if ok {
	// 	val.RegisterStructValidation(timePeroid, db.TimePeroid{})
	// }
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
