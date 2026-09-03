package httpx

import (
	"encoding/json"
	"net/http"
)

type Code string

const (
	CodeInvalidID     Code = "invalid_id"
	CodeInternalError Code = "Internal_error"
)

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status int, message string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{Error: errorPayload{
		Code:    code,
		Message: message,
	}})
}
