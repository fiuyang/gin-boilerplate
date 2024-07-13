package helper

import (
	"encoding/json"
)

func Automapper(objOrigin interface{}, objDestination interface{}) {
	jsonOrigin := StructToJson(objOrigin)
	json.Unmarshal([]byte(jsonOrigin), objDestination)
}
