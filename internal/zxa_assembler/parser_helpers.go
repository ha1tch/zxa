// file: internal/zxa_assembler/parser_helpers.go

package zxa_assembler

import (
	"strings"
	"unicode"
)

// isSpace returns true if the character is whitespace
func isSpace(c rune) bool {
	return c == ' ' || c == '\t'
}

// isAlpha returns true if the character is alphabetic
func isAlpha(c rune) bool {
	return unicode.IsLetter(c) || c == '_'
}

// isDigit returns true if the character is a digit
func isDigit(c rune) bool {
	return unicode.IsDigit(c)
}

// isAlphaNum returns true if the character is alphanumeric
func isAlphaNum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

// isRegister checks if a string represents a Z80 register name
func isRegister(s string) bool {
	registers := map[string]bool{
		"A": true, "B": true, "C": true, "D": true,
		"E": true, "H": true, "L": true, "I": true,
		"R": true, "BC": true, "DE": true, "HL": true,
		"SP": true, "IX": true, "IY": true, "AF": true,
	}
	return registers[strings.ToUpper(s)]
}

// isDirective checks if a string represents an assembler directive
func isDirective(s string) bool {
	directives := map[string]bool{
		"ORG": true, "EQU": true, "DEFB": true,
		"DEFW": true, "DEFS": true, "INCLUDE": true,
		"INCBIN": true,
	}
	return directives[strings.ToUpper(s)]
}

