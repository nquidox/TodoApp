package service

import "net/http"

type ErrorResponse struct {
	ResultCode int         `json:"resultCode"`
	ErrorCode  int         `json:"errorCode"`
	Messages   string      `json:"messages"`
	Data       interface{} `json:"data"`
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
