package lexer

import (
	"reflect"
	"testing"
)

func TestLexer_Tokenize_ValidEmptyObject(t *testing.T) {
	input := "{}"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []Token{
		{Type: TokenLeftBrace, Literal: "{"},
		{Type: TokenRightBrace, Literal: "}"},
		{Type: TokenEOF, Literal: ""},
	}

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}
}

func TestLexer_SimpleStrings(t *testing.T) {
	input := `"hello" "world"`
	expectedTokens := []Token{
		{Type: TokenString, Literal: "hello"},
		{Type: TokenString, Literal: "world"},
		{Type: TokenEOF, Literal: ""},
	}

	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("expected tokens %v, got %v", expectedTokens, tokens)
	}
}
