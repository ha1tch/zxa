// file: internal/zxa_assembler/instruction_parser.go

package zxa_assembler

import (
	"fmt"
	"strconv"
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

		case TokenNumber, TokenIdentifier:
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
		if next.Type != TokenComma && next.Type != TokenNone && next.Type != TokenDirective {
			return fmt.Errorf("expected comma between operands at line %d",
				next.Line)
		}
	}

	// Build instruction format for lookup
	fullInst := buildInstructionString(mnemonic, operands)

	// Look up the instruction
	inst, exists := p.assembler.instructions[fullInst]
	if !exists {
		// Try generic format for immediate/extended instructions
		genericInst := buildGenericInstruction(mnemonic, operands)
		inst, exists = p.assembler.instructions[genericInst]
		if !exists {
			return fmt.Errorf("unknown instruction at line %d: %s",
				token.Line, fullInst)
		}
	}

	// Generate the instruction code
	if err := p.generateInstructionCode(inst, operands); err != nil {
		return err
	}

	return nil
}

// buildGenericInstruction creates a generic instruction format for lookup
func buildGenericInstruction(mnemonic string, operands []string) string {
	// Replace numeric values with format placeholders
	genericOps := make([]string, len(operands))
	for i, op := range operands {
		if isNumeric(op) {
			// Use 'n' for 8-bit immediates, 'nn' for 16-bit values
			if isEightBitValue(op) {
				genericOps[i] = "n"
			} else {
				genericOps[i] = "nn"
			}
		} else {
			genericOps[i] = op
		}
	}

	if len(genericOps) == 0 {
		return mnemonic
	}
	return fmt.Sprintf("%s %s", mnemonic, strings.Join(genericOps, ","))
}

// isNumeric checks if a string represents a numeric value
func isNumeric(s string) bool {
	// Remove common prefixes
	s = strings.TrimPrefix(s, "$")
	s = strings.TrimPrefix(s, "0x")
	s = strings.TrimPrefix(s, "%")
	s = strings.TrimPrefix(s, "0b")

	// Try parsing as different number formats
	_, err := strconv.ParseInt(s, 0, 32)
	return err == nil
}

// isEightBitValue checks if a numeric value fits in 8 bits
func isEightBitValue(s string) bool {
	val, err := strconv.ParseInt(s, 0, 32)
	if err != nil {
		return false
	}
	return val >= -128 && val <= 255
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
		if val < -128 || val > 255 {
			return fmt.Errorf("immediate value out of range: %d", val)
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
		if val < -32768 || val > 65535 {
			return fmt.Errorf("extended immediate value out of range: %d", val)
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
		if disp < -128 || disp > 127 {
			return fmt.Errorf("index displacement out of range: %d", disp)
		}
		p.assembler.emitByte(byte(disp))

	case Relative:
		if len(operands) < 1 {
			return fmt.Errorf("relative instruction requires target")
		}
		target, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		offset := target - (p.assembler.currentAddr + 2)
		if offset < -128 || offset > 127 {
			return fmt.Errorf("relative jump out of range")
		}
		p.assembler.emitByte(byte(offset))

	case Extended:
		if len(operands) < 1 {
			return fmt.Errorf("extended instruction requires address")
		}
		val, err := p.evaluateExpression(operands[0])
		if err != nil {
			return err
		}
		if val < 0 || val > 65535 {
			return fmt.Errorf("address out of range: %d", val)
		}
		p.assembler.emitByte(byte(val))
		p.assembler.emitByte(byte(val >> 8))
	}

	return nil
}
