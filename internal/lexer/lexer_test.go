package lexer

import "testing"

func TestNextToken_EmptyInput(t *testing.T) {
	// Test empty input string
	input := ""

	// Expected token is EOF
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{EOF, ""}, // End of input
	}

	// Initialize the lexer with the empty input
	lexer := New(input)

	// Check for EOF token
	tok := lexer.NextToken()

	// Log the token for debugging purposes
	t.Logf("Token: Type=%q, Literal=%q, Line=%d, Column=%d", tok.Type, tok.Literal, tok.Line, tok.Column)

	// Validate the token type
	if tok.Type != tests[0].expectedType {
		t.Fatalf("expected=%q, got=%q (literal=%q)", tests[0].expectedType, tok.Type, tok.Literal)
	}

	// Validate the token literal
	if tok.Literal != tests[0].expectedLiteral {
		t.Fatalf("expected=%q, got=%q (type=%q)", tests[0].expectedLiteral, tok.Literal, tok.Type)
	}
}

func TestNextToken_ValidCharacters(t *testing.T) {
	// Define the input string with various edge cases and scenarios
	input := `{"key": "value", "numbers": [1, 2, 3], "emptyObject": {}, "nestedArray": [[1, 2], [3, 4]]}`

	// Define the expected sequence of tokens for these scenarios
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LBRACE, "{"},           // Start of object
		{STRING, "key"},         // String key
		{COLON, ":"},            // Colon separator
		{STRING, "value"},       // String value
		{COMMA, ","},            // Comma separator
		{STRING, "numbers"},     // String key
		{COLON, ":"},            // Colon separator
		{LBRACKET, "["},         // Start of array
		{NUMBER, "1"},           // Number
		{COMMA, ","},            // Comma separator
		{NUMBER, "2"},           // Number
		{COMMA, ","},            // Comma separator
		{NUMBER, "3"},           // Number
		{RBRACKET, "]"},         // End of array
		{COMMA, ","},            // Comma separator
		{STRING, "emptyObject"}, // String key
		{COLON, ":"},            // Colon separator
		{LBRACE, "{"},           // Start of empty object
		{RBRACE, "}"},           // End of empty object
		{COMMA, ","},            // Comma separator
		{STRING, "nestedArray"}, // String key
		{COLON, ":"},            // Colon separator
		{LBRACKET, "["},         // Start of nested array
		{LBRACKET, "["},         // Start of inner array
		{NUMBER, "1"},           // Number in inner array
		{COMMA, ","},            // Comma separator
		{NUMBER, "2"},           // Number in inner array
		{RBRACKET, "]"},         // End of inner array
		{COMMA, ","},            // Comma separator
		{LBRACKET, "["},         // Start of another inner array
		{NUMBER, "3"},           // Number in inner array
		{COMMA, ","},            // Comma separator
		{NUMBER, "4"},           // Number in inner array
		{RBRACKET, "]"},         // End of inner array
		{RBRACKET, "]"},         // End of nested array
		{EOF, ""},               // End of input
	}

	// Initialize the lexer with the input
	lexer := New(input)

	// Iterate over the expected tokens and compare with lexer output
	for i, tt := range tests {
		tok := lexer.NextToken()

		// Log the token for debugging purposes
		t.Logf("Token[%d]: Type=%q, Literal=%q, Line=%d, Column=%d", i, tok.Type, tok.Literal, tok.Line, tok.Column)

		// Validate the token type
		if tok.Type != tt.expectedType {
			t.Fatalf(
				"tests[%d] - wrong token type. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal,
			)
		}

		// Validate the token literal
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf(
				"tests[%d] - wrong literal. expected=%q, got=%q (type=%q)",
				i, tt.expectedLiteral, tok.Literal, tok.Type,
			)
		}
	}
}
