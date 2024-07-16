package exception

type UnauthorizedErrorStruct struct {
	ErrorMsg string
}

func NewUnauthorizedHandler(msg string) *UnauthorizedErrorStruct {
	return &UnauthorizedErrorStruct{
		ErrorMsg: msg,
	}
}

func (e *UnauthorizedErrorStruct) Error() string {
	return e.ErrorMsg
}
