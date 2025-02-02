// file: internal/zxa_assembler/parser.go

package zxa_assembler

import (
	"fmt"
	"strings"
)

// TokenType represents different types of assembly tokens
type TokenType int

const (
	TokenNone TokenType = iota
	TokenLabel
	TokenInstruction
	TokenDirective
	TokenRegister
	TokenNumber
	TokenString
	TokenComma
	TokenColon
	TokenLParen
	TokenRParen
	TokenPlus
	TokenMinus
	TokenIdentifier
)

// Token represents a lexical token from the assembly source
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

// Parser handles the parsing of assembly source code
type Parser struct {
	assembler *Assembler
	input     string
	pos       int
	line      int
	column    int
	tokens    []Token
	current   int
	debug     bool
}

// NewParser creates a new parser instance
func NewParser(input string, debug bool) *Parser {
	return &Parser{
		input:   input,
		pos:     0,
		line:    1,
		column:  1,
		tokens:  make([]Token, 0),
		current: 0,
		debug:   debug,
	}
}

// nextToken returns the next token from the input
func (p *Parser) nextToken() (Token, error) {
	// Return buffered token if any
	if len(p.tokens) > 0 {
		token := p.tokens[0]
		p.tokens = p.tokens[1:]
		if p.debug {
			fmt.Printf("DEBUG: nextToken: returning buffered token type=%v value='%s'\n", 
				token.Type, token.Value)
		}
		return token, nil
	}

	p.skipWhitespace()

	if p.pos >= len(p.input) {
		if p.debug {
			fmt.Printf("DEBUG: nextToken: EOF\n")
		}
		return Token{TokenNone, "", p.line, p.column}, nil
	}

	c := p.input[p.pos]

	if p.debug {
		fmt.Printf("DEBUG: nextToken: processing character '%c'\n", c)
	}

	switch {
	case c == ';':
		p.skipComment()
		return p.nextToken()

	case c == '\n':
		p.line++
		p.column = 1
		p.pos++
		return p.nextToken()

	case isAlpha(rune(c)):
		return p.readIdentifier()

	case isDigit(rune(c)) || c == '$' || c == '%':
		return p.readNumber()

	case c == '"':
		return p.readString()

	case c == ',':
		p.pos++
		p.column++
		return Token{TokenComma, ",", p.line, p.column - 1}, nil

	case c == ':':
		p.pos++
		p.column++
		return Token{TokenColon, ":", p.line, p.column - 1}, nil

	case c == '(':
		p.pos++
		p.column++
		return Token{TokenLParen, "(", p.line, p.column - 1}, nil

	case c == ')':
		p.pos++
		p.column++
		return Token{TokenRParen, ")", p.line, p.column - 1}, nil

	case c == '+':
		p.pos++
		p.column++
		return Token{TokenPlus, "+", p.line, p.column - 1}, nil

	case c == '-':
		p.pos++
		p.column++
		return Token{TokenMinus, "-", p.line, p.column - 1}, nil
	}

	return Token{}, fmt.Errorf("unexpected character '%c' at line %d, column %d",
		c, p.line, p.column)
}

// parseLine parses a single line of assembly
func (p *Parser) parseLine() error {
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if p.debug {
		fmt.Printf("DEBUG: parseLine: first token type=%v value='%s'\n", 
			token.Type, token.Value)
	}

	// Empty line or end of file
	if token.Type == TokenNone {
		return nil
	}

	// Handle label definitions
	if token.Type == TokenIdentifier {
		if p.debug {
			fmt.Printf("DEBUG: Processing identifier '%s'\n", token.Value)
		}
		
		nextToken, err := p.nextToken()
		if err != nil {
			return err
		}

		if p.debug {
			fmt.Printf("DEBUG: After identifier, next token type=%v value='%s'\n", 
				nextToken.Type, nextToken.Value)
		}

		switch nextToken.Type {
		case TokenColon:
			// Traditional label with colon
			if p.debug {
				fmt.Printf("DEBUG: Found label definition '%s:'\n",
					token.Value)
			}
			if err := p.assembler.addSymbol(token.Value, p.assembler.currentAddr); err != nil {
				return err
			}
			// Get next token for instruction processing
			token, err = p.nextToken()
			if err != nil {
				return err
			}

		case TokenDirective:
			// Handle case like "LABEL EQU value"
			if strings.ToUpper(nextToken.Value) == "EQU" {
				// Save the label for EQU processing
				p.assembler.currentLabel = token.Value
				return p.parseDirective(nextToken)
			}
			// Not EQU, treat as normal identifier
			p.tokens = append(p.tokens, nextToken)
			token = Token{TokenIdentifier, token.Value, token.Line, token.Column}

		default:
			// Not a label definition, put back both tokens
			p.tokens = append(p.tokens, nextToken)
			token = Token{TokenIdentifier, token.Value, token.Line, token.Column}
		}
	}

	// Process instruction or directive
	if p.debug {
		fmt.Printf("DEBUG: Processing instruction/directive token type=%v value='%s'\n",
			token.Type, token.Value)
	}

	switch token.Type {
	case TokenInstruction:
		return p.parseInstruction(token)
	case TokenDirective:
		return p.parseDirective(token)
	case TokenNone:
		return nil
	default:
		return fmt.Errorf("unexpected token at line %d: %s", token.Line, token.Value)
	}
}