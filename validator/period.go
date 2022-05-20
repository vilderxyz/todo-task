package valid

import "github.com/go-playground/validator/v10"

const (
	TODAY    = "today"
	TOMORROW = "tomorrow"
	WEEK     = "week"
	ALL      = ""
)

// Custom validator that returns false when string is not one of
//
// [ "today" , "tomorrow" , "week" , ""]
var ValidPeriod validator.Func = func(fl validator.FieldLevel) bool {
	if period, ok := fl.Field().Interface().(string); ok {
		switch period {
		case TODAY, TOMORROW, WEEK, ALL:
			return true
		}
	}
	return false
}
