// file: /cmd/zxa/internal/zxa_assembler/instruction.go

package zxa_assembler

import (
	"encoding/json"
	"fmt"
	"os"
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

type Statistics struct {
	BytesGenerated int `json:"bytesGenerated"`
	LinesProcessed int `json:"linesProcessed"`
	SymbolsDefined int `json:"symbolsDefined"`
}

type AssemblyResult struct {
	Success    bool
	Output     []byte
	HexOutput  bool
	JSONOutput bool
	Symbols    map[string]Symbol
	Statistics Statistics
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

	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read include file: %v", err)
	}

	parser := NewParser(string(content))
	parser.assembler = a

	for {
		err := parser.parseLine()
		if err != nil {
			return fmt.Errorf("error parsing include file: %v", err)
		}
		if parser.pos >= len(parser.input) {
			break
		}
	}

	return nil
}

func (a *Assembler) recordBinaryFile(filename string, skip, length int) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read binary file: %v", err)
	}

	if skip >= len(content) {
		return fmt.Errorf("skip offset beyond file size")
	}

	if length < 0 {
		length = len(content) - skip
	}

	if skip+length > len(content) {
		return fmt.Errorf("requested length exceeds file size")
	}

	for i := 0; i < length; i++ {
		a.emitByte(content[skip+i])
	}

	return nil
}

func (a *Assembler) Assemble(filename string) (*AssemblyResult, error) {
	result := &AssemblyResult{
		Success:    true,
		Output:     nil,
		HexOutput:  a.hexOutput,
		JSONOutput: a.jsonOutput,
		Symbols:    a.symbols,
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}

	parser := NewParser(string(content))
	parser.assembler = a

	// First pass: Parse symbols and resolve forward references
	for {
		err := parser.parseLine()
		if err != nil {
			return nil, fmt.Errorf("parsing error: %v", err)
		}

		if parser.pos >= len(parser.input) {
			break
		}
	}

	// Reset for second pass
	a.output = make([]byte, 0, 1024)
	a.currentAddr = 0
	a.originSet = false
	parser.pos = 0
	parser.line = 1
	parser.column = 1

	// Second pass: Generate code
	for {
		err := parser.parseLine()
		if err != nil {
			return nil, fmt.Errorf("code generation error: %v", err)
		}

		result.Statistics.LinesProcessed++

		if parser.pos >= len(parser.input) {
			break
		}
	}

	if err := a.resolveForwardRefs(); err != nil {
		return nil, fmt.Errorf("forward reference resolution error: %v", err)
	}

	result.Statistics.BytesGenerated = len(a.output)
	result.Statistics.SymbolsDefined = len(a.symbols)
	result.Output = a.output

	return result, nil
}

func (r *AssemblyResult) WriteFiles(basename string) error {
	binFile := basename + ".bin"
	if err := os.WriteFile(binFile, r.Output, 0644); err != nil {
		return fmt.Errorf("failed to write binary output: %v", err)
	}

	if r.HexOutput {
		hexFile := basename + ".hex"
		f, err := os.Create(hexFile)
		if err != nil {
			return fmt.Errorf("failed to create hex output file: %v", err)
		}
		defer f.Close()

		for i := 0; i < len(r.Output); i += 16 {
			fmt.Fprintf(f, "%04X: ", i)
			
			for j := 0; j < 16; j++ {
				if i+j < len(r.Output) {
					fmt.Fprintf(f, "%02X ", r.Output[i+j])
				} else {
					fmt.Fprintf(f, "   ")
				}
			}
			
			fmt.Fprintf(f, " |")
			for j := 0; j < 16 && i+j < len(r.Output); j++ {
				b := r.Output[i+j]
				if b >= 32 && b <= 126 {
					fmt.Fprintf(f, "%c", b)
				} else {
					fmt.Fprintf(f, ".")
				}
			}
			fmt.Fprintf(f, "|\n")
		}
	}

	if r.JSONOutput {
		jsonFile := basename + ".json"
		report := struct {
			Statistics Statistics         `json:"statistics"`
			Symbols    map[string]Symbol `json:"symbols"`
		}{
			Statistics: r.Statistics,
			Symbols:    r.Symbols,
		}
		
		jsonData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to generate JSON report: %v", err)
		}
		
		if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
			return fmt.Errorf("failed to write JSON report: %v", err)
		}
	}

	return nil
}