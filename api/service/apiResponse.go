package service

import (
	"fmt"
	"net/http"
)

type DefaultResponse struct {
	ResultCode int    `json:"resultCode" extensions:"x-order=1"`
	HttpCode   int    `json:"httpCode" extensions:"x-order=2"`
	Messages   string `json:"messages" extensions:"x-order=3"`
	Data       any    `json:"data" extensions:"x-order=4"`
}

type errorResponse struct {
	ResultCode int    `json:"resultCode" extensions:"x-order=1"`
	HttpCode   int    `json:"httpCode" extensions:"x-order=2"`
	Messages   string `json:"messages" extensions:"x-order=3"`
	Data       any    `json:"data" extensions:"x-order=4"`
}

func serverResponse(w http.ResponseWriter, dataInterface interface{}) {
	bytes, err := SerializeJSON(dataInterface)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func OkResponse(w http.ResponseWriter, data interface{}) {
	serverResponse(w, data)
}

func BadRequestResponse(w http.ResponseWriter, errType string, errMsg error) {
	serverResponse(w, errorResponse{
		ResultCode: 1,
		HttpCode:   http.StatusBadRequest,
		Messages:   "Bad Request",
		Data:       fmt.Sprintf("%s: %s", errType, errMsg),
	})
}

func UnauthorizedResponse(w http.ResponseWriter, msg any) {
	serverResponse(w, errorResponse{
		ResultCode: 1,
		HttpCode:   http.StatusUnauthorized,
		Messages:   "Unauthorized",
		Data:       msg,
	})
}

func UnprocessableEntity(w http.ResponseWriter, errType string, errMsg error) {
	serverResponse(w, errorResponse{
		ResultCode: 1,
		HttpCode:   http.StatusUnprocessableEntity,
		Messages:   "Unprocessable Entity",
		Data:       fmt.Sprintf("%s: %s", errType, errMsg),
	})
}

func InternalServerErrorResponse(w http.ResponseWriter, errType string, errMsg error) {
	serverResponse(w, errorResponse{
		ResultCode: 1,
		HttpCode:   http.StatusInternalServerError,
		Messages:   "Internal Server Error",
		Data:       fmt.Sprintf("%s: %s", errType, errMsg),
	})
}
