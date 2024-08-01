package service

import "net/http"

type ErrorResponse struct {
	ResultCode int         `json:"resultCode" extensions:"x-order=1"`
	ErrorCode  int         `json:"errorCode" extensions:"x-order=2"`
	Messages   string      `json:"messages" extensions:"x-order=3"`
	Data       interface{} `json:"data" extensions:"x-order=4"`
}

func ServerResponse(w http.ResponseWriter, dataInterface interface{}) {
	bytes, err := SerializeJSON(dataInterface)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
