package parser

import (
	"fmt"
	"github.com/letsmakecakes/jsonparser/internal/lexer"
	"strconv"
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

// Parse parses the input starting from the root and returns the root Node.
func (p *Parser) Parse() (Node, error) {
	if p.curToken.Type != lexer.LBRACE {
		return nil, fmt.Errorf("expected '{', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}

	return p.parseObject()
}

// parseObject parses an object and returns an ObjectValue node.
func (p *Parser) parseObject() (*ObjectValue, error) {
	object := &ObjectValue{Pairs: make(map[string]Value)}

	p.nextToken()

	// Handle an empty object
	if p.curToken.Type == lexer.RBRACE {
		p.nextToken()
		return object, nil
	}

	// Parse object contents.
	for p.curToken.Type != lexer.EOF {
		key, err := p.parseKey()
		if err != nil {
			return nil, err
		}

		if p.curToken.Type != lexer.COLON {
			return nil, fmt.Errorf("expected ':', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		object.Pairs[key] = value

		if p.curToken.Type == lexer.RBRACE {
			p.nextToken()
			return object, nil
		}

		if p.curToken.Type != lexer.COMMA {
			return nil, fmt.Errorf("expected ',' or '}', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
	}

	return nil, fmt.Errorf("unexpected end of input")
}

// parseKey parses a key in an object.
func (p *Parser) parseKey() (string, error) {
	if p.curToken.Type != lexer.STRING {
		return "", fmt.Errorf("expected string key, got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}
	key := p.curToken.Literal
	p.nextToken()
	return key, nil
}

// parseValue parses a value in an object or array and returns a Value node.
func (p *Parser) parseValue() (Value, error) {
	switch p.curToken.Type {
	case lexer.STRING:
		value := &StringValue{Value: p.curToken.Literal}
		p.nextToken()
		return value, nil
	case lexer.NUMBER:
		// Convert the string literal to a float64
		numValue, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse number: %v", err)
		}
		value := &NumberValue{Value: numValue}
		p.nextToken()
		return value, nil
	case lexer.TRUE:
		value := &BooleanValue{Value: true}
		p.nextToken()
		return value, nil
	case lexer.FALSE:
		value := &BooleanValue{Value: false}
		p.nextToken()
		return value, nil
	case lexer.NULL:
		value := &NullValue{}
		p.nextToken()
		return value, nil
	case lexer.LBRACE:
		return p.parseObject()
	case lexer.LBRACKET:
		return p.parseArray()
	default:
		return nil, fmt.Errorf("unexpected token %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
	}
}

// parseArray parses an array and returns an ArrayValue node.
func (p *Parser) parseArray() (*ArrayValue, error) {
	array := &ArrayValue{Elements: []Value{}}

	p.nextToken()

	// Handle an empty array.
	if p.curToken.Type == lexer.RBRACKET {
		p.nextToken()
		return array, nil
	}

	for {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		array.Elements = append(array.Elements, value)

		if p.curToken.Type == lexer.RBRACKET {
			p.nextToken()
			return array, nil
		}

		if p.curToken.Type != lexer.COMMA {
			return nil, fmt.Errorf("expected ',' or ']', got %s at line %d, column %d", p.curToken.Type, p.curToken.Line, p.curToken.Column)
		}
		p.nextToken()
	}
}
