package utils

import (
	"github.com/gin-gonic/gin"
	"scylla/entity"
)

func ResponseInterceptor(ctx *gin.Context, resp *entity.Response) {
	traceIdInf, _ := ctx.Get("trace_id")
	traceId := ""
	if traceIdInf != nil {
		traceId = traceIdInf.(string)
	}
	resp.TraceID = traceId
}
