package exception

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"scylla/entity"
	"strings"
	"unicode"
)

func ErrorHandlers(ctx *gin.Context, err interface{}) {
	if notFoundError(ctx, err) {
		return
	} else if validationError(ctx, err) {
		return
	} else if badRequestError(ctx, err) {
		return
	} else if unauthorizedError(ctx, err) {
		return
	} else if validationExcelError(ctx, err) {
		return
	} else if validationExcel(ctx, err) {
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
				report[fieldName] = fmt.Sprintf("%s has already been taken %s", fieldName)
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
				report[fieldName] = fmt.Sprintf("%s value must be char and numeric", fieldName, e.Param())
			case "notEmptyStringSlice":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is string", fieldName)
			case "dive":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty", fieldName)
			case "date":
				report[fieldName] = fmt.Sprintf("%s value must be date (yyyy-mm-dd)", fieldName)
			case "notEmptyIntSlice":
				report[fieldName] = fmt.Sprintf("%s value ​​in the array cannot be empty is int", fieldName)
			case "isInt":
				report[fieldName] = fmt.Sprintf("%s value must be of type int", fieldName)
			case "isString":
				report[fieldName] = fmt.Sprintf("%s value must be of type string", fieldName)
			}
		}

		ctx.JSON(http.StatusBadRequest, entity.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  report,
			TraceID: uuid.New().String(),
		})

		return true
	}
	return false
}

func badRequestError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*BadRequestErrorStruct)
	if ok {
		ctx.JSON(http.StatusBadRequest, entity.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Error(),
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func unauthorizedError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*UnauthorizedErrorStruct)
	if ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, entity.Error{
			Code:    http.StatusUnauthorized,
			Status:  "UNAUTHORIZED",
			Errors:  exception.Error(),
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func notFoundError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*NotFoundErrorStruct)
	if ok {
		ctx.JSON(http.StatusNotFound, entity.Error{
			Code:    http.StatusNotFound,
			Status:  "NOT FOUND",
			Errors:  exception.Error,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func internalServerError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*InternalServerErrorStruct)
	if ok {
		ctx.JSON(http.StatusInternalServerError, entity.Error{
			Code:    http.StatusInternalServerError,
			Status:  "INTERNAL SERVER ERROR",
			Errors:  exception.Error(),
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func validationExcelError(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*NewValidationExcelError)
	if ok {
		ctx.JSON(http.StatusBadRequest, entity.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Errors,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}

func validationExcel(ctx *gin.Context, err interface{}) bool {
	exception, ok := err.(*ValidationExcel)
	if ok {
		ctx.JSON(http.StatusBadRequest, entity.Error{
			Code:    http.StatusBadRequest,
			Status:  "BAD REQUEST",
			Errors:  exception.Errors,
			TraceID: uuid.New().String(),
		})
		return true
	}
	return false
}
