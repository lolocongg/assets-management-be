package error_middleware

const (
	// HTTP status 4xx
	CodeBadRequest          = "BAD_REQUEST"
	CodeValidationFailed    = "VALIDATION_FAILED"
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeForbidden           = "FORBIDDEN"
	CodeNotFound            = "NOT_FOUND"
	CodeUnprocessableEntity = "UNPROCESSABLE_ENTITY"
	CodeTooManyRequests     = "TOO_MANY_REQUESTS"

	// HTTP status 5xx
	CodeInternal           = "INTERNAL_SERVER_ERROR"
	CodeServiceUnavailable = "SERVICE_UNAVAILABLE"
	CodeBadGateway         = "BAD_GATEWAY"
)

type AppError struct {
	HTTPStatus int            `json:"http_status"`
	Code       string         `json:"code"`
	Message    string         `json:"message"`
	Details    map[string]any `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string) *AppError {
	return &AppError{
		HTTPStatus: status,
		Code:       code,
		Message:    message,
	}
}
