package api

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool `json:"success"`
	Code    int  `json:"code"`
	Data    T    `json:"data"`
}

type IDResp struct {
	ID int64 `json:"id"`
}

type Message struct {
	Message string `json:"message"`
}

func writeMessageResponse(w http.ResponseWriter, message string, code int, success bool) {
	resp := Response[Message]{
		Success: success,
		Code:    code,
		Data:    Message{Message: message},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		return
	}
}

var (
	RequestUnauthorisedHandler = func(w http.ResponseWriter, message string) {
		writeMessageResponse(w, message, http.StatusUnauthorized, false)
	}

	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeMessageResponse(w, err.Error(), http.StatusBadRequest, false)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeMessageResponse(w, "Internal error", http.StatusInternalServerError, false)
	}
	RequestSuccessHandler = func(w http.ResponseWriter, message string) {
		writeMessageResponse(w, message, http.StatusOK, true)
	}
)
