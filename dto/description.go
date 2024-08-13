package dto

type JsonBadRequest struct {
	Code    int                    `json:"code" example:"400"`
	Status  string                 `json:"status" example:"BAD REQUEST"`
	Errors  map[string]interface{} `json:"errors" swaggertype:"object,string" example:"username:username is required,email:email is required"`
	TraceID string                 `json:"trace_id" example:"dedc5250-5c20-48c9-9383-fac3ccff2679"`
}

type JsonSuccess struct {
	Code    int         `json:"code" example:"200"`
	Status  string      `json:"status" example:"OK"`
	Message string      `json:"message,omitempty" example:"Success"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id" example:"dedc5250-5c20-48c9-9383-fac3ccff2679"`
}

type JsonCreated struct {
	Code    int         `json:"code" example:"201"`
	Status  string      `json:"status" example:"CREATED"`
	Message string      `json:"message,omitempty" example:"Created"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id" example:"dedc5250-5c20-48c9-9383-fac3ccff2679"`
}

type JsonInternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Status  string `json:"status" example:"INTERNAL SERVER ERROR"`
	Errors  string `json:"errors,omitempty" example:"error database or third party"`
	TraceID string `json:"trace_id" example:"dedc5250-5c20-48c9-9383-fac3ccff2679"`
}

type JsonNotFound struct {
	Code    int    `json:"code" example:"404"`
	Status  string `json:"status" example:"NOT FOUND"`
	Errors  string `json:"errors,omitempty" example:"record not found"`
	TraceID string `json:"trace_id" example:"dedc5250-5c20-48c9-9383-fac3ccff2679"`
}
