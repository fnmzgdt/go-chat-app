package errors

import "net/http"

type RestError struct {
	Message string
	Status  int
	Error   string
}

func NewInternalServerError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error: "Internal Server Error",
	}
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Message: message,
		Status:  http.StatusBadRequest,
		Error: "Bad Request",
	}
}