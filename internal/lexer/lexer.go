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

func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token

	for {
		l.skipWhiteSpace()
		var tok Token

		switch l.ch {
		case '{':
			tok = Token{Type: TokenLeftBrace, Literal: "{"}
		case '}':
			tok = Token{Type: TokenRightBrace, Literal: "}"}
		case ':':
			tok = Token{Type: TokenColon, Literal: ":"}
		case ',':
			tok = Token{Type: TokenComma, Literal: ","}
		case '"':
			str, err := l.readString()
			if err != nil {
				return nil, err
			}
			tok = Token{Type: TokenColon}
		}
	}
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch =='\r'
}