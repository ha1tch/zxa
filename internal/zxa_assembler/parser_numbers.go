// file: internal/zxa_assembler/parser_numbers.go

package zxa_assembler

import (
	"fmt"
	"strconv"
	"strings"
)

// Number formats supported by the assembler
const (
	fmtDecimal = iota
	fmtHexDollar    // $NNNN format
	fmtHexC         // 0xNNNN format
	fmtHexSuffix    // NNNNh format
	fmtBinaryPercent // %NNNN format
	fmtBinaryC      // 0bNNNN format
)

// parseNumber takes a string and determines its format and value
func parseNumber(s string) (int64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty number")
	}

	// Determine format based on prefix or suffix
	format := fmtDecimal
	numStr := s

	switch {
	case strings.HasPrefix(s, "$"):
		format = fmtHexDollar
		numStr = strings.TrimPrefix(s, "$")
	case strings.HasPrefix(s, "0x"):
		format = fmtHexC
		numStr = strings.TrimPrefix(s, "0x")
	case strings.HasSuffix(s, "h") || strings.HasSuffix(s, "H"):
		format = fmtHexSuffix
		numStr = strings.TrimSuffix(strings.TrimSuffix(s, "h"), "H")
	case strings.HasPrefix(s, "%"):
		format = fmtBinaryPercent
		numStr = strings.TrimPrefix(s, "%")
	case strings.HasPrefix(s, "0b"):
		format = fmtBinaryC
		numStr = strings.TrimPrefix(s, "0b")
	}

	// Validate no empty number after prefix/before suffix
	if numStr == "" {
		switch format {
		case fmtHexDollar:
			return 0, fmt.Errorf("empty number after $ prefix")
		case fmtHexC:
			return 0, fmt.Errorf("empty number after 0x prefix")
		case fmtHexSuffix:
			return 0, fmt.Errorf("empty number before h suffix")
		case fmtBinaryPercent:
			return 0, fmt.Errorf("empty number after %% prefix")
		case fmtBinaryC:
			return 0, fmt.Errorf("empty number after 0b prefix")
		}
	}

	// Parse according to format
	var base int
	switch format {
	case fmtDecimal:
		base = 10
	case fmtHexDollar, fmtHexC, fmtHexSuffix:
		base = 16
		// For hex suffix format, first character must be 0-9 if it starts with a letter
		if format == fmtHexSuffix && len(numStr) > 0 {
			firstChar := numStr[0]
			if (firstChar >= 'a' && firstChar <= 'f') || (firstChar >= 'A' && firstChar <= 'F') {
				return 0, fmt.Errorf("hex number with h suffix must start with 0-9: %s", s)
			}
		}
	case fmtBinaryPercent, fmtBinaryC:
		base = 2
	}

	// Parse the value
	val, err := strconv.ParseInt(numStr, base, 32)
	if err != nil {
		switch format {
		case fmtDecimal:
			return 0, fmt.Errorf("invalid decimal number: %s", s)
		case fmtHexDollar, fmtHexC, fmtHexSuffix:
			return 0, fmt.Errorf("invalid hexadecimal number: %s", s)
		case fmtBinaryPercent, fmtBinaryC:
			return 0, fmt.Errorf("invalid binary number: %s", s)
		}
	}

	return val, nil
}

// formatNumber converts a number to a string in the specified format
func formatNumber(val int64, format int) string {
	switch format {
	case fmtHexDollar:
		return fmt.Sprintf("$%X", val)
	case fmtHexC:
		return fmt.Sprintf("0x%X", val)
	case fmtHexSuffix:
		return fmt.Sprintf("%Xh", val)
	case fmtBinaryPercent:
		return fmt.Sprintf("%%%b", val)
	case fmtBinaryC:
		return fmt.Sprintf("0b%b", val)
	default:
		return fmt.Sprintf("%d", val)
	}
}


