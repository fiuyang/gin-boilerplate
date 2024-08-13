package utils

import (
	"github.com/gin-gonic/gin"
	"scylla/dto"
)

func ResponseInterceptor(ctx *gin.Context, resp *dto.Response) {
	traceIdInf, _ := ctx.Get("trace_id")
	traceId := ""
	if traceIdInf != nil {
		traceId = traceIdInf.(string)
	}
	resp.TraceID = traceId
}
