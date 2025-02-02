// file: internal/zxa_assembler/parser_directives.go

package zxa_assembler

import (
	"fmt"
	"strings"
)

// parseDirective handles the parsing of assembler directives
func (p *Parser) parseDirective(token Token) error {
	directive := strings.ToUpper(token.Value)

	switch directive {
	case "ORG":
		return p.parseORG()
	case "EQU":
		return p.parseEQU()
	case "DEFB":
		return p.parseDEFB()
	case "DEFW":
		return p.parseDEFW()
	case "DEFS":
		return p.parseDEFS()
	case "INCLUDE":
		return p.parseINCLUDE()
	case "INCBIN":
		return p.parseINCBIN()
	default:
		return fmt.Errorf("unknown directive at line %d: %s",
			token.Line, directive)
	}
}

// parseORG handles the ORG directive
func (p *Parser) parseORG() error {
	// Get the address expression
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Type != TokenNumber && token.Type != TokenIdentifier {
		return fmt.Errorf("ORG requires address at line %d", token.Line)
	}

	// Evaluate the address
	addr, err := p.evaluateExpression(token.Value)
	if err != nil {
		return fmt.Errorf("invalid ORG address at line %d: %v", token.Line, err)
	}

	// Set the current address
	p.assembler.setOrigin(addr)

	return nil
}

// parseEQU handles the EQU directive
func (p *Parser) parseEQU() error {
	// Get the value expression
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Type != TokenNumber && token.Type != TokenIdentifier {
		return fmt.Errorf("EQU requires value at line %d", token.Line)
	}

	// Evaluate the value
	value, err := p.evaluateExpression(token.Value)
	if err != nil {
		return fmt.Errorf("invalid EQU value at line %d: %v", token.Line, err)
	}

	// Either use the current label or the last seen label
	if p.assembler.currentLabel == "" {
		return fmt.Errorf("EQU without label at line %d", token.Line)
	}

	// Add or update the symbol
	p.assembler.symbols[p.assembler.currentLabel] = Symbol{
		Name:  p.assembler.currentLabel,
		Value: value,
		Type:  "equ",
	}

	// Clear the current label
	p.assembler.currentLabel = ""

	return nil
}

// parseDEFB handles the DEFB directive
func (p *Parser) parseDEFB() error {
	for {
		token, err := p.nextToken()
		if err != nil {
			return err
		}

		if token.Type == TokenNone {
			break
		}

		switch token.Type {
		case TokenString:
			// Emit each character as a byte
			for _, c := range token.Value {
				p.assembler.emitByte(byte(c))
			}

		case TokenNumber, TokenIdentifier:
			// Evaluate and emit the byte value
			value, err := p.evaluateExpression(token.Value)
			if err != nil {
				return err
			}
			if value < -128 || value > 255 {
				return fmt.Errorf("DEFB value out of range at line %d: %d",
					token.Line, value)
			}
			p.assembler.emitByte(byte(value))

		case TokenComma:
			continue

		default:
			return fmt.Errorf("unexpected token in DEFB at line %d: %s",
				token.Line, token.Value)
		}
	}

	return nil
}

// parseDEFW handles the DEFW directive
func (p *Parser) parseDEFW() error {
	for {
		token, err := p.nextToken()
		if err != nil {
			return err
		}

		if token.Type == TokenNone {
			break
		}

		switch token.Type {
		case TokenNumber, TokenIdentifier:
			// Evaluate and emit the word value
			value, err := p.evaluateExpression(token.Value)
			if err != nil {
				return err
			}
			if value < -32768 || value > 65535 {
				return fmt.Errorf("DEFW value out of range at line %d: %d",
					token.Line, value)
			}
			p.assembler.emitByte(byte(value))
			p.assembler.emitByte(byte(value >> 8))

		case TokenComma:
			continue

		default:
			return fmt.Errorf("unexpected token in DEFW at line %d: %s",
				token.Line, token.Value)
		}
	}

	return nil
}

// parseDEFS handles the DEFS directive
func (p *Parser) parseDEFS() error {
	// Get size
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Type != TokenNumber && token.Type != TokenIdentifier {
		return fmt.Errorf("DEFS requires size at line %d", token.Line)
	}

	size, err := p.evaluateExpression(token.Value)
	if err != nil {
		return err
	}

	// Check for fill value
	fillValue := 0
	token, err = p.nextToken()
	if err != nil {
		return err
	}

	if token.Type == TokenComma {
		token, err = p.nextToken()
		if err != nil {
			return err
		}
		if token.Type != TokenNumber && token.Type != TokenIdentifier {
			return fmt.Errorf("invalid DEFS fill value at line %d", token.Line)
		}
		fillValue, err = p.evaluateExpression(token.Value)
		if err != nil {
			return err
		}
	}

	// Emit the fill bytes
	for i := 0; i < size; i++ {
		p.assembler.emitByte(byte(fillValue))
	}

	return nil
}

// parseINCLUDE handles the INCLUDE directive
func (p *Parser) parseINCLUDE() error {
	// Get filename
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Type != TokenString {
		return fmt.Errorf("INCLUDE requires filename at line %d", token.Line)
	}

	filename := token.Value

	// Process the included file
	if err := p.assembler.processIncludeFile(filename); err != nil {
		return fmt.Errorf("error processing include file %s: %v", filename, err)
	}

	return nil
}

// parseINCBIN handles the INCBIN directive
func (p *Parser) parseINCBIN() error {
	// Get filename
	token, err := p.nextToken()
	if err != nil {
		return err
	}

	if token.Type != TokenString {
		return fmt.Errorf("INCBIN requires filename at line %d", token.Line)
	}

	filename := token.Value

	// Check for optional skip and length parameters
	var skip, length int = 0, -1

	token, err = p.nextToken()
	if err != nil {
		return err
	}

	if token.Type == TokenComma {
		// Parse skip value
		token, err = p.nextToken()
		if err != nil {
			return err
		}
		skip, err = p.evaluateExpression(token.Value)
		if err != nil {
			return err
		}

		token, err = p.nextToken()
		if err != nil {
			return err
		}

		if token.Type == TokenComma {
			// Parse length value
			token, err = p.nextToken()
			if err != nil {
				return err
			}
			length, err = p.evaluateExpression(token.Value)
			if err != nil {
				return err
			}
		}
	}

	// Record the binary file for later processing
	p.assembler.recordBinaryFile(filename, skip, length)

	return nil
}
