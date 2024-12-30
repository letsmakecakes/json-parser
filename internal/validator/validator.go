package validator

import (
	"github.com/letsmakecakes/jsonparser/internal/parser"
	"github.com/letsmakecakes/jsonparser/pkg/errors"
	"regexp"
	"strconv"
	"unicode"
)

// Validator validates JSON elements, such as depth, strings, and numbers.
type Validator struct {
	maxDepth int
}

// New creates a new Validator with a specified maximum nesting depth.
func New(maxDepth int) *Validator {
	return &Validator{
		maxDepth: maxDepth,
	}
}

// Validate traverses the AST and validates the JSON structure.
func (v *Validator) Validate(node parser.Node) error {
	return v.validateNode(node, 0)
}

// validateNode recursively validates a node and its children, checking depth and other constraints.
func (v *Validator) validateNode(node parser.Node, depth int) error {
	if depth > v.maxDepth {
		return errors.NewValidationError("exceeded maximum nesting depth")
	}

	switch n := node.(type) {
	case *parser.ObjectValue:
		for key, value := range n.Pairs {
			if err := v.ValidateString(key); err != nil {
				return err
			}
			if err := v.validateNode(value, depth+1); err != nil {
				return err
			}
		}
	case *parser.ArrayValue:
		for _, elem := range n.Elements {
			if err := v.validateNode(elem, depth+1); err != nil {
				return err
			}
		}
	case *parser.StringValue:
		return v.ValidateString(n.Value)
	case *parser.NumberValue:
		return v.ValidateNumber(n.Value)
	case *parser.BooleanValue, *parser.NullValue:
	// No specific validation needed for boolean or null
	default:
		return errors.NewValidationError("unknown node type")
	}

	return nil
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
func (v *Validator) ValidateNumber(num float64) error {
	// Convert the float64 to a string for regex validation
	numStr := strconv.FormatFloat(num, 'f', -1, 64)
	// JSON number regex based on ECMA-404 specification
	numberRegex := `^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$`

	// compile and match the regex
	re := regexp.MustCompile(numberRegex)
	if !re.MatchString(numStr) {
		return errors.NewValidationError("invalid number format")
	}
	return nil
}
