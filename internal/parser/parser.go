package parser

import (
	"fmt"
	"github.com/letsmakecakes/jsonparser/internal/lexer"
)

// Parser is responsible for parsing tokens into a structured format.
type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

// New creates a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances the parser to the next token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Parse parses the input starting from the root.
func (p *Parser) Parse() error {
	if p.curToken.Type != lexer.LBRACE {
		return fmt.Errorf("expected '{', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}

	p.nextToken()

	// Handle an empty object.
	if p.curToken.Type == lexer.RBRACE {
		p.nextToken()
		return nil
	}

	// Parse object contents.
	for p.curToken.Type != lexer.EOF {
		if err := p.parseKeyValue(); err != nil {
			return err
		}

		if p.curToken.Type == lexer.RBRACE {
			p.nextToken()
			return nil
		}

		if p.curToken.Type != lexer.COMMA {
			return fmt.Errorf("expected , or }, got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
	}

	return fmt.Errorf("unexpected end of input")
}

// parseKeyValue parses a key-value pair in an object.
func (p *Parser) parseKeyValue() error {
	// Parse key.
	if p.curToken.Type != lexer.STRING {
		return fmt.Errorf("expected string key, got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}
	p.nextToken()

	// Parse colon.
	if p.curToken.Type != lexer.COLON {
		return fmt.Errorf("expected ':', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}
	p.nextToken()

	// Parse value.
	return p.parseValue()
}

// parseValue parses a value in an object or array.
func (p *Parser) parseValue() error {
	switch p.curToken.Type {
	case lexer.STRING, lexer.NUMBER, lexer.TRUE, lexer.FALSE, lexer.NULL:
		p.nextToken()
		return nil
	case lexer.LBRACE:
		return p.Parse()
	case lexer.LBRACKET:
		return p.parseArray()
	default:
		return fmt.Errorf("unexpected token %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}
}

// parseArray parses an array.
func (p *Parser) parseArray() error {
	p.nextToken()

	// Handle an empty array.
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken()
		return nil
	}

	for {
		if err := p.parseValue(); err != nil {
			return err
		}

		if p.curToken.Type == lexer.RBRACKET {
			p.nextToken()
			return nil
		}

		if p.curToken.Type != lexer.COMMA {
			return fmt.Errorf("expected ',' or ']', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
	}
}
