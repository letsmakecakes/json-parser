package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	COLON    TokenType = ":"
	COMMA    TokenType = ","

	STRING TokenType = "STRING"
	NUMBER TokenType = "NUMBER"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	NULL   TokenType = "NULL"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
