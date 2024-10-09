package lexer

import (
	"fmt"
)

type ParseError struct {
	Message string
}

func (e *ParseError) Error() string {
	return e.Message
}

func NewUnexpectedCharacterError(ch byte) error {
	return &ParseError{Message: fmt.Sprintf("Unexpected charatcer: %q", ch)}
}

func NewUnexpectedTokenError(got Token, expected interface{}) error {
	return &ParseError{Message: fmt.Sprintf("Unexpected token: %v, expected: %v", got, expected)}
}
