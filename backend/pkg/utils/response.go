package utils

import (
	"encoding/json"
	"net/http"
)

type dataResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

type errorResponse struct {
    Error   string `json:"error"`
    Code    int    `json:"code"`
    Details string `json:"details,omitempty"`
}


func JSONResponse(w http.ResponseWriter,statusCode int,message string,data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(dataResponse{
		Success: statusCode < 400,
		Message: message,
		Data: data,
	})
}


func ErrorJSON(w http.ResponseWriter, code int, err error) {
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(errorResponse{
        Error: err.Error(),
        Code:  code,
    })
}

