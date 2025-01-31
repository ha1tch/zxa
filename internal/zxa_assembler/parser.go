package zxa_assembler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
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
}

// NewParser creates a new parser instance
func NewParser(input string) *Parser {
	return &Parser{
		input:   input,
		pos:     0,
		line:    1,
		column:  1,
		tokens:  make([]Token, 0),
		current: 0,
	}
}

// isSpace returns true if the character is whitespace
func isSpace(c rune) bool {
	return c == ' ' || c == '\t'
}

// isAlpha returns true if the character is alphabetic
func isAlpha(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}

// isAlphaNum returns true if the character is alphanumeric
func isAlphaNum(c rune) bool {
	return isAlpha(c) || unicode.IsDigit(c)
}

// skipWhitespace skips any whitespace characters
func (p *Parser) skipWhitespace() {
	for p.pos < len(p.input) && isSpace(rune(p.input[p.pos])) {
		p.pos++
		p.column++
	}
}

// skipComment skips a comment to the end of line
func (p *Parser) skipComment() {
	for p.pos < len(p.input) && p.input[p.pos] != '\n' {
		p.pos++
	}
}

// isEOF checks if we've reached the end of input
func (p *Parser) isEOF() bool {
	return p.pos >= len(p.input)
}

// readIdentifier reads an identifier token
func (p *Parser) readIdentifier() (Token, error) {
	start := p.pos
	startCol := p.column

	for p.pos < len(p.input) && isAlphaNum(rune(p.input[p.pos])) {
		p.pos++
		p.column++
	}

	value := p.input[start:p.pos]

	// Check if it's a register
	if isRegister(value) {
		return Token{TokenRegister, value, p.line, startCol}, nil
	}

	// Check if it's an instruction
	if p.isInstruction(value) {
		return Token{TokenInstruction, value, p.line, startCol}, nil
	}

	// Check if it's a directive
	if isDirective(value) {
		return Token{TokenDirective, value, p.line, startCol}, nil
	}

	return Token{TokenIdentifier, value, p.line, startCol}, nil
}

// readNumber reads a numeric token
func (p *Parser) readNumber() (Token, error) {
	start := p.pos
	startCol := p.column
	isHex := false
	isBin := false

	// Check for hex/binary prefixes
	if p.pos+1 < len(p.input) {
		prefix := p.input[p.pos : p.pos+2]
		if prefix == "0x" || prefix == "0X" {
			isHex = true
			p.pos += 2
			p.column += 2
		} else if prefix == "0b" || prefix == "0B" {
			isBin = true
			p.pos += 2
			p.column += 2
		}
	} else if p.pos < len(p.input) {
		if p.input[p.pos] == '$' {
			isHex = true
			p.pos++
			p.column++
		} else if p.input[p.pos] == '%' {
			isBin = true
			p.pos++
			p.column++
		}
	}

	// Read digits
	for p.pos < len(p.input) &&
		(unicode.IsDigit(rune(p.input[p.pos])) ||
			(isHex && strings.ContainsRune("ABCDEFabcdef", rune(p.input[p.pos])))) {
		p.pos++
		p.column++
	}

	value := p.input[start:p.pos]

	// Validate number format based on prefix
	if isHex {
		if _, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(value, "0x"), "$"), 16, 32); err != nil {
			return Token{}, fmt.Errorf("invalid hex number at line %d: %s", p.line, value)
		}
	} else if isBin {
		if _, err := strconv.ParseInt(strings.TrimPrefix(strings.TrimPrefix(value, "0b"), "%"), 2, 32); err != nil {
			return Token{}, fmt.Errorf("invalid binary number at line %d: %s", p.line, value)
		}
	} else {
		if _, err := strconv.Atoi(value); err != nil {
			return Token{}, fmt.Errorf("invalid number at line %d: %s", p.line, value)
		}
	}

	return Token{TokenNumber, value, p.line, startCol}, nil
}

// readString reads a string token
func (p *Parser) readString() (Token, error) {
	startCol := p.column
	p.pos++ // Skip opening quote
	p.column++

	start := p.pos
	for p.pos < len(p.input) && p.input[p.pos] != '"' {
		if p.input[p.pos] == '\\' {
			p.pos += 2
			p.column += 2
		} else {
			p.pos++
			p.column++
		}
	}

	if p.pos >= len(p.input) {
		return Token{}, fmt.Errorf("unterminated string at line %d", p.line)
	}

	value := p.input[start:p.pos]
	p.pos++ // Skip closing quote
	p.column++

	return Token{TokenString, value, p.line, startCol}, nil
}

// nextToken returns the next token from the input
func (p *Parser) nextToken() (Token, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return Token{TokenNone, "", p.line, p.column}, nil
	}

	c := p.input[p.pos]

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

	case unicode.IsDigit(rune(c)) || c == '$' || c == '%':
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

// evaluateExpression evaluates numeric expressions with various prefixes
func (p *Parser) evaluateExpression(expr string) (int, error) {
	expr = strings.TrimSpace(expr)

	// Handle hex values (both $FFFF and 0xFFFF format)
	if strings.HasPrefix(expr, "$") {
		hex := strings.TrimPrefix(expr, "$")
		val, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	}
	if strings.HasPrefix(expr, "0x") {
		hex := strings.TrimPrefix(expr, "0x")
		val, err := strconv.ParseInt(hex, 16, 32)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	}

	// Handle binary values (both %1010 and 0b1010 format)
	if strings.HasPrefix(expr, "%") {
		bin := strings.TrimPrefix(expr, "%")
		val, err := strconv.ParseInt(bin, 2, 32)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	}
	if strings.HasPrefix(expr, "0b") {
		bin := strings.TrimPrefix(expr, "0b")
		val, err := strconv.ParseInt(bin, 2, 32)
		if err != nil {
			return 0, err
		}
		return int(val), nil
	}

	// Handle symbols
	if sym, exists := p.assembler.symbols[expr]; exists {
		return sym.Value, nil
	}

	// Default to decimal
	return strconv.Atoi(expr)
}

// Helper functions to check token types
func isRegister(s string) bool {
	registers := map[string]bool{
		"A": true, "B": true, "C": true, "D": true,
		"E": true, "H": true, "L": true, "I": true,
		"R": true, "BC": true, "DE": true, "HL": true,
		"SP": true, "IX": true, "IY": true, "AF": true,
	}
	return registers[strings.ToUpper(s)]
}

func (p *Parser) isInstruction(s string) bool {
	_, ok := p.assembler.instructions[strings.ToUpper(s)]
	return ok
}

func isDirective(s string) bool {
	directives := map[string]bool{
		"ORG": true, "EQU": true, "DEFB": true,
		"DEFW": true, "DEFS": true, "INCLUDE": true,
		"INCBIN": true,
	}
	return directives[strings.ToUpper(s)]
}

// parseLine parses a single line of assembly
func (p *Parser) parseLine() error {
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	// Empty line or end of file
	if token.Type == TokenNone {
		return nil
	}

	// Handle label definitions
	if token.Type == TokenIdentifier {
		nextToken, err := p.nextToken()
		if err != nil {
			return err
		}

		if nextToken.Type == TokenColon {
			// Process label
			if err := p.assembler.addSymbol(token.Value, p.assembler.currentAddr); err != nil {
				return err
			}
			// Get next token for instruction processing
			token, err = p.nextToken()
			if err != nil {
				return err
			}
		} else {
			// Not a label, put back the token
			p.tokens = append(p.tokens, nextToken)
		}
	}

	// Process instruction or directive
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
