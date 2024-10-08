package errors

import "fmt"

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

func NewUnexpectedCharacterError(ch byte) error {
	return &ParseError{Message: fmt.Sprintf("Unexpected charatcer: %q", ch)}
}