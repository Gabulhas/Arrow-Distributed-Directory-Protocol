package utils

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
)

func StructToString(i interface{}) string {

	struct_json, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(struct_json)
}

func RandomRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomSleep(duration time.Duration, min int, max int) {

	time.Sleep(duration * time.Duration(RandomRange(min, max)))
}
