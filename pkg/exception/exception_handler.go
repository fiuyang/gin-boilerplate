package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"scylla/dto"
	"strings"
	"unicode"
)

func ExceptionHandlers(ctx *gin.Context, err interface{}) {
	if notFoundError(ctx, err) {
		return
	} else if validationError(ctx, err) {
		return
	} else if badRequestError(ctx, err) {
		return
	} else if unauthorizedError(ctx, err) {
		return
	} else if excelValidationError(ctx, err) {
		return
	} else if excelValidation(ctx, err) {
		return
	} else {
		internalServerError(ctx, err)
		return
	}
}

func validationError(ctx *gin.Context, err interface{}) bool {

	if castedObject, ok := err.(validator.ValidationErrors); ok {
		report := make(map[string]string)
		var fieldName string

		for _, e := range castedObject {
			if len(e.Namespace()) > 0 && unicode.IsUpper(rune(e.Namespace()[0])) {
				dotIndex := strings.Index(e.Namespace(), ".")
				if dotIndex != -1 {
					fieldName = e.Namespace()[dotIndex+1:]
				}
			} else {
				fieldName = e.Field()
			}
			switch e.Tag() {
			case "required":
				report[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				report[fieldName] = fmt.Sprintf("%s is not valid email", fieldName)
			case "gte":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "lte":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "unique":
				report[fieldName] = fmt.Sprintf("%s has already been taken", fieldName)
			case "max":
				report[fieldName] = fmt.Sprintf("%s value must be lower than %s", fieldName, e.Param())
			case "min":
				report[fieldName] = fmt.Sprintf("%s value must be greater than %s", fieldName, e.Param())
			case "numeric":
				report[fieldName] = fmt.Sprintf("%s value must be numeric", fieldName)
			case "number":
				report[fieldName] = fmt.Sprintf("%s value must be number", fieldName)
			case "oneof":
				report[fieldName] = fmt.Sprintf("%s value must be %s", fieldName, e.Param())
			case "len":
				report[fieldName] = fmt.Sprintf("%s value must be exactly %s characters long", fieldName, e.Param())
			case "alphanum":
				report[fieldName] = fmt.Sprintf("%s value must be char and numeric %s", fieldName, e.Param())
			case "sliceString":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is string", fieldName)
			case "dive":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty", fieldName)
			case "datetime":
				report[fieldName] = fmt.Sprintf("%s value must be date (yyyy-mm-dd)", fieldName)
			case "required_if":
				report[fieldName] = fmt.Sprintf("%s must be filled in if %s", fieldName, e.Param())
			case "sliceInt":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is int", fieldName)
			case "equal":
				report[fieldName] = fmt.Sprintf("%s and %s do not match do not match", fieldName, e.Param())
			case "image":
				report[fieldName] = fmt.Sprintf("%s file must be of type jpg, jpeg, png", fieldName)
			case "base64Image":
				report[fieldName] = fmt.Sprintf("%s value must be base64 encoded image", fieldName)
			}
		}

		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  report,
			TraceID: traceID.(string),
		})

		return true
	}
	return false
}

func badRequestError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*BadRequestErrorStruct)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Error(),
			TraceID: traceID.(string),
		})

		return true
	}
	return false
}

func unauthorizedError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*UnauthorizedErrorStruct)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Error{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Errors:  exception.Error(),
			TraceID: traceID.(string),
		})
		return true
	}
	return false
}

func notFoundError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*NotFoundErrorStruct)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusNotFound, dto.Error{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Errors:  exception.Error,
			TraceID: traceID.(string),
		})
		return true
	}
	return false
}

func internalServerError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*InternalServerErrorStruct)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusInternalServerError, dto.Error{
			Code:    http.StatusInternalServerError,
			Status:  "INTERNAL SERVER ERROR",
			Errors:  exception.Error(),
			TraceID: traceID.(string),
		})
		return true
	}
	return false
}

func excelValidationError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*NewExcelValidationError)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Errors,
			TraceID: traceID.(string),
		})
		return true
	}
	return false
}

func excelValidation(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*ExcelValidation)
	if ok {
		traceID, _ := ctx.Get("trace_id")
		ctx.JSON(http.StatusBadRequest, dto.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Errors,
			TraceID: traceID.(string),
		})
		return true
	}
	return false
}
