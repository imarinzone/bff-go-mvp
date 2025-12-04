package httpx

import (
	"encoding/json"
	"net/http"

	"bff-go-mvp/internal/model"
)

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes an error response matching the swagger Error schema.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	errBody := model.Error{
		Error: model.ErrorBody{
			Code:    code,
			Message: message,
		},
	}
	WriteJSON(w, status, errBody)
}


