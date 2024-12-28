package lexer

import (
	"testing"
)

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
	input := `{"name":"Alice","age":25,"isStudent":false,"courses":["Math","Science","History"],"address":{"street":"123 Main St","city":"Metropolis","zipcode":"12345"},"graduationYear":null,"isGraduated":true,"height":5.7,"siblings":[{"name":"Bob","age":22},{"name":"Charlie","age":28}]}`

	// Define the expected sequence of tokens for these scenarios
	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LBRACE, "{"},              // Start of object
		{STRING, "name"},           // String key
		{COLON, ":"},               // Colon separator
		{STRING, "Alice"},          // String value
		{COMMA, ","},               // Comma separator
		{STRING, "age"},            // String key
		{COLON, ":"},               // Colon separator
		{NUMBER, "25"},             // Number
		{COMMA, ","},               // Comma separator
		{STRING, "isStudent"},      // String key
		{COLON, ":"},               // Colon separator
		{FALSE, "false"},           // Boolean value (false)
		{COMMA, ","},               // Comma separator
		{STRING, "courses"},        // String key
		{COLON, ":"},               // Colon separator
		{LBRACKET, "["},            // Start of array
		{STRING, "Math"},           // String in array
		{COMMA, ","},               // Comma separator
		{STRING, "Science"},        // String in array
		{COMMA, ","},               // Comma separator
		{STRING, "History"},        // String in array
		{RBRACKET, "]"},            // End of array
		{COMMA, ","},               // Comma separator
		{STRING, "address"},        // String key
		{COLON, ":"},               // Colon separator
		{LBRACE, "{"},              // Start of nested object
		{STRING, "street"},         // String key
		{COLON, ":"},               // Colon separator
		{STRING, "123 Main St"},    // String value
		{COMMA, ","},               // Comma separator
		{STRING, "city"},           // String key
		{COLON, ":"},               // Colon separator
		{STRING, "Metropolis"},     // String value
		{COMMA, ","},               // Comma separator
		{STRING, "zipcode"},        // String key
		{COLON, ":"},               // Colon separator
		{STRING, "12345"},          // String value
		{RBRACE, "}"},              // End of nested object
		{COMMA, ","},               // Comma separator
		{STRING, "graduationYear"}, // String key
		{COLON, ":"},               // Colon separator
		{NULL, "null"},             // Null value
		{COMMA, ","},               // Comma separator
		{STRING, "isGraduated"},    // String key
		{COLON, ":"},               // Colon separator
		{TRUE, "true"},             // Boolean value (true)
		{COMMA, ","},               // Comma separator
		{STRING, "height"},         // String key
		{COLON, ":"},               // Colon separator
		{NUMBER, "5.7"},            // Number value
		{COMMA, ","},               // Comma separator
		{STRING, "siblings"},       // String key
		{COLON, ":"},               // Colon separator
		{LBRACKET, "["},            // Start of nested array
		{LBRACE, "{"},              // Start of inner object
		{STRING, "name"},           // String key
		{COLON, ":"},               // Colon separator
		{STRING, "Bob"},            // String value
		{COMMA, ","},               // Comma separator
		{STRING, "age"},            // String key
		{COLON, ":"},               // Colon separator
		{NUMBER, "22"},             // Number value
		{RBRACE, "}"},              // End of inner object
		{COMMA, ","},               // Comma separator
		{LBRACE, "{"},              // Start of another inner object
		{STRING, "name"},           // String key
		{COLON, ":"},               // Colon separator
		{STRING, "Charlie"},        // String value
		{COMMA, ","},               // Comma separator
		{STRING, "age"},            // String key
		{COLON, ":"},               // Colon separator
		{NUMBER, "28"},             // Number value
		{RBRACE, "}"},              // End of inner object
		{RBRACKET, "]"},            // End of nested array
		{RBRACE, "}"},              // End of main object
		{EOF, ""},                  // End of input
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
