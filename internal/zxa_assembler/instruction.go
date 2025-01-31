package zxa_assembler

import (
	"fmt"
	"os"
	"path/filepath"
)

// AddressingMode represents different Z80 addressing modes
type AddressingMode int

const (
	Implied AddressingMode = iota  // No operand
	Register               // Single register
	RegisterPair          // Register pair
	Immediate             // 8-bit immediate
	ImmediateExt         // 16-bit immediate
	Extended             // Extended addressing
	Indexed              // Indexed addressing (IX+d, IY+d)
	Relative             // Relative addressing (for jr, djnz)
	RegisterIndirect     // Register indirect (HL)
	BitIndex             // Bit operations
)

// CPUVariant represents different Z80 CPU variants
type CPUVariant uint

const (
	Z80Standard CPUVariant = iota
	Z80Next
)

// Symbol represents a label or constant in the assembly
type Symbol struct {
	Name  string
	Value int
	Type  string // "label", "equ", or "forward"
}

// Instruction represents a Z80 instruction definition
type Instruction struct {
	Opcode      byte            // Base opcode
	Prefix      byte            // Instruction prefix (0 = none, 0xCB, 0xDD, 0xED, 0xFD)
	Mode        AddressingMode  // Addressing mode
	Length      int            // Instruction length in bytes
	Cycles      int            // Base cycle count
	Condition   bool           // True if instruction can be conditional
}

// InstructionMap holds all Z80 instructions indexed by mnemonic
type InstructionMap map[string]Instruction

// AssemblerOptions contains configuration for the assembler
type AssemblerOptions struct {
	Variant CPUVariant
}

// ForwardRef represents a forward reference to be resolved
type ForwardRef struct {
	Address    int             // Where to patch
	Type      AddressingMode  // How to patch (relative vs absolute)
	Length     int            // How many bytes to patch
	Target     string         // Target symbol name
}

// BinaryFile represents a binary file to be included
type BinaryFile struct {
	Filename string
	Offset   int    // Where to include in output
	Skip     int    // Bytes to skip from start of file
	Length   int    // Bytes to include (-1 for all)
}

// AssemblyResult represents the result of assembly
type AssemblyResult struct {
	Success    bool
	Binary     []byte
	HexDump    string
	JSONReport string
	Statistics AssemblyStats
}

// AssemblyStats contains assembly statistics
type AssemblyStats struct {
	BytesGenerated int
	LinesProcessed int
	SymbolsDefined int
}

// Assembler represents the assembler state
type Assembler struct {
	instructions  InstructionMap
	output        []byte
	currentAddr   int
	currentLabel  string
	symbols       map[string]Symbol
	forwardRefs   []ForwardRef
	originSet     bool
	includes      map[string]bool
	includePath   []string
	options       AssemblerOptions
	binaryFiles   []BinaryFile
	hexOutput     bool
	jsonOutput    bool
}

// NewAssembler creates a new assembler instance
func NewAssembler(opts AssemblerOptions) *Assembler {
	a := &Assembler{
		output:       make([]byte, 0, 1024),
		symbols:      make(map[string]Symbol),
		forwardRefs:  make([]ForwardRef, 0),
		includes:     make(map[string]bool),
		includePath:  []string{"."},
		options:      opts,
		instructions: make(InstructionMap),
	}

	// Initialize instruction set
	a.instructions.initBaseInstructions()
	a.instructions.initCBInstructions()
	a.instructions.initEDInstructions()
	a.instructions.initIndexInstructions()

	if opts.Variant == Z80Next {
		a.instructions.initZ80NInstructions()
	}

	return a
}

// emitByte adds a byte to the output
func (a *Assembler) emitByte(b byte) {
	a.output = append(a.output, b)
	a.currentAddr++
}

// addSymbol adds a symbol to the symbol table
func (a *Assembler) addSymbol(name string, value int) error {
	if _, exists := a.symbols[name]; exists {
		return fmt.Errorf("duplicate symbol: %s", name)
	}
	a.symbols[name] = Symbol{
		Name:  name,
		Value: value,
		Type:  "label",
	}
	return nil
}

// updateSymbol updates an existing symbol's value
func (a *Assembler) updateSymbol(name string, value int) error {
	if _, exists := a.symbols[name]; !exists {
		return fmt.Errorf("undefined symbol: %s", name)
	}
	sym := a.symbols[name]
	sym.Value = value
	a.symbols[name] = sym
	return nil
}

