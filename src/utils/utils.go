package utils

import (
	"encoding/json"
	"log"
	"math/rand"
)

func StructToString(i interface{}) string {

	structJson, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(structJson)
}

func RandomRange(min, max int) int {
	return rand.Intn(max - min + 1) + min
}

