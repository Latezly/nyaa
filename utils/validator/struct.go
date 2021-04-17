package validator

import (
	"github.com/Latezly/nyaa_go/utils/messages"
	"github.com/go-playground/validator"
)

func dateErrors(fe validator.FieldError, mes *messages.Messages) error {
	switch fe.Tag() {
	case "gt":
		return mes.AddErrorTf(fe.Field(), "error_greater_date", fe.Field())
	case "gte":
		return mes.AddErrorTf(fe.Field(), "error_greater_equal_date", fe.Field())
	case "lt":
		return mes.AddErrorTf(fe.Field(), "error_less_date", fe.Field())
	case "lte":
		return mes.AddErrorTf(fe.Field(), "error_less_equal_date", fe.Field())
	}
	return mes.AddErrorTf(fe.Field(), "error_field", fe.Field())
}
