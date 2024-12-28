package lexer

// TokenType represents the type of lexical token.
type TokenType string

// TokenType constants define all possible token types in the lexer.
const (
	ILLEGAL TokenType = "ILLEGAL" // Represents an invalid or unrecognized token
	EOF     TokenType = "EOF"     // End of input

	LBRACE   TokenType = "{" // Left brace
	RBRACE   TokenType = "}" // Right brace
	LBRACKET TokenType = "[" // Left bracket
	RBRACKET TokenType = "]" // Right bracket
	COLON    TokenType = ":" // Colon
	COMMA    TokenType = "," // Comma

	STRING TokenType = "STRING" // String literal
	NUMBER TokenType = "NUMBER" // Numeric literal

	TRUE  TokenType = "TRUE"  // Boolean literal: true
	FALSE TokenType = "FALSE" // Boolean literal: false
	NULL  TokenType = "NULL"  // Null literal
)

// Token represents a single lexical token with its type, literal value, and position in the input.
type Token struct {
	Type    TokenType // The type of the token
	Literal string    // The literal value of the token
	Line    int       // Line number where the token appears
	Column  int       // Column number where the token appears
}
