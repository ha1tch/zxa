package zxa_assembler

import (
	"fmt"
	"path/filepath"
)

// ErrorCategory represents the type of assembler error
type ErrorCategory int

const (
	ErrNone ErrorCategory = iota
	ErrSyntax            // Syntax errors in assembly code
	ErrSymbol            // Symbol-related errors (undefined, duplicate)
	ErrValue            // Value errors (out of range, invalid)
	ErrFile             // File handling errors
	ErrDirective        // Directive errors
	ErrRange            // Range errors (jumps too far, etc)
	ErrInternal         // Internal assembler errors
)

// String returns the string representation of an error category
func (e ErrorCategory) String() string {
	switch e {
	case ErrSyntax:
		return "syntax error"
	case ErrSymbol:
		return "symbol error"
	case ErrValue:
		return "value error"
	case ErrFile:
		return "file error"
	case ErrDirective:
		return "directive error"
	case ErrRange:
		return "range error"
	case ErrInternal:
		return "internal error"
	default:
		return "unknown error"
	}
}

// AssemblerError represents a detailed error with category and location
type AssemblerError struct {
	Category ErrorCategory
	Message  string
	File     string
	Line     int
	Column   int // Column where error was detected
}

// Error implements the error interface for AssemblerError
func (e AssemblerError) Error() string {
	// Get base filename for cleaner output
	filename := filepath.Base(e.File)
	
	if e.Column > 0 {
		return fmt.Sprintf("%s:%d:%d: %s: %s", 
			filename, e.Line, e.Column, e.Category, e.Message)
	}
	return fmt.Sprintf("%s:%d: %s: %s", 
		filename, e.Line, e.Category, e.Message)
}

// ErrorList represents a collection of assembler errors
type ErrorList struct {
	errors []AssemblerError
}

// Add appends a new error to the list
func (l *ErrorList) Add(err AssemblerError) {
	l.errors = append(l.errors, err)
}

// HasErrors returns true if the list contains any errors
func (l *ErrorList) HasErrors() bool {
	return len(l.errors) > 0
}

// Errors returns the slice of errors
func (l *ErrorList) Errors() []AssemblerError {
	return l.errors
}

// Error implements the error interface for ErrorList
func (l *ErrorList) Error() string {
	if len(l.errors) == 0 {
		return "no errors"
	}
	
	if len(l.errors) == 1 {
		return l.errors[0].Error()
	}
	
	return fmt.Sprintf("%s (and %d more errors)", 
		l.errors[0].Error(), len(l.errors)-1)
}

// Error creation helper functions
func syntaxError(file string, line, col int, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrSyntax,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
		Line:     line,
		Column:   col,
	}
}

func symbolError(file string, line int, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrSymbol,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
		Line:     line,
	}
}

func valueError(file string, line int, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrValue,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
		Line:     line,
	}
}

func fileError(file string, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrFile,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
	}
}

func directiveError(file string, line int, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrDirective,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
		Line:     line,
	}
}

func rangeError(file string, line int, msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrRange,
		Message:  fmt.Sprintf(msg, args...),
		File:     file,
		Line:     line,
	}
}

func internalError(msg string, args ...interface{}) AssemblerError {
	return AssemblerError{
		Category: ErrInternal,
		Message:  fmt.Sprintf(msg, args...),
	}
}