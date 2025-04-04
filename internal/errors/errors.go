// Package errors provides centralized error handling for the SMTP application.
package errors

import (
	"fmt"
)

// ErrorCode represents specific error types in the application
type ErrorCode string

const (
	// Mail related errors
	ErrMailValidation ErrorCode = "MAIL_VALIDATION"
	ErrMailDelivery   ErrorCode = "MAIL_DELIVERY"
	ErrMailProcessing ErrorCode = "MAIL_PROCESSING"

	// SMTP related errors
	ErrSMTPConnection ErrorCode = "SMTP_CONNECTION"
	ErrSMTPAuth       ErrorCode = "SMTP_AUTH"
	ErrSMTPDelivery   ErrorCode = "SMTP_DELIVERY"

	// Configuration errors
	ErrConfig     ErrorCode = "CONFIG"
	ErrDKIMConfig ErrorCode = "DKIM_CONFIG"

	// DNS related errors
	ErrDNSLookup ErrorCode = "DNS_LOOKUP"
	ErrMXRecord  ErrorCode = "MX_RECORD"
)

// AppError represents an application-specific error with context
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
	Context map[string]interface{}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates a new AppError
func NewError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to an AppError
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	e.Context[key] = value
	return e
}

type ConfigError struct {
	Field   string
	Message string
	Err     error
}

func (e *ConfigError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Field, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Field, e.Message)
}

// Create custom error types
type MailError struct {
	Code    string
	Message string
	Err     error
}

func (e *MailError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}
