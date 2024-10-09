package lexer

import "testing"

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
