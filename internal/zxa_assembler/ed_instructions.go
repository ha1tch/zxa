package zxa_assembler

// initEDInstructions adds all ED-prefixed instructions to the instruction map
func (m InstructionMap) initEDInstructions() {
	// 16-bit load group
	m["LD BC,(nn)"] = Instruction{0x4B, 0xED, Extended, 4, 20, false}
	m["LD DE,(nn)"] = Instruction{0x5B, 0xED, Extended, 4, 20, false}
	m["LD HL,(nn)"] = Instruction{0x6B, 0xED, Extended, 4, 20, false}
	m["LD SP,(nn)"] = Instruction{0x7B, 0xED, Extended, 4, 20, false}
	m["LD (nn),BC"] = Instruction{0x43, 0xED, Extended, 4, 20, false}
	m["LD (nn),DE"] = Instruction{0x53, 0xED, Extended, 4, 20, false}
	m["LD (nn),HL"] = Instruction{0x63, 0xED, Extended, 4, 20, false}
	m["LD (nn),SP"] = Instruction{0x73, 0xED, Extended, 4, 20, false}

	// Block transfer and search group
	m["LDI"] = Instruction{0xA0, 0xED, Implied, 2, 16, false}
	m["LDIR"] = Instruction{0xB0, 0xED, Implied, 2, 21, false}
	m["LDD"] = Instruction{0xA8, 0xED, Implied, 2, 16, false}
	m["LDDR"] = Instruction{0xB8, 0xED, Implied, 2, 21, false}
	m["CPI"] = Instruction{0xA1, 0xED, Implied, 2, 16, false}
	m["CPIR"] = Instruction{0xB1, 0xED, Implied, 2, 21, false}
	m["CPD"] = Instruction{0xA9, 0xED, Implied, 2, 16, false}
	m["CPDR"] = Instruction{0xB9, 0xED, Implied, 2, 21, false}

	// Block I/O group
	m["INI"] = Instruction{0xA2, 0xED, Implied, 2, 16, false}
	m["INIR"] = Instruction{0xB2, 0xED, Implied, 2, 21, false}
	m["IND"] = Instruction{0xAA, 0xED, Implied, 2, 16, false}
	m["INDR"] = Instruction{0xBA, 0xED, Implied, 2, 21, false}
	m["OUTI"] = Instruction{0xA3, 0xED, Implied, 2, 16, false}
	m["OTIR"] = Instruction{0xB3, 0xED, Implied, 2, 21, false}
	m["OUTD"] = Instruction{0xAB, 0xED, Implied, 2, 16, false}
	m["OTDR"] = Instruction{0xBB, 0xED, Implied, 2, 21, false}

	// 16-bit arithmetic group
	m["ADC HL,BC"] = Instruction{0x4A, 0xED, RegisterPair, 2, 15, false}
	m["ADC HL,DE"] = Instruction{0x5A, 0xED, RegisterPair, 2, 15, false}
	m["ADC HL,HL"] = Instruction{0x6A, 0xED, RegisterPair, 2, 15, false}
	m["ADC HL,SP"] = Instruction{0x7A, 0xED, RegisterPair, 2, 15, false}
	m["SBC HL,BC"] = Instruction{0x42, 0xED, RegisterPair, 2, 15, false}
	m["SBC HL,DE"] = Instruction{0x52, 0xED, RegisterPair, 2, 15, false}
	m["SBC HL,HL"] = Instruction{0x62, 0xED, RegisterPair, 2, 15, false}
	m["SBC HL,SP"] = Instruction{0x72, 0xED, RegisterPair, 2, 15, false}

	// Exchange group
	m["NEG"] = Instruction{0x44, 0xED, Implied, 2, 8, false}

	// Interrupt mode and interrupt handling
	m["IM 0"] = Instruction{0x46, 0xED, Implied, 2, 8, false}
	m["IM 1"] = Instruction{0x56, 0xED, Implied, 2, 8, false}
	m["IM 2"] = Instruction{0x5E, 0xED, Implied, 2, 8, false}
	m["RETI"] = Instruction{0x4D, 0xED, Implied, 2, 14, false}
	m["RETN"] = Instruction{0x45, 0xED, Implied, 2, 14, false}

	// I/O group
	m["IN B,(C)"] = Instruction{0x40, 0xED, RegisterIndirect, 2, 12, false}
	m["IN C,(C)"] = Instruction{0x48, 0xED, RegisterIndirect, 2, 12, false}
	m["IN D,(C)"] = Instruction{0x50, 0xED, RegisterIndirect, 2, 12, false}
	m["IN E,(C)"] = Instruction{0x58, 0xED, RegisterIndirect, 2, 12, false}
	m["IN H,(C)"] = Instruction{0x60, 0xED, RegisterIndirect, 2, 12, false}
	m["IN L,(C)"] = Instruction{0x68, 0xED, RegisterIndirect, 2, 12, false}
	m["IN A,(C)"] = Instruction{0x78, 0xED, RegisterIndirect, 2, 12, false}
	m["IN F,(C)"] = Instruction{0x70, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),B"] = Instruction{0x41, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),C"] = Instruction{0x49, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),D"] = Instruction{0x51, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),E"] = Instruction{0x59, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),H"] = Instruction{0x61, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),L"] = Instruction{0x69, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),A"] = Instruction{0x79, 0xED, RegisterIndirect, 2, 12, false}
	m["OUT (C),0"] = Instruction{0x71, 0xED, RegisterIndirect, 2, 12, false}

	// Special register group
	m["LD I,A"] = Instruction{0x47, 0xED, Register, 2, 9, false}
	m["LD R,A"] = Instruction{0x4F, 0xED, Register, 2, 9, false}
	m["LD A,I"] = Instruction{0x57, 0xED, Register, 2, 9, false}
	m["LD A,R"] = Instruction{0x5F, 0xED, Register, 2, 9, false}

	// Special rotate and shift group
	m["RLD"] = Instruction{0x6F, 0xED, Implied, 2, 18, false}
	m["RRD"] = Instruction{0x67, 0xED, Implied, 2, 18, false}
}