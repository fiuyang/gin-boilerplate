package helper

import (
	"fmt"
	"scylla/pkg/engine"
	"strings"

	"github.com/go-playground/validator/v10"
)

func IsUnique(fl validator.FieldLevel) bool {
	return unique(fl)
}

func unique(fl validator.FieldLevel) bool {
	var count int64
	param := fl.Param()
	params := strings.Split(param, ";")
	switch len(params) {
	case 1:
		return true
	case 2:
		dField := fl.Field().String()
		engine.Instance.Table(params[0]).Where(fmt.Sprintf("%s = ?", params[1]), dField).Debug().Count(&count)
		if count > 0 {
			return false
		} else {
			return true
		}
	case 3:
		ignore := fl.Parent().FieldByName(params[2]).String()
		dField := fl.Field().String()
		engine.Instance.Table(params[0]).Where(fmt.Sprintf("%s = ?", params[1]), dField).Not(map[string]any{fmt.Sprintf("%s", strings.ToLower(params[2])): []string{ignore}}).Count(&count)
		if count > 0 {
			return false
		} else {
			return true
		}
	default:
		return true
	}
}
