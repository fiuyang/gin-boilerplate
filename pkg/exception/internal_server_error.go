package exception

type InternalServerErrorStruct struct {
	ErrorMsg string
}

func NewInternalServerErrorHandler(msg string) *InternalServerErrorStruct {
	return &InternalServerErrorStruct{
		ErrorMsg: msg,
	}
}

func (e *InternalServerErrorStruct) Error() string {
	return e.ErrorMsg
}
