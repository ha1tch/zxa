// file: internal/zxa_assembler/parser_methods.go

package zxa_assembler

import (
	"fmt"
	"strconv"
	"strings"
)

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
	if p.debug {
		fmt.Printf("DEBUG: readIdentifier: value='%s'\n", value)
	}

	// Check if it's a register
	if isRegister(value) {
		if p.debug {
			fmt.Printf("DEBUG: Found register: %s\n", value)
		}
		return Token{TokenRegister, value, p.line, startCol}, nil
	}

	// Check if it's an instruction
	if p.isInstruction(value) {
		if p.debug {
			fmt.Printf("DEBUG: Found instruction: %s\n", value)
		}
		return Token{TokenInstruction, value, p.line, startCol}, nil
	}

	// Check if it's a directive
	if isDirective(value) {
		if p.debug {
			fmt.Printf("DEBUG: Found directive: %s\n", value)
		}
		return Token{TokenDirective, value, p.line, startCol}, nil
	}

	// Check for hex suffix
	if len(value) > 1 && (strings.HasSuffix(value, "h") || strings.HasSuffix(value, "H")) {
		// Validate hex number
		numPart := value[:len(value)-1]
		if _, err := strconv.ParseInt(numPart, 16, 32); err == nil {
			if p.debug {
				fmt.Printf("DEBUG: Found hex number (suffix): %s\n", value)
			}
			return Token{TokenNumber, value, p.line, startCol}, nil
		}
	}

	if p.debug {
		fmt.Printf("DEBUG: Found identifier: %s\n", value)
	}
	return Token{TokenIdentifier, value, p.line, startCol}, nil
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

	if p.debug {
		fmt.Printf("DEBUG: readString: value='%s'\n", value)
	}

	return Token{TokenString, value, p.line, startCol}, nil
}


func (p *Parser) isInstruction(s string) bool {
	_, ok := p.assembler.instructions[strings.ToUpper(s)]
	return ok
}