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

func TestLexer_StringsWithEscapes(t *testing.T) {
	input := `"hello\nworld" "escaped \"quote\""`
	expectedTokens := []Token{
		{Type: TokenString, Literal: "hello\nworld"},
		{Type: TokenString, Literal: `escaped "quote"`},
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

func TestLexer_UnicodeStrings(t *testing.T) {
	input := `"unicode \u0041" "emoji \u1F600"`
	expectedTokens := []Token{
		{Type: TokenString, Literal: "unicode A"},
		{Type: TokenString, Literal: "emoji 😀"}, // \u1F600 is 😀
		{Type: TokenEOF, Literal: ""},
	}

	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Note: Ensure your readUnicode method correctly parses \u1F600
	if !reflect.DeepEqual(tokens, expectedTokens) {
		t.Errorf("expected tokens %v, got %v", expectedTokens, tokens)
	}
}
