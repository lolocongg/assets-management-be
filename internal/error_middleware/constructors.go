package error_middleware

import "net/http"

func NewBadRequest(message string) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		CodeBadRequest,
		message,
	)
}

func NewValidationFailed(message string) *AppError {
	return NewAppError(
		http.StatusBadRequest,
		CodeValidationFailed,
		message,
	)
}

func NewUnauthorized(message string) *AppError {
	return NewAppError(
		http.StatusUnauthorized,
		CodeUnauthorized,
		message,
	)
}

func NewForbidden(message string) *AppError {
	return NewAppError(
		http.StatusForbidden,
		CodeForbidden,
		message,
	)
}

func NewNotFound(message string) *AppError {
	return NewAppError(
		http.StatusNotFound,
		CodeNotFound,
		message,
	)
}

func NewUnprocessableEntity(message string) *AppError {
	return NewAppError(
		http.StatusUnprocessableEntity,
		CodeUnprocessableEntity,
		message,
	)
}

func NewTooManyRequests(message string) *AppError {
	return NewAppError(
		http.StatusTooManyRequests,
		CodeTooManyRequests,
		message,
	)
}

func NewInternal(message string) *AppError {
	return NewAppError(
		http.StatusInternalServerError,
		CodeInternal,
		message,
	)
}

func NewServiceUnavailable(message string) *AppError {
	return NewAppError(
		http.StatusServiceUnavailable,
		CodeServiceUnavailable,
		message,
	)
}

func NewBadGateway(message string) *AppError {
	return NewAppError(
		http.StatusBadGateway,
		CodeBadGateway,
		message,
	)
}

func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}
