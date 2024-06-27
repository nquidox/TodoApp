package todoList

import (
	"net/http"
	"todoApp/api/service"
)

type RequestTitle struct {
	Title string `json:"title"`
}

type ErrorResponse struct {
	ResultCode int    `json:"resultCode"`
	Messages   string `json:"messages"`
	Data       string `json:"data"`
}

func serverResponse(w http.ResponseWriter, dataInterface interface{}) {
	bytes, err := service.SerializeJSON(dataInterface)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
