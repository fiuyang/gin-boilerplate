package exception

import "fmt"

type ExcelValidation struct {
	Errors map[string][]string
}

func (e *ExcelValidation) AddHandler(field string, row int, message string) {
	if e.Errors == nil {
		e.Errors = make(map[string][]string)
	}
	e.Errors[field] = append(e.Errors[field], fmt.Sprintf("%s row %d", message, row))
}

func (e *ExcelValidation) Error() string {
	return "validation errors"
}
