package lexer

// Lexer tokenizes input string for parsing.
type Lexer struct {
	input        string // The input being tokenized.
	position     int    // Current position in the input (points to the current char)
	readPosition int    // Next position to read from input
	ch           byte   // Current character being examined
	line         int    // Current line in the input
	column       int    // Current column in the input
}

// New initializes and returns a new lexer instance.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar() // Initialize the first character
	return l
}

// readChar advances the lexer to the next character in the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // End of input
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// Update the line and column numbers
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// NextToken extracts the next token from the input.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{', '}', '[', ']', ':', ',':
		tok = l.makeSingleCharToken()
	case '"':
		tok = l.readString()
	case 0:
		tok = l.newToken(EOF, "")
	default:
		if isLetter(l.ch) {
			return l.readIdentifier()
		} else if isDigit(l.ch) {
			return l.newToken(NUMBER, l.readNumber())
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// skipWhitespace skips over spaces, tabs, and newlines in the input.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// makeSingleCharToken creates tokens for single-character symbols.
func (l *Lexer) makeSingleCharToken() Token {
	tokenType := singleCharTokenType(l.ch)
	return l.newToken(tokenType, string(l.ch))
}

// readString reads a string literal, handling any errors.
func (l *Lexer) readString() Token {
	startLine, startColumn := l.line, l.column
	start := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	if l.ch == 0 {
		return Token{
			Type:    ILLEGAL,
			Literal: "Unterminated string",
			Line:    startLine,
			Column:  startColumn,
		}
	}

	return l.newToken(STRING, l.input[start:l.position])
}

// readNumber reads a numeric literal.
func (l *Lexer) readNumber() string {
	start := l.position
	decimalSeen := false

	for isDigit(l.ch) || (l.ch == '.' && !decimalSeen) {
		if l.ch == '.' {
			decimalSeen = true
		}
		l.readChar()
	}
	return l.input[start:l.position]
}

// readIdentifier reads an identifier or keyword and returns the appropriate token.
func (l *Lexer) readIdentifier() Token {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	ident := l.input[start:l.position]
	tokenType := lookupKeyword(ident)
	return l.newToken(tokenType, ident)
}

// newToken creates a new Token with the current lexer state.
func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

// singleCharTokenType maps single-character symbols to their token types.
func singleCharTokenType(ch byte) TokenType {
	switch ch {
	case '{':
		return LBRACE
	case '}':
		return RBRACE
	case '[':
		return LBRACKET
	case ']':
		return RBRACKET
	case ':':
		return COLON
	case ',':
		return COMMA
	default:
		return ILLEGAL
	}
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// isLetter checks if a character is a letter (a-z or A-Z).
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// lookupKeyword determines if an identifier matches a keyword.
func lookupKeyword(ident string) TokenType {
	switch ident {
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "null":
		return NULL
	default:
		return ILLEGAL
	}
}
