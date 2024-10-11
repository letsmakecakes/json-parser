package parser

import "github.com/letsmakecakes/jsonparser/internal/lexer"

type Parser struct {
	tokens  []lexer.Token
	current int
}
