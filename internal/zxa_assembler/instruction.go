package zxa_assembler

import (
	"fmt"
	"path/filepath"
)

type AddressingMode int

const (
	Implied AddressingMode = iota
	Register
	RegisterPair
	Immediate
	ImmediateExt
	Extended
	Indexed
	Relative
	RegisterIndirect
	BitIndex
)

type CPUVariant uint

const (
	Z80Standard CPUVariant = iota
	Z80Next
)

type Instruction struct {
	Opcode    byte
	Prefix    byte
	Mode      AddressingMode
	Length    int
	Cycles    int
	Condition bool
}

type InstructionMap map[string]Instruction

type AssemblerOptions struct {
	Variant CPUVariant
}

type Symbol struct {
	Name  string
	Value int
	Type  string
}

type ForwardRef struct {
	Address int
	Type    AddressingMode
	Length  int
	Target  string
}

type AssemblyResult struct {
	Success    bool
	Output     []byte
	Statistics struct {
		BytesGenerated int
		LinesProcessed int
		SymbolsDefined int
	}
}

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
	hexOutput    bool
	jsonOutput   bool
}

func NewInstructionSet(opts AssemblerOptions) (InstructionMap, error) {
	m := make(InstructionMap)
	m.initBaseInstructions()
	m.initCBInstructions()
	m.initEDInstructions()
	m.initIndexInstructions()

	if opts.Variant == Z80Next {
		m.initZ80NInstructions()
	}

	return m, nil
}

func NewAssembler(instructions InstructionMap) *Assembler {
	return &Assembler{
		instructions: instructions,
		output:      make([]byte, 0, 1024),
		symbols:     make(map[string]Symbol),
		forwardRefs: make([]ForwardRef, 0),
		includes:    make(map[string]bool),
		includePath: []string{"."},
	}
}

func (a *Assembler) SetHexOutput(enable bool) {
	a.hexOutput = enable
}

func (a *Assembler) SetJSONOutput(enable bool) {
	a.jsonOutput = enable
}

func (a *Assembler) AddIncludePath(path string) {
	a.includePath = append(a.includePath, path)
}

func (a *Assembler) emitByte(b byte) {
	a.output = append(a.output, b)
	a.currentAddr++
}

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

func (a *Assembler) updateSymbol(name string, value int) error {
	if _, exists := a.symbols[name]; !exists {
		return fmt.Errorf("undefined symbol: %s", name)
	}
	sym := a.symbols[name]
	sym.Value = value
	a.symbols[name] = sym
	return nil
}

func (a *Assembler) addForwardRef(target string, addr int, mode AddressingMode, length int) {
	a.forwardRefs = append(a.forwardRefs, ForwardRef{
		Address: addr,
		Type:    mode,
		Length:  length,
		Target:  target,
	})
}

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

func (a *Assembler) Lookup(mnemonic string) (Instruction, bool) {
	inst, ok := a.instructions[mnemonic]
	return inst, ok
}

func (a *Assembler) GetCurrentAddress() int {
	return a.currentAddr
}

func (a *Assembler) GetOutput() []byte {
	return a.output
}

func (a *Assembler) setOrigin(addr int) {
	a.currentAddr = addr
	a.originSet = true
}

func (a *Assembler) processIncludeFile(filename string) error {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return err
	}
	if a.includes[absPath] {
		return fmt.Errorf("circular include detected: %s", filename)
	}
	a.includes[absPath] = true
	return fmt.Errorf("INCLUDE directive not yet implemented")
}

func (a *Assembler) recordBinaryFile(filename string, skip, length int) error {
	return fmt.Errorf("INCBIN directive not yet implemented")
}

func (a *Assembler) Assemble(filename string) (*AssemblyResult, error) {
	return nil, fmt.Errorf("assembly not yet implemented")
}

func (r *AssemblyResult) WriteFiles(basename string) error {
	return fmt.Errorf("output writing not yet implemented")
}