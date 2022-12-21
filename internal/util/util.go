package util

import (
	"encoding/json"
	"log"
)

func StructToString(data interface{}) string {
	d, err := json.Marshal(data)
	if err != nil {
		log.Println("Cannot Convert Struct -> String ", err)
		return ""
	}
	return string(d)
}
