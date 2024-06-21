package db

import (
	"encoding/json"
	"log"
)

func deserializeJSON(data []byte, s interface{}) {
	err := json.Unmarshal(data, &s)
	if err != nil {
		log.Println(err)
	}
}
