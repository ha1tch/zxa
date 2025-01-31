package zxa_assembler

// initZ80NInstructions adds Z80N (ZX Spectrum Next) specific instructions
func (m InstructionMap) initZ80NInstructions() {
	// Addition with carry of registers pairs
	m["ADD BC,A"] = Instruction{0x76, 0xED, RegisterPair, 2, 8, false}
	m["ADD DE,A"] = Instruction{0x77, 0xED, RegisterPair, 2, 8, false}
	m["ADD HL,A"] = Instruction{0x7C, 0xED, RegisterPair, 2, 8, false}

	// 16-bit load immediate pseudo-instructions
	m["PUSH nn"] = Instruction{0x8A, 0xED, ImmediateExt, 4, 12, false}

	// Misc Z80N instructions
	m["OUTINB"] = Instruction{0x90, 0xED, Implied, 2, 16, false}
	m["MUL"] = Instruction{0x30, 0xED, Implied, 2, 8, false}
	m["SWAPNIB"] = Instruction{0x23, 0xED, Implied, 2, 8, false}
	m["MIRROR"] = Instruction{0x24, 0xED, Implied, 2, 8, false}
	m["NEXTREG n,n"] = Instruction{0x91, 0xED, Immediate, 4, 20, false}
	m["NEXTREG n,A"] = Instruction{0x92, 0xED, Immediate, 3, 17, false}
	m["PIXELDN"] = Instruction{0x93, 0xED, Implied, 2, 8, false}
	m["PIXELAD"] = Instruction{0x94, 0xED, Implied, 2, 8, false}
	m["SETAE"] = Instruction{0x95, 0xED, Implied, 2, 8, false}
	m["TEST n"] = Instruction{0x27, 0xED, Immediate, 3, 11, false}
	m["BSLA DE,B"] = Instruction{0x28, 0xED, RegisterPair, 2, 8, false}
	m["BSRA DE,B"] = Instruction{0x29, 0xED, RegisterPair, 2, 8, false}
	m["BSRL DE,B"] = Instruction{0x2A, 0xED, RegisterPair, 2, 8, false}
	m["BSRF DE,B"] = Instruction{0x2B, 0xED, RegisterPair, 2, 8, false}
	m["BRLC DE,B"] = Instruction{0x2C, 0xED, RegisterPair, 2, 8, false}

	// Copper instructions
	m["CUP"] = Instruction{0xB5, 0xED, Implied, 2, 8, false}

	// DMA instructions
	m["LDIX"] = Instruction{0xA4, 0xED, Implied, 2, 16, false}
	m["LDIRX"] = Instruction{0xB4, 0xED, Implied, 2, 16, false}
	m["LDDX"] = Instruction{0xAC, 0xED, Implied, 2, 16, false}
	m["LDDRX"] = Instruction{0xBC, 0xED, Implied, 2, 16, false}
	m["LDPIRX"] = Instruction{0xB7, 0xED, Implied, 2, 16, false}
}
