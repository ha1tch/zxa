package zxa_assembler

import (
	"fmt"
	"strings"
)

// parseInstruction handles the parsing of Z80 instructions and their operands
func (p *Parser) parseInstruction(token Token) error {
	mnemonic := strings.ToUpper(token.Value)
	operands := []string{}

	// Read operands until end of line
	for {
		tok, err := p.nextToken()
		if err != nil {
			return err
		}

		if tok.Type == TokenNone || tok.Type == TokenDirective {
			break
		}

		// Handle different operand types
		switch tok.Type {
		case TokenRegister:
			operands = append(operands, tok.Value)

		case TokenNumber:
			operands = append(operands, tok.Value)

		case TokenIdentifier:
			operands = append(operands, tok.Value)

		case TokenLParen:
			// Handle indirect addressing
			indirectOp, err := p.parseIndirectOperand()
			if err != nil {
				return err
			}
			operands = append(operands, fmt.Sprintf("(%s)", indirectOp))

		case TokenComma:
			continue

		default:
			return fmt.Errorf("unexpected token in instruction at line %d: %s",
				tok.Line, tok.Value)
		}

		// Check for comma between operands
		next, err := p.nextToken()
		if err != nil {
			return err
		}
		if next.Type != TokenComma && next.Type != TokenNone {
			return fmt.Errorf("expected comma between operands at line %d",
				next.Line)
		}
	}

	// Build the complete instruction
	fullInst := buildInstructionString(mnemonic, operands)

	// Look up the instruction in our instruction set
	inst, exists := p.assembler.instructions[fullInst]
	if !exists {
		return fmt.Errorf("unknown instruction at line %d: %s",
			token.Line, fullInst)
	}

	// Generate binary for the instruction
	if err := p.generateInstructionCode(inst, operands); err != nil {
		return err
	}

	return nil
}

// parseIndirectOperand handles (HL), (IX+d), etc.
func (p *Parser) parseIndirectOperand() (string, error) {
	var result strings.Builder
	
	// Read tokens until closing parenthesis
	for {
		tok, err := p.nextToken()
		if err != nil {
			return "", err
		}

		switch tok.Type {
		case TokenRParen:
			return result.String(), nil

		case TokenRegister, TokenIdentifier:
			result.WriteString(tok.Value)

		case TokenPlus, TokenMinus:
			result.WriteString(tok.Value)

		case TokenNumber:
			result.WriteString(tok.Value)

		default:
			return "", fmt.Errorf("unexpected token in indirect addressing at line %d: %s",
				tok.Line, tok.Value)
		}
	}
}

// buildInstructionString creates the instruction lookup key
func buildInstructionString(mnemonic string, operands []string) string {
	if len(operands) == 0 {
		return mnemonic
	}
	return fmt.Sprintf("%s %s", mnemonic, strings.Join(operands, ","))
}

// generateInstructionCode outputs the binary for an instruction
func (p *Parser) generateInstructionCode(inst Instruction, operands []string) error {
	// Output prefix byte if any
	if inst.Prefix != 0 {
		p.assembler.emitByte(inst.Prefix)
	}

	// Output main opcode
	p.assembler.emitByte(inst.Opcode)

	// Handle operands based on addressing mode
	switch inst.Mode {
	case Immediate:
		if len(operands) < 1 {
			return fmt.Errorf("immediate instruction requires operand")
		}
		val, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		p.assembler.emitByte(byte(val))

	case ImmediateExt:
		if len(operands) < 1 {
			return fmt.Errorf("extended immediate instruction requires operand")
		}
		val, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		p.assembler.emitByte(byte(val))
		p.assembler.emitByte(byte(val >> 8))

	case Indexed:
		if len(operands) < 1 {
			return fmt.Errorf("indexed addressing requires displacement")
		}
		disp, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		p.assembler.emitByte(byte(disp))

	case Relative:
		if len(operands) < 1 {
			return fmt.Errorf("relative instruction requires target")
		}
		// Handle relative jump calculation
		target, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		offset := target - (p.assembler.currentAddr + 2)
		if offset < -128 || offset > 127 {
			return fmt.Errorf("relative jump out of range")
		}
		p.assembler.emitByte(byte(offset))
	}

	return nil
}