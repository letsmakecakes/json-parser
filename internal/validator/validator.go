package validator

import (
	"github.com/letsmakecakes/jsonparser/pkg/errors"
	"regexp"
	"unicode"
)

// Validator validates JSON elements, such as depth, strings, and numbers.
type Validator struct {
	maxDepth     int
	currentDepth int
}

// New creates a new Validator with a specified maximum nesting depth.
func New(maxDepth int) *Validator {
	return &Validator{
		maxDepth: maxDepth,
	}
}

// ValidateDepth increments the current depth and checks if it exceeds the maximum depth.
func (v *Validator) ValidateDepth() error {
	v.currentDepth++
	if v.currentDepth > v.maxDepth {
		return errors.NewValidationError("exceeded maximum nesting depth")
	}
	return nil
}

// ExitDepth decrements the current depth, ensuring it does not go below zero.
func (v *Validator) ExitDepth() {
	if v.currentDepth > 0 {
		v.currentDepth--
	}
}

// ValidateString checks if a string contains invalid control characters.
func (v *Validator) ValidateString(s string) error {
	for _, r := range s {
		if r < 0x20 && !unicode.IsSpace(r) {
			return errors.NewValidationError("invalid control character in string")
		}
	}
	return nil
}

// ValidateNumber validates if a string represents a valid number in JSON format.
func (v *Validator) ValidateNumber(numStr string) error {
	// JSON number regex based on ECMA-404 specification
	numberRegex := `^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`

	// compile and match the regex
	re := regexp.MustCompile(numberRegex)
	if !re.MatchString(numStr) {
		return errors.NewValidationError("invalid number format")
	}
	return nil
}
