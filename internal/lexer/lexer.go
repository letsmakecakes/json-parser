package lexer

type TokenType string

const (
	TokenLeftBrace  TokenType = "{"
	TokenRightBrace TokenType = "}"
	TokenColon      TokenType = ":"
	TokenComma      TokenType = ","
	TokenString     TokenType = "STRING"
	TokenNumber     TokenType = "NUMBER"
	TokenTrue       TokenType = "TRUE"
	TokenFalse      TokenType = "FALSE"
	TokenNull       TokenType = "NULL"
	TokenEOF        TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}