// Package pmail provides core mail handling structures and functionality.
package pmail

import (
	"github.com/mjl-/mox/smtp"
	"github.com/mjl-/mox/smtpclient"
	"github.com/stlimtat/remiges-smtp/internal/errors"
)

// Mail represents an email message with all its components and metadata.
// It provides a structured way to handle email content and delivery information.
type Mail struct {
	// Body contains the raw message body
	Body []byte `validate:"required" json:"body"`

	// Headers contains the raw message headers
	Headers []byte `validate:"required" json:"headers"`

	// HeadersMap provides quick access to header values
	// Key is the header name (case-sensitive)
	// Value is the header value in raw bytes
	HeadersMap map[string][]byte `json:"headers_map"`

	// ContentType specifies the MIME type of the message
	ContentType []byte `validate:"required" json:"content_type"`

	// FinalBody contains the complete message after processing
	FinalBody []byte `json:"final_body"`

	// From specifies the sender's email address
	From smtp.Address `validate:"required" json:"from"`

	// Metadata stores additional processing information
	Metadata map[string][]byte `json:"metadata"`

	// MsgID contains the unique message identifier
	MsgID []byte `validate:"required" json:"msg_id"`

	// MsgPrefix contains any prefix data for the message
	MsgPrefix []byte `json:"msg_prefix"`

	// Subject contains the email subject
	Subject []byte `validate:"required" json:"subject"`

	// To contains the list of recipients
	To []smtp.Address `validate:"required,min=1" json:"to"`
}

// Validate performs validation of the Mail structure
func (m *Mail) Validate() error {
	if m == nil {
		return errors.NewError(errors.ErrMailValidation, "mail cannot be nil", nil)
	}

	if len(m.Body) == 0 {
		return errors.NewError(errors.ErrMailValidation, "body cannot be empty", nil)
	}

	if len(m.Headers) == 0 {
		return errors.NewError(errors.ErrMailValidation, "headers cannot be empty", nil)
	}

	if len(m.ContentType) == 0 {
		return errors.NewError(errors.ErrMailValidation, "content type cannot be empty", nil)
	}

	if m.From.String() == "" {
		return errors.NewError(errors.ErrMailValidation, "from address cannot be empty", nil)
	}

	if len(m.To) == 0 {
		return errors.NewError(errors.ErrMailValidation, "to addresses cannot be empty", nil)
	}

	return nil
}

// SetHeader safely sets a header value in the HeadersMap
func (m *Mail) SetHeader(name string, value []byte) {
	if m.HeadersMap == nil {
		m.HeadersMap = make(map[string][]byte)
	}
	m.HeadersMap[name] = value
}

// GetHeader safely retrieves a header value from HeadersMap
func (m *Mail) GetHeader(name string) ([]byte, bool) {
	if m.HeadersMap == nil {
		return nil, false
	}
	value, exists := m.HeadersMap[name]
	return value, exists
}

// Response wraps the SMTP client response with additional functionality
type Response struct {
	smtpclient.Response
}

type HeaderMap map[string][]byte
type MetadataMap map[string][]byte
