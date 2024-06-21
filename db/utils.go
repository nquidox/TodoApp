package db

import (
	"encoding/json"
)

func serializeJSON(v interface{}) ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func deserializeJSON(data []byte, s interface{}) error {
	err := json.Unmarshal(data, &s)
	return err
}
