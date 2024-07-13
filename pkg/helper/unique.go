package helper

import (
	"reflect"
	"scylla/model"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var modelMap = map[string]reflect.Type{
	"customers": reflect.TypeOf(model.Customer{}),
	"users":     reflect.TypeOf(model.User{}),
}

func ValidateUnique(db *gorm.DB, fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	tableName := getModelFromTag(fl)

	exists := UniqueExistsInTable(db, value, tableName)

	return !exists
}

func UniqueExistsInTable(db *gorm.DB, value interface{}, tableName string) bool {
	parts := strings.Split(tableName, ";")
	modelName := parts[0]
	columnName := parts[1]

	modelType, ok := modelMap[modelName]
	if !ok {
		return false
	}

	modelInstance := reflect.New(modelType).Interface()

	var err error
	if len(parts) > 2 {
		err = db.Table(modelName).Where(columnName+" = ? AND "+parts[2], value).First(modelInstance).Error
	} else {
		err = db.Table(modelName).Where(columnName+" = ?", value).First(modelInstance).Error
	}
	if err != nil {
		return false
	}
	return true
}

func getModelFromTag(fl validator.FieldLevel) string {
	// Assuming 'validate' tag is in the format "unique=tableName;columnName;columnID" columnID is optional when update data
	validateTag := fl.Param()

	parts := strings.Split(validateTag, "=")
	if len(parts) >= 2 {
		return parts[1]
	}

	return validateTag
}