// addForwardRef adds a forward reference to be resolved later
func (a *Assembler) addForwardRef(target string, addr int, mode AddressingMode, length int) {
	a.forwardRefs = append(a.forwardRefs, ForwardRef{
		Address: addr,
		Type:   mode,
		Length: length,
		Target: target,
	})
}

// resolveForwardRefs resolves all forward references
func (a *Assembler) resolveForwardRefs() error {
	for _, ref := range a.forwardRefs {
		sym, exists := a.symbols[ref.Target]
		if !exists {
			return fmt.Errorf("undefined symbol: %s", ref.Target)
		}

		switch ref.Type {
		case Relative:
			offset := sym.Value - (ref.Address + 2)
			if offset < -128 || offset > 127 {
				return fmt.Errorf("relative jump out of range to %s", ref.Target)
			}
			a.output[ref.Address] = byte(offset)

		case Extended, ImmediateExt:
			value := uint16(sym.Value)
			a.output[ref.Address] = byte(value)
			a.output[ref.Address+1] = byte(value >> 8)

		default:
			a.output[ref.Address] = byte(sym.Value)
		}
	}
	return nil
}

// Lookup finds an instruction definition by mnemonic
func (a *Assembler) Lookup(mnemonic string) (Instruction, bool) {
	inst, ok := a.instructions[mnemonic]
	return inst, ok
}

// GetCurrentAddress returns the current assembly address
func (a *Assembler) GetCurrentAddress() int {
	return a.currentAddr
}

// GetOutput returns the assembled binary
func (a *Assembler) GetOutput() []byte {
	return a.output
}

// setOrigin sets the assembly origin point
func (a *Assembler) setOrigin(addr int) {
	a.currentAddr = addr
	a.originSet = true
}

// recordBinaryFile records a binary file for inclusion
func (a *Assembler) recordBinaryFile(filename string, skip, length int) {
	binFile := BinaryFile{
		Filename: filename,
		Offset:   a.currentAddr,
		Skip:     skip,
		Length:   length,
	}

	// Store the binary file information for later processing
	a.binaryFiles = append(a.binaryFiles, binFile)

	// Update current address based on file size
	// This will be validated during actual file processing
	if length == -1 {
		info, err := os.Stat(filename)
		if err == nil {
			fileSize := int(info.Size()) - skip
			a.currentAddr += fileSize
		}
	} else {
		a.currentAddr += length
	}
}

// processIncludeFile processes an included source file
func (a *Assembler) processIncludeFile(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	// Check for circular includes
	if a.includes[absPath] {
		return fmt.Errorf("circular include detected: %s", filename)
	}

	// Mark this file as included
	a.includes[absPath] = true
	defer delete(a.includes, absPath)

	// Read and process the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read include file %s: %v", filename, err)
	}

	// Create a new parser for this file
	parser := NewParser(string(content))
	parser.assembler = a

	// Process each line
	for !parser.isEOF() {
		if err := parser.parseLine(); err != nil {
			return fmt.Errorf("error in included file %s: %v", filename, err)
		}
	}

	return nil
}

// SetHexOutput configures hex dump output
func (a *Assembler) SetHexOutput(enabled bool) {
	a.hexOutput = enabled
}

// SetJSONOutput configures JSON report output
func (a *Assembler) SetJSONOutput(enabled bool) {
	a.jsonOutput = enabled
}

// AddIncludePath adds a directory to the include search path
func (a *Assembler) AddIncludePath(path string) {
	a.includePath = append(a.includePath, path)
}

// Assemble processes the input file and generates output
func (a *Assembler) Assemble(filename string) (AssemblyResult, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return AssemblyResult{}, fmt.Errorf("failed to read input file: %v", err)
	}

	parser := NewParser(string(content))
	parser.assembler = a

	// Process each line
	for !parser.isEOF() {
		if err := parser.parseLine(); err != nil {
			return AssemblyResult{}, err
		}
	}

	// Resolve forward references
	if err := a.resolveForwardRefs(); err != nil {
		return AssemblyResult{}, err
	}

	return AssemblyResult{
		Success: true,
		Binary:  a.output,
		Statistics: AssemblyStats{
			BytesGenerated: len(a.output),
			// Add other stats as needed
		},
	}, nil
}