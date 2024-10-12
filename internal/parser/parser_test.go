package parser

import (
	"testing"

	"github.com/letsmakecakes/jsonparser/internal/lexer"
)

func TestParse_EmptyObject(t *testing.T) {
	input := "{}"
	lex := lexer.NewLexer(input)
	tokens, err := lex.Tokenize()
	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	obj, err := Parse(tokens)
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if len(obj.Pairs) != 0 {
		t.Errorf("expected empty object, got %v", obj.Pairs)
	}
}
