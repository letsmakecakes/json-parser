package parser

import (
	"github.com/letsmakecakes/jsonparser/internal/ast"
	"github.com/letsmakecakes/jsonparser/internal/lexer"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func (p *Parser) parseObject() (*ast.Object, error) {
	obj := &ast.Object{}

	if !p.expectCurrent(lexer.TokenLeftBrace) {
		return nil, lexer.NewUnexpectedTokenError(p.peek(), lexer.TokenLeftBrace)
	}

	p.nextToken()

	for !p.peekTypeIs(lexer.TokenRightBrace) && !p.peekTypeIs(lexer.TokenEOF) {
		keyToken := p.peek()
		if keyToken.Type != lexer.TokenString {
			return nil, lexer.NewUnexpectedTokenError(keyToken, lexer.TokenString)
		}
		key := keyToken.Literal
		p.nextToken()

		if !p.expectedCurrent(lexer.TokenColon) {
			return nil, lexer.NewUnexpectedTokenError(p.peek(), lexer.TokenColon)
		}
		p.nextToken()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj.Pairs[key] = value

		if p.peekTypeIs(lexer.TokenComma) {
			p.nextToken()
		} else {
			break
		}
	}

	if !p.expectCurrent(lexer.TokenRightBrace) {
		return nil, lexer.NewUnexpectedTokenError(p.peek(), lexer.TokenRightBrace)
	}

	return obj, nil
}

func (p *Parser) expectCurrent(tokenType lexer.TokenType) bool {
	return p.tokens[p.current].Type == tokenType
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) nextToken() {
	p.current = p.current + 1
}

func (p *Parser) peekTypeIs(tokenType lexer.TokenType) bool {
	return p.tokens[p.current].Type == tokenType
}

func (p *Parser) parseValue() (ast.Value, error) {
	tok := p.peek()
	switch tok.Type {
	case lexer.TokenString:
		p.nextToken()
		return &ast.String{Value: tok.Literal}, nil
	case lexer.TokenNumber:
		p.nextToken()
		return &ast.Number{Value: tok.Literal}, nil
	case lexer.TokenNull:
		p.nextToken()
		return &ast.Null{}, nil
	case lexer.TokenLeftBrace:
		return p.parseObject()
	case lexer.TokenLeftBracket:
		return p.parseArray()
	default:
		return nil, lexer.NewUnexpectedTokenError(tok, "a valid value")
	}
}

func (p *Parser) parseArray() (*ast.Array, error) {
	array := &ast.Array{}

	p.nextToken() // skip the opening bracket

	// Handle empty array case
	if p.peekTypeIs(lexer.TokenRightBracket) {
		p.nextToken() // consume the closing bracket
		return array, nil
	}

	for !p.peekTypeIs(lexer.TokenRightBracket) && !p.peekTypeIs(lexer.TokenEOF) {
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		array.Elements = append(array.Elements, value)

		if p.peekTypeIs(lexer.TokenComma) {
			p.nextToken() // skip the comma
		} else {
			break
		}
	}

	if !p.expectCurrent(lexer.TokenRightBracket) {
		return nil, lexer.NewUnexpectedTokenError(p.peek(), lexer.TokenRightBracket)
	}

	return nil, nil
}