// validateNumberString validates a number string based on its format
func validateNumberString(s string) error {
	if s == "" {
		return fmt.Errorf("empty number string")
	}

	// First check for hex suffix
	if strings.HasSuffix(s, "h") || strings.HasSuffix(s, "H") {
		numStr := strings.TrimSuffix(strings.TrimSuffix(s, "h"), "H")
		// First digit must be 0-9 for hex suffix format
		if len(numStr) > 0 {
			firstChar := numStr[0]
			if (firstChar >= 'a' && firstChar <= 'f') || (firstChar >= 'A' && firstChar <= 'F') {
				return fmt.Errorf("hex number with h suffix must start with 0-9: %s", s)
			}
		}
		for _, c := range []byte(numStr) {
			if !isValidHexDigit(c) {
				return fmt.Errorf("invalid hex digit: %c", c)
			}
		}
		return nil
	}

	// Then check other formats
	switch {
	case strings.HasPrefix(s, "$"):
		for _, c := range []byte(strings.TrimPrefix(s, "$")) {
			if !isValidHexDigit(c) {
				return fmt.Errorf("invalid hex digit: %c", c)
			}
		}
	case strings.HasPrefix(s, "0x"):
		for _, c := range []byte(strings.TrimPrefix(s, "0x")) {
			if !isValidHexDigit(c) {
				return fmt.Errorf("invalid hex digit: %c", c)
			}
		}
	case strings.HasPrefix(s, "%"):
		for _, c := range []byte(strings.TrimPrefix(s, "%")) {
			if !isValidBinaryDigit(c) {
				return fmt.Errorf("invalid binary digit: %c", c)
			}
		}
	case strings.HasPrefix(s, "0b"):
		for _, c := range []byte(strings.TrimPrefix(s, "0b")) {
			if !isValidBinaryDigit(c) {
				return fmt.Errorf("invalid binary digit: %c", c)
			}
		}
	default:
		// Decimal
		for _, c := range []byte(s) {
			if c < '0' || c > '9' {
				return fmt.Errorf("invalid decimal digit: %c", c)
			}
		}
	}

	return nil
}

// isValidHexDigit checks if a byte is a valid hexadecimal digit
func isValidHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')
}

// isValidBinaryDigit checks if a byte is a valid binary digit
func isValidBinaryDigit(c byte) bool {
	return c == '0' || c == '1'
}




// readNumber reads a numeric token
func (p *Parser) readNumber() (Token, error) {
	start := p.pos
	startCol := p.column

	// Handle hex and binary prefixes
	var isHex, isBin bool
	if p.pos < len(p.input) {
		if p.input[p.pos] == '$' {
			isHex = true
			p.pos++
			p.column++
		} else if p.input[p.pos] == '%' {
			isBin = true
			p.pos++
			p.column++
		} else if p.pos+1 < len(p.input) {
			if p.input[p.pos:p.pos+2] == "0x" {
				isHex = true
				p.pos += 2
				p.column += 2
			} else if p.input[p.pos:p.pos+2] == "0b" {
				isBin = true
				p.pos += 2
				p.column += 2
			}
		}
	}

	numberStart := p.pos
	// Read the number part
	for p.pos < len(p.input) {
		c := rune(p.input[p.pos])
		if isSpace(c) || c == ',' || c == ')' || c == ';' || c == '\n' {
			break
		}

		// Handle hex digits
		if isHex {
			if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
				break
			}
		} else if isBin {
			if c != '0' && c != '1' {
				break
			}
		} else {
			if !isDigit(c) {
				break
			}
		}

		p.pos++
		p.column++
	}

	if p.pos <= numberStart {
		return Token{}, fmt.Errorf("empty number at line %d, column %d", p.line, p.column)
	}

	value := p.input[start:p.pos]
	if p.debug {
		fmt.Printf("DEBUG: readNumber: value='%s' isHex=%v isBin=%v\n", value, isHex, isBin)
	}

	return Token{TokenNumber, value, p.line, startCol}, nil
}


// evaluateExpression evaluates numeric expressions with various prefixes
func (p *Parser) evaluateExpression(expr string) (int, error) {
	expr = strings.TrimSpace(expr)

	if p.debug {
		fmt.Printf("DEBUG: evaluateExpression: expr='%s'\n", expr)
	}

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

	// Handle hex values with h suffix
	if strings.HasSuffix(expr, "h") || strings.HasSuffix(expr, "H") {
		hex := strings.TrimSuffix(strings.TrimSuffix(expr, "h"), "H")
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
	val, err := strconv.ParseInt(expr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}
