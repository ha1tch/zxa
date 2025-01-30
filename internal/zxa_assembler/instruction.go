package zxa_assembler

import (
	"fmt"
)

// AddressingMode represents different Z80 addressing modes
type AddressingMode int

const (
	Implied          AddressingMode = iota // No operand
	Register                               // Single register
	RegisterPair                           // Register pair
	Immediate                              // 8-bit immediate
	ImmediateExt                           // 16-bit immediate
	Extended                               // Extended addressing
	Indexed                                // Indexed addressing (IX+d, IY+d)
	Relative                               // Relative addressing (for jr, djnz)
	RegisterIndirect                       // Register indirect (HL)
	BitIndex                               // Bit operations
)

// CPUVariant represents different Z80 CPU variants
type CPUVariant uint

const (
	Z80Standard CPUVariant = iota
	Z80Next
)

// Instruction represents a Z80 instruction definition
type Instruction struct {
	Opcode    byte           // Base opcode
	Prefix    byte           // Instruction prefix (0 = none, 0xCB, 0xDD, 0xED, 0xFD)
	Mode      AddressingMode // Addressing mode
	Length    int            // Instruction length in bytes
	Cycles    int            // Base cycle count
	Condition bool           // True if instruction can be conditional
}

// InstructionMap holds all Z80 instructions indexed by mnemonic
type InstructionMap map[string]Instruction

// AssemblerOptions contains configuration for the assembler
type AssemblerOptions struct {
	Variant CPUVariant
	// Other options could be added here
}

// ForwardRef represents a forward reference to be resolved
type ForwardRef struct {
	Address int            // Where to patch
	Type    AddressingMode // How to patch (relative vs absolute)
	Length  int            // How many bytes to patch
	Target  string         // Target symbol name
}

// Assembler represents the assembler state
type Assembler struct {
	instructions InstructionMap
	output       []byte
	currentAddr  int
	currentLabel string
	symbols      map[string]Symbol
	forwardRefs  []ForwardRef
	originSet    bool
	includes     map[string]bool
	includePath  []string
	options      AssemblerOptions
}

// NewAssembler creates a new assembler instance
func NewAssembler(opts AssemblerOptions) *Assembler {
	a := &Assembler{
		output:      make([]byte, 0, 1024),
		symbols:     make(map[string]Symbol),
		forwardRefs: make([]ForwardRef, 0),
		includes:    make(map[string]bool),
		includePath: []string{"."},
		options:     opts,
	}

	// Initialize instruction set
	a.instructions = make(InstructionMap)
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
		Type:    mode,
		Length:  length,
		Target:  target,
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
