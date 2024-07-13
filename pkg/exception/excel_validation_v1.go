package exception

import "fmt"

type NewValidationExcelError struct {
	Errors map[string][]string
}

func (e *NewValidationExcelError) Add(field string, row int, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string][]string)
	}
	if _, ok := e.Errors[field]; !ok {
		e.Errors[field] = []string{}
	}
	e.Errors[field] = append(e.Errors[field], fmt.Sprintf("%s row %d", message, row))
}

func (e *NewValidationExcelError) Error() string {
	return "validation errors"
}
