# zxa - Z80 cross assembler

A simple Z80 cross-assembler for the ZX Spectrum written in Go, featuring a clean syntax and modern development workflow features.

## Features

- Aspires to full Z80 instruction set support (possibly incomplete as of Jan 2025)
- Z80N variant support for the ZX Spectrum Next
- Assembler directives (ORG, INCLUDE, INCBIN, etc.)
- Multiple output formats (binary, hex dump, JSON)
- Detailed error reporting with source locations
- Support for labels and constants
- Include file support for modular code
- Binary file inclusion support

### Missing, may be added:
- Integration with plus3 (https://github.com/ha1tch/plus3)
- Integration with zxgotools (https://github.com/ha1tch/zxgotools)
- Macros
- Assembler explainer with Ollama

## Installation

```bash
go install github.com/ha1tch/zxa/cmd/zxa@latest
```

## Quick Start

1. Create an assembly file (example.asm):
```assembly
        ORG $8000
start:  LD A,42
        LD B,10
loop:   DJNZ loop
        RET
```

2. Assemble the file:
```bash
zxa example.asm
```

3. Check the outputs:
- example.bin (binary output)
- example.hex (hex dump if --hex specified)
- example.json (assembly report if --json specified)

## Usage

```bash
zxa [options] <input.asm>

Options:
  --next                Enable ZX Spectrum Next's Z80N processor support 
  -o, --output string   Output file name (default: input base name)
  -I, --include string  Add include search path
  --hex                 Generate hex dump output
  --json                Generate JSON assembly report
  -v, --verbose         Enable verbose output
  -q                    Quiet mode (suppress non-error output)
  --version             Show version information
```

## Directives

- `ORG address` - Set the origin address
- `INCLUDE "file"` - Include source file
- `INCBIN "file"[,skip[,length]]` - Include binary file
- `DEFB expressions,...` - Define bytes
- `DEFW expressions,...` - Define words
- `DEFS length[,fill]` - Define storage
- `label: EQU expression` - Define constant

## Error Handling

The assembler provides detailed error messages with categories:

- Syntax errors
- Symbol errors (undefined, duplicate)
- Value errors (out of range)
- File handling errors
- Directive errors
- Range errors (jump out of range)

## Library Usage

The assembler can also be used as a Go library:

```go
package main

import "github.com/ha1tch/zxa/internal/zxa_assembler"

func main() {
    opts := zxa_assembler.AssemblerOptions{
        Variant: zxa_assembler.Z80Standard,
    }
    
    instructions, err := zxa_assembler.NewInstructionSet(opts)
    if err != nil {
        // Handle error
    }
    
    asm := zxa_assembler.NewAssembler(instructions)
    asm.SetHexOutput(true)
    asm.SetJSONOutput(true)
    
    result, err := asm.Assemble("program.asm")
    if err != nil {
        // Handle error
    }
}
```

## Examples

Check the `examples/` directory for complete example programs:

- `hello/` - Simple hello world program
- `game/` - Basic game example with sprites

## Testing

Run the test suite:

```bash
go test ./...
```

### Contact:
haitch@duck.com
https://oldbytes.space/@haitchfive

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
