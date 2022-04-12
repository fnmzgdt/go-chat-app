package errors

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	Error   error
	Message string
	Status  int
}

func NewInternalServerError(err error, message string) *HttpError {
	return &HttpError{
		Error:   err,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewBadRequestError(err error, message string) *HttpError {
	return &HttpError{
		Error:   err,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, "error", message)
}

func RespondWithJSON(w http.ResponseWriter, code int, messageType string, message interface{}) {
	payload := map[string]interface{}{messageType: message}
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
