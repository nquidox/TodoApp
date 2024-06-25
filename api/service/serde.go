package service

import "encoding/json"

func SerializeJSON(v interface{}) ([]byte, error) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func DeserializeJSON(data []byte, s interface{}) error {
	err := json.Unmarshal(data, &s)
	return err
}
