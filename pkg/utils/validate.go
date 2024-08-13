package utils

import (
	"encoding/base64"
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"reflect"
	"scylla/pkg/helper"
	"strings"
)

func InitializeValidator() *validator.Validate {
	validate := validator.New()

	_ = validate.RegisterValidation("sliceString", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]string)
		if len(slices) == 0 {
			return false
		}
		for _, s := range slices {
			if s == "" {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("sliceInt", func(fl validator.FieldLevel) bool {
		slices := fl.Field().Interface().([]int)
		if len(slices) == 0 {
			return false
		}
		for _, val := range slices {
			if val == 0 {
				return false
			}
		}
		return true
	})

	_ = validate.RegisterValidation("equal", func(fl validator.FieldLevel) bool {
		fieldName := fl.Param()
		field := fl.Parent().FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}
		return field.Interface() == fl.Field().Interface()
	})

	_ = validate.RegisterValidation("image", func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(*multipart.FileHeader)
		if !ok || file == nil {
			return false
		}

		allowedTypes := []string{"image/jpeg", "image/png", "image/jpg"}
		if !contains(allowedTypes, file.Header.Get("Content-Type")) {
			return false
		}

		const maxSize = 5 * 1024 * 1024
		if file.Size > maxSize {
			return false
		}

		return true
	})

	_ = validate.RegisterValidation("base64Image", func(fl validator.FieldLevel) bool {
		data := fl.Field().String()

		if !strings.HasPrefix(data, "data:image/") {
			return false
		}

		base64Data := strings.SplitN(data, ",", 2)
		if len(base64Data) != 2 {
			return false
		}

		_, err := base64.StdEncoding.DecodeString(base64Data[1])
		return err == nil
	})

	_ = validate.RegisterValidation("unique", func(fl validator.FieldLevel) bool {
		return helper.IsUnique(fl)
	})

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			return name
		}
		return name
	})

	return validate
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
