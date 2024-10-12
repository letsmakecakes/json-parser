package lexer

import (
    "reflect"
    "testing"
)

func TestLexer_EmptyObject(t *testing.T) {
    input := "{}"
    expectedTokens := []Token{
        {Type: TokenLeftBrace, Literal: "{"},
        {Type: TokenRightBrace, Literal: "}"},
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
    // 😀 is represented by the surrogate pair \uD83D\uDE00
    input := `"unicode \u0041" "emoji \uD83D\uDE00"`
    expectedTokens := []Token{
        {Type: TokenString, Literal: "unicode A"},
        {Type: TokenString, Literal: "emoji 😀"}, // \uD83D\uDE00 represents 😀
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

func TestLexer_Numbers(t *testing.T) {
    input := `123 -456 78.90 -0.12`
    expectedTokens := []Token{
        {Type: TokenNumber, Literal: "123"},
        {Type: TokenNumber, Literal: "-456"},
        {Type: TokenNumber, Literal: "78.90"},
        {Type: TokenNumber, Literal: "-0.12"},
        {Type: TokenEOF, Literal: ""},
    }

    lexer := NewLexer(input)
    tokens, err := lexer.Tokenize()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !reflect.DeepEqual(tokens, expectedTokens) {
        t.Errorf("expected tokens %v, got %v", expectedTokens, tokens)
        for i, tok := range tokens {
            t.Logf("Token %d: Type=%s, Literal=%s", i, tok.Type, tok.Literal)
        }
    }
}

func TestLexer_Literals(t *testing.T) {
    input := `true false null`
    expectedTokens := []Token{
        {Type: TokenTrue, Literal: "true"},
        {Type: TokenFalse, Literal: "false"},
        {Type: TokenNull, Literal: "null"},
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

func TestLexer_ComplexStructure(t *testing.T) {
    input := `{
        "name": "John Doe",
        "age": 30,
        "isStudent": false,
        "scores": [85, 90, 92.5],
        "address": {
            "street": "123 Main St",
            "city": "Anytown"
        }
    }`
    expectedTokens := []Token{
        {Type: TokenLeftBrace, Literal: "{"},
        {Type: TokenString, Literal: "name"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenString, Literal: "John Doe"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenString, Literal: "age"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenNumber, Literal: "30"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenString, Literal: "isStudent"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenFalse, Literal: "false"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenString, Literal: "scores"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenLeftBracket, Literal: "["},
        {Type: TokenNumber, Literal: "85"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenNumber, Literal: "90"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenNumber, Literal: "92.5"},
        {Type: TokenRightBracket, Literal: "]"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenString, Literal: "address"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenLeftBrace, Literal: "{"},
        {Type: TokenString, Literal: "street"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenString, Literal: "123 Main St"},
        {Type: TokenComma, Literal: ","},
        {Type: TokenString, Literal: "city"},
        {Type: TokenColon, Literal: ":"},
        {Type: TokenString, Literal: "Anytown"},
        {Type: TokenRightBrace, Literal: "}"},
        {Type: TokenRightBrace, Literal: "}"},
        {Type: TokenEOF, Literal: ""},
    }

    lexer := NewLexer(input)
    tokens, err := lexer.Tokenize()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if !reflect.DeepEqual(tokens, expectedTokens) {
        t.Errorf("expected tokens %v, got %v", expectedTokens, tokens)
        for i, tok := range tokens {
            t.Logf("Token %d: Type=%s, Literal=%s", i, tok.Type, tok.Literal)
        }
    }
}
