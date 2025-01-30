package zxa_assembler

// initIndexInstructions adds all DD/FD-prefixed (IX/IY) instructions to the instruction map
func (m InstructionMap) initIndexInstructions() {
	// IX instructions (DD prefix)
	m["ADD IX,BC"] = Instruction{0x09, 0xDD, RegisterPair, 2, 15, false}
	m["ADD IX,DE"] = Instruction{0x19, 0xDD, RegisterPair, 2, 15, false}
	m["ADD IX,IX"] = Instruction{0x29, 0xDD, RegisterPair, 2, 15, false}
	m["ADD IX,SP"] = Instruction{0x39, 0xDD, RegisterPair, 2, 15, false}

	m["LD IX,nn"] = Instruction{0x21, 0xDD, ImmediateExt, 4, 14, false}
	m["LD (nn),IX"] = Instruction{0x22, 0xDD, Extended, 4, 20, false}
	m["LD IX,(nn)"] = Instruction{0x2A, 0xDD, Extended, 4, 20, false}
	m["LD SP,IX"] = Instruction{0xF9, 0xDD, Register, 2, 10, false}

	m["INC IX"] = Instruction{0x23, 0xDD, Register, 2, 10, false}
	m["DEC IX"] = Instruction{0x2B, 0xDD, Register, 2, 10, false}

	// IX with displacement (d)
	m["LD (IX+d),n"] = Instruction{0x36, 0xDD, Indexed, 4, 19, false}
	m["LD (IX+d),B"] = Instruction{0x70, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),C"] = Instruction{0x71, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),D"] = Instruction{0x72, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),E"] = Instruction{0x73, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),H"] = Instruction{0x74, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),L"] = Instruction{0x75, 0xDD, Indexed, 3, 19, false}
	m["LD (IX+d),A"] = Instruction{0x77, 0xDD, Indexed, 3, 19, false}

	m["LD B,(IX+d)"] = Instruction{0x46, 0xDD, Indexed, 3, 19, false}
	m["LD C,(IX+d)"] = Instruction{0x4E, 0xDD, Indexed, 3, 19, false}
	m["LD D,(IX+d)"] = Instruction{0x56, 0xDD, Indexed, 3, 19, false}
	m["LD E,(IX+d)"] = Instruction{0x5E, 0xDD, Indexed, 3, 19, false}
	m["LD H,(IX+d)"] = Instruction{0x66, 0xDD, Indexed, 3, 19, false}
	m["LD L,(IX+d)"] = Instruction{0x6E, 0xDD, Indexed, 3, 19, false}
	m["LD A,(IX+d)"] = Instruction{0x7E, 0xDD, Indexed, 3, 19, false}

	// IY instructions (FD prefix) - Mirror of IX instructions
	m["ADD IY,BC"] = Instruction{0x09, 0xFD, RegisterPair, 2, 15, false}
	m["ADD IY,DE"] = Instruction{0x19, 0xFD, RegisterPair, 2, 15, false}
	m["ADD IY,IY"] = Instruction{0x29, 0xFD, RegisterPair, 2, 15, false}
	m["ADD IY,SP"] = Instruction{0x39, 0xFD, RegisterPair, 2, 15, false}

	m["LD IY,nn"] = Instruction{0x21, 0xFD, ImmediateExt, 4, 14, false}
	m["LD (nn),IY"] = Instruction{0x22, 0xFD, Extended, 4, 20, false}
	m["LD IY,(nn)"] = Instruction{0x2A, 0xFD, Extended, 4, 20, false}
	m["LD SP,IY"] = Instruction{0xF9, 0xFD, Register, 2, 10, false}

	m["INC IY"] = Instruction{0x23, 0xFD, Register, 2, 10, false}
	m["DEC IY"] = Instruction{0x2B, 0xFD, Register, 2, 10, false}

	// IY with displacement (d)
	m["LD (IY+d),n"] = Instruction{0x36, 0xFD, Indexed, 4, 19, false}
	m["LD (IY+d),B"] = Instruction{0x70, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),C"] = Instruction{0x71, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),D"] = Instruction{0x72, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),E"] = Instruction{0x73, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),H"] = Instruction{0x74, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),L"] = Instruction{0x75, 0xFD, Indexed, 3, 19, false}
	m["LD (IY+d),A"] = Instruction{0x77, 0xFD, Indexed, 3, 19, false}

	m["LD B,(IY+d)"] = Instruction{0x46, 0xFD, Indexed, 3, 19, false}
	m["LD C,(IY+d)"] = Instruction{0x4E, 0xFD, Indexed, 3, 19, false}
	m["LD D,(IY+d)"] = Instruction{0x56, 0xFD, Indexed, 3, 19, false}
	m["LD E,(IY+d)"] = Instruction{0x5E, 0xFD, Indexed, 3, 19, false}
	m["LD H,(IY+d)"] = Instruction{0x66, 0xFD, Indexed, 3, 19, false}
	m["LD L,(IY+d)"] = Instruction{0x6E, 0xFD, Indexed, 3, 19, false}
	m["LD A,(IY+d)"] = Instruction{0x7E, 0xFD, Indexed, 3, 19, false}

	// Arithmetic and logic with indexed addressing
	// For IX
	m["ADD A,(IX+d)"] = Instruction{0x86, 0xDD, Indexed, 3, 19, false}
	m["ADC A,(IX+d)"] = Instruction{0x8E, 0xDD, Indexed, 3, 19, false}
	m["SUB (IX+d)"] = Instruction{0x96, 0xDD, Indexed, 3, 19, false}
	m["SBC A,(IX+d)"] = Instruction{0x9E, 0xDD, Indexed, 3, 19, false}
	m["AND (IX+d)"] = Instruction{0xA6, 0xDD, Indexed, 3, 19, false}
	m["XOR (IX+d)"] = Instruction{0xAE, 0xDD, Indexed, 3, 19, false}
	m["OR (IX+d)"] = Instruction{0xB6, 0xDD, Indexed, 3, 19, false}
	m["CP (IX+d)"] = Instruction{0xBE, 0xDD, Indexed, 3, 19, false}

	// For IY
	m["ADD A,(IY+d)"] = Instruction{0x86, 0xFD, Indexed, 3, 19, false}
	m["ADC A,(IY+d)"] = Instruction{0x8E, 0xFD, Indexed, 3, 19, false}
	m["SUB (IY+d)"] = Instruction{0x96, 0xFD, Indexed, 3, 19, false}
	m["SBC A,(IY+d)"] = Instruction{0x9E, 0xFD, Indexed, 3, 19, false}
	m["AND (IY+d)"] = Instruction{0xA6, 0xFD, Indexed, 3, 19, false}
	m["XOR (IY+d)"] = Instruction{0xAE, 0xFD, Indexed, 3, 19, false}
	m["OR (IY+d)"] = Instruction{0xB6, 0xFD, Indexed, 3, 19, false}
	m["CP (IY+d)"] = Instruction{0xBE, 0xFD, Indexed, 3, 19, false}

	// Inc/Dec indexed memory
	m["INC (IX+d)"] = Instruction{0x34, 0xDD, Indexed, 3, 23, false}
	m["DEC (IX+d)"] = Instruction{0x35, 0xDD, Indexed, 3, 23, false}
	m["INC (IY+d)"] = Instruction{0x34, 0xFD, Indexed, 3, 23, false}
	m["DEC (IY+d)"] = Instruction{0x35, 0xFD, Indexed, 3, 23, false}
}