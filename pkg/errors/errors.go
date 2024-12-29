package errors

import "fmt"

// ParseError represents an error that occurs during parsing.
type ParseError struct {
	Line    int
	Column  int
	Message string
}

// Error formats the ParseError into a readable string.
func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at line %d, column %d : %s", e.Line, e.Column, e.Message)
}

// NewParseError creates a new instance of ParseError.
func NewParseError(line, column int, message string) error {
	return &ParseError{
		Line:    line,
		Column:  column,
		Message: message,
	}
}

// ValidationError represents an error that occurs during validation.
type ValidationError struct {
	Message string
}

// Error formats the ValidationError into a readable string.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new instance of ValidationError.
func NewValidationError(message string) error {
	return &ValidationError{
		Message: message,
	}
}
