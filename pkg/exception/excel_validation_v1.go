package exception

import "fmt"

type NewExcelValidationError struct {
	Errors map[string][]string
}

func (e *NewExcelValidationError) AddHandler(field string, row int, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string][]string)
	}
	if _, ok := e.Errors[field]; !ok {
		e.Errors[field] = []string{}
	}
	e.Errors[field] = append(e.Errors[field], fmt.Sprintf("%s row %d", message, row))
}

func (e *NewExcelValidationError) Error() string {
	return "validation errors"
}
