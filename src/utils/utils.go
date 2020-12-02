package utils

import (
	"encoding/json"
	"log"
)

func Struct_to_string(i interface{}) string {

	struct_json, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(struct_json)
}
