package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// APIError represents a structured error message.
type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"` // Changed to interface{} to support different data types
}

// NewAPIError creates a new APIError instance.
func NewAPIError(code int, msg string, detail interface{}) *APIError {
	return &APIError{Code: code, Message: msg, Detail: detail}
}

// HandleValidationError handles validation errors and formats them into an APIError.
func HandleValidationError(c echo.Context, validationErrors map[string]any) error {
	apiError := NewAPIError(http.StatusBadRequest, "Validation failed", validationErrors)
	return c.JSON(apiError.Code, apiError)
}

// HandleError handles generic API errors.
func HandleError(c echo.Context, err *APIError, detail ...string) error {
	// Check if a detail was provided and append it to the error's detail field
	if len(detail) > 0 && detail[0] != "" {
		err.Detail = detail[0]
	}
	return c.JSON(err.Code, err)
}
