package exception

type NotFoundErrorStruct struct {
	ErrorMsg string
}

func NewNotFoundHandler(msg string) *NotFoundErrorStruct {
	return &NotFoundErrorStruct{
		ErrorMsg: msg,
	}
}

func (e *NotFoundErrorStruct) Error() string {
	return e.ErrorMsg
}
