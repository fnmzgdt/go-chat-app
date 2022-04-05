package errors

import "net/http"

type HttpError struct {
	Error error
	Message   string
	Status  int
}

func NewInternalServerError(err error, message string) *HttpError {
	return &HttpError{
		Error: err,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewBadRequestError(err error, message string) *HttpError {
	return &HttpError{
		Error: err,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}