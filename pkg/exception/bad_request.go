package exception

type BadRequestErrorStruct struct {
	ErrorMsg string
}

func NewBadRequestHandler(msg string) *BadRequestErrorStruct {
	return &BadRequestErrorStruct{
		ErrorMsg: msg,
	}
}

func (e *BadRequestErrorStruct) Error() string {
	return e.ErrorMsg
}
