package exception

import "fmt"

type ValidationExcel struct {
	Errors map[string][]string
}

func (e *ValidationExcel) Add(field string, row int, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string][]string)
	}
	e.Errors[field] = append(e.Errors[field], fmt.Sprintf("%s row %d", message, row))
}

func (e *ValidationExcel) Error() string {
	return "validation errors"
}
