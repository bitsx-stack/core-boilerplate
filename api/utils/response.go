package utils

import (
	"encoding/json"
	"net/http"
)


type ErroResponseType struct {
	Tatizo string `json:"tatizo"`
}

func ErrorResponse(w http.ResponseWriter, code int, message string) {
	response := ErroResponseType{
		Tatizo: message,
	}
	SuccessResponse(w, code, response)
}

func SuccessResponse(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
