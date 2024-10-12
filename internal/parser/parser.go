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
		return nil, errors.NewUnexpectedCharacterError(p.peek(), lexer.TokenLeftBrace)
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
	case lexer.TokenLeftBrace:
		return p.parseArray()
	default:
		return nil, lexer.NewUnexpectedTokenError(tok, "a valid value")
	}
}