package zxa_assembler

// initBaseInstructions adds all non-prefixed instructions to the instruction map
func (m InstructionMap) initBaseInstructions() {
	// 8-bit load group
	m["LD A,A"] = Instruction{0x7F, 0x00, Register, 1, 4, false}
	m["LD A,B"] = Instruction{0x78, 0x00, Register, 1, 4, false}
	m["LD A,C"] = Instruction{0x79, 0x00, Register, 1, 4, false}
	m["LD A,D"] = Instruction{0x7A, 0x00, Register, 1, 4, false}
	m["LD A,E"] = Instruction{0x7B, 0x00, Register, 1, 4, false}
	m["LD A,H"] = Instruction{0x7C, 0x00, Register, 1, 4, false}
	m["LD A,L"] = Instruction{0x7D, 0x00, Register, 1, 4, false}
	m["LD A,(HL)"] = Instruction{0x7E, 0x00, RegisterIndirect, 1, 7, false}
	m["LD A,n"] = Instruction{0x3E, 0x00, Immediate, 2, 7, false}
	m["LD A,(BC)"] = Instruction{0x0A, 0x00, RegisterIndirect, 1, 7, false}
	m["LD A,(DE)"] = Instruction{0x1A, 0x00, RegisterIndirect, 1, 7, false}
	m["LD A,(nn)"] = Instruction{0x3A, 0x00, Extended, 3, 13, false}

	m["LD B,A"] = Instruction{0x47, 0x00, Register, 1, 4, false}
	m["LD B,B"] = Instruction{0x40, 0x00, Register, 1, 4, false}
	m["LD B,C"] = Instruction{0x41, 0x00, Register, 1, 4, false}
	m["LD B,D"] = Instruction{0x42, 0x00, Register, 1, 4, false}
	m["LD B,E"] = Instruction{0x43, 0x00, Register, 1, 4, false}
	m["LD B,H"] = Instruction{0x44, 0x00, Register, 1, 4, false}
	m["LD B,L"] = Instruction{0x45, 0x00, Register, 1, 4, false}
	m["LD B,n"] = Instruction{0x06, 0x00, Immediate, 2, 7, false}

	// 16-bit load group
	m["LD BC,nn"] = Instruction{0x01, 0x00, ImmediateExt, 3, 10, false}
	m["LD DE,nn"] = Instruction{0x11, 0x00, ImmediateExt, 3, 10, false}
	m["LD HL,nn"] = Instruction{0x21, 0x00, ImmediateExt, 3, 10, false}
	m["LD SP,nn"] = Instruction{0x31, 0x00, ImmediateExt, 3, 10, false}
	m["LD SP,HL"] = Instruction{0xF9, 0x00, Register, 1, 6, false}

	// Exchange group
	m["EX DE,HL"] = Instruction{0xEB, 0x00, Implied, 1, 4, false}
	m["EX AF,AF'"] = Instruction{0x08, 0x00, Implied, 1, 4, false}
	m["EXX"] = Instruction{0xD9, 0x00, Implied, 1, 4, false}
	m["EX (SP),HL"] = Instruction{0xE3, 0x00, RegisterIndirect, 1, 19, false}

	// 8-bit arithmetic and logical group
	m["ADD A,A"] = Instruction{0x87, 0x00, Register, 1, 4, false}
	m["ADD A,B"] = Instruction{0x80, 0x00, Register, 1, 4, false}
	m["ADD A,n"] = Instruction{0xC6, 0x00, Immediate, 2, 7, false}
	m["ADC A,A"] = Instruction{0x8F, 0x00, Register, 1, 4, false}
	m["SUB A"] = Instruction{0x97, 0x00, Register, 1, 4, false}
	m["SUB n"] = Instruction{0xD6, 0x00, Immediate, 2, 7, false}
	m["AND A"] = Instruction{0xA7, 0x00, Register, 1, 4, false}
	m["AND n"] = Instruction{0xE6, 0x00, Immediate, 2, 7, false}
	m["OR A"] = Instruction{0xB7, 0x00, Register, 1, 4, false}
	m["OR n"] = Instruction{0xF6, 0x00, Immediate, 2, 7, false}
	m["XOR A"] = Instruction{0xAF, 0x00, Register, 1, 4, false}
	m["XOR n"] = Instruction{0xEE, 0x00, Immediate, 2, 7, false}
	m["CP A"] = Instruction{0xBF, 0x00, Register, 1, 4, false}
	m["CP n"] = Instruction{0xFE, 0x00, Immediate, 2, 7, false}
	m["INC A"] = Instruction{0x3C, 0x00, Register, 1, 4, false}
	m["DEC A"] = Instruction{0x3D, 0x00, Register, 1, 4, false}

	// Jump group
	m["JP nn"] = Instruction{0xC3, 0x00, Extended, 3, 10, false}
	m["JP NZ,nn"] = Instruction{0xC2, 0x00, Extended, 3, 10, true}
	m["JP Z,nn"] = Instruction{0xCA, 0x00, Extended, 3, 10, true}
	m["JP NC,nn"] = Instruction{0xD2, 0x00, Extended, 3, 10, true}
	m["JP C,nn"] = Instruction{0xDA, 0x00, Extended, 3, 10, true}
	m["JR e"] = Instruction{0x18, 0x00, Relative, 2, 12, false}
	m["JR NZ,e"] = Instruction{0x20, 0x00, Relative, 2, 12, true}
	m["JR Z,e"] = Instruction{0x28, 0x00, Relative, 2, 12, true}
	m["JR NC,e"] = Instruction{0x30, 0x00, Relative, 2, 12, true}
	m["JR C,e"] = Instruction{0x38, 0x00, Relative, 2, 12, true}
	m["DJNZ e"] = Instruction{0x10, 0x00, Relative, 2, 13, false}

	// Call and return group
	m["CALL nn"] = Instruction{0xCD, 0x00, Extended, 3, 17, false}
	m["CALL NZ,nn"] = Instruction{0xC4, 0x00, Extended, 3, 17, true}
	m["CALL Z,nn"] = Instruction{0xCC, 0x00, Extended, 3, 17, true}
	m["CALL NC,nn"] = Instruction{0xD4, 0x00, Extended, 3, 17, true}
	m["CALL C,nn"] = Instruction{0xDC, 0x00, Extended, 3, 17, true}
	m["RET"] = Instruction{0xC9, 0x00, Implied, 1, 10, false}
	m["RET NZ"] = Instruction{0xC0, 0x00, Implied, 1, 11, true}
	m["RET Z"] = Instruction{0xC8, 0x00, Implied, 1, 11, true}
	m["RET NC"] = Instruction{0xD0, 0x00, Implied, 1, 11, true}
	m["RET C"] = Instruction{0xD8, 0x00, Implied, 1, 11, true}

	// RST group
	m["RST 00H"] = Instruction{0xC7, 0x00, Implied, 1, 11, false}
	m["RST 08H"] = Instruction{0xCF, 0x00, Implied, 1, 11, false}
	m["RST 10H"] = Instruction{0xD7, 0x00, Implied, 1, 11, false}
	m["RST 18H"] = Instruction{0xDF, 0x00, Implied, 1, 11, false}
	m["RST 20H"] = Instruction{0xE7, 0x00, Implied, 1, 11, false}
	m["RST 28H"] = Instruction{0xEF, 0x00, Implied, 1, 11, false}
	m["RST 30H"] = Instruction{0xF7, 0x00, Implied, 1, 11, false}
	m["RST 38H"] = Instruction{0xFF, 0x00, Implied, 1, 11, false}
}