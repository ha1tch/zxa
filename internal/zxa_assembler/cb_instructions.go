package zxa_assembler

import "fmt"

// initCBInstructions adds all CB-prefixed instructions to the instruction map
func (m InstructionMap) initCBInstructions() {
	// Rotation instructions
	m["RLC B"] = Instruction{0x00, 0xCB, Register, 2, 8, false}
	m["RLC C"] = Instruction{0x01, 0xCB, Register, 2, 8, false}
	m["RLC D"] = Instruction{0x02, 0xCB, Register, 2, 8, false}
	m["RLC E"] = Instruction{0x03, 0xCB, Register, 2, 8, false}
	m["RLC H"] = Instruction{0x04, 0xCB, Register, 2, 8, false}
	m["RLC L"] = Instruction{0x05, 0xCB, Register, 2, 8, false}
	m["RLC (HL)"] = Instruction{0x06, 0xCB, RegisterIndirect, 2, 15, false}
	m["RLC A"] = Instruction{0x07, 0xCB, Register, 2, 8, false}

	m["RRC B"] = Instruction{0x08, 0xCB, Register, 2, 8, false}
	m["RRC C"] = Instruction{0x09, 0xCB, Register, 2, 8, false}
	m["RRC D"] = Instruction{0x0A, 0xCB, Register, 2, 8, false}
	m["RRC E"] = Instruction{0x0B, 0xCB, Register, 2, 8, false}
	m["RRC H"] = Instruction{0x0C, 0xCB, Register, 2, 8, false}
	m["RRC L"] = Instruction{0x0D, 0xCB, Register, 2, 8, false}
	m["RRC (HL)"] = Instruction{0x0E, 0xCB, RegisterIndirect, 2, 15, false}
	m["RRC A"] = Instruction{0x0F, 0xCB, Register, 2, 8, false}

	m["RL B"] = Instruction{0x10, 0xCB, Register, 2, 8, false}
	m["RL C"] = Instruction{0x11, 0xCB, Register, 2, 8, false}
	m["RL D"] = Instruction{0x12, 0xCB, Register, 2, 8, false}
	m["RL E"] = Instruction{0x13, 0xCB, Register, 2, 8, false}
	m["RL H"] = Instruction{0x14, 0xCB, Register, 2, 8, false}
	m["RL L"] = Instruction{0x15, 0xCB, Register, 2, 8, false}
	m["RL (HL)"] = Instruction{0x16, 0xCB, RegisterIndirect, 2, 15, false}
	m["RL A"] = Instruction{0x17, 0xCB, Register, 2, 8, false}

	m["RR B"] = Instruction{0x18, 0xCB, Register, 2, 8, false}
	m["RR C"] = Instruction{0x19, 0xCB, Register, 2, 8, false}
	m["RR D"] = Instruction{0x1A, 0xCB, Register, 2, 8, false}
	m["RR E"] = Instruction{0x1B, 0xCB, Register, 2, 8, false}
	m["RR H"] = Instruction{0x1C, 0xCB, Register, 2, 8, false}
	m["RR L"] = Instruction{0x1D, 0xCB, Register, 2, 8, false}
	m["RR (HL)"] = Instruction{0x1E, 0xCB, RegisterIndirect, 2, 15, false}
	m["RR A"] = Instruction{0x1F, 0xCB, Register, 2, 8, false}

	// Shift instructions
	m["SLA B"] = Instruction{0x20, 0xCB, Register, 2, 8, false}
	m["SLA C"] = Instruction{0x21, 0xCB, Register, 2, 8, false}
	m["SLA D"] = Instruction{0x22, 0xCB, Register, 2, 8, false}
	m["SLA E"] = Instruction{0x23, 0xCB, Register, 2, 8, false}
	m["SLA H"] = Instruction{0x24, 0xCB, Register, 2, 8, false}
	m["SLA L"] = Instruction{0x25, 0xCB, Register, 2, 8, false}
	m["SLA (HL)"] = Instruction{0x26, 0xCB, RegisterIndirect, 2, 15, false}
	m["SLA A"] = Instruction{0x27, 0xCB, Register, 2, 8, false}

	m["SRA B"] = Instruction{0x28, 0xCB, Register, 2, 8, false}
	m["SRA C"] = Instruction{0x29, 0xCB, Register, 2, 8, false}
	m["SRA D"] = Instruction{0x2A, 0xCB, Register, 2, 8, false}
	m["SRA E"] = Instruction{0x2B, 0xCB, Register, 2, 8, false}
	m["SRA H"] = Instruction{0x2C, 0xCB, Register, 2, 8, false}
	m["SRA L"] = Instruction{0x2D, 0xCB, Register, 2, 8, false}
	m["SRA (HL)"] = Instruction{0x2E, 0xCB, RegisterIndirect, 2, 15, false}
	m["SRA A"] = Instruction{0x2F, 0xCB, Register, 2, 8, false}

	m["SRL B"] = Instruction{0x38, 0xCB, Register, 2, 8, false}
	m["SRL C"] = Instruction{0x39, 0xCB, Register, 2, 8, false}
	m["SRL D"] = Instruction{0x3A, 0xCB, Register, 2, 8, false}
	m["SRL E"] = Instruction{0x3B, 0xCB, Register, 2, 8, false}
	m["SRL H"] = Instruction{0x3C, 0xCB, Register, 2, 8, false}
	m["SRL L"] = Instruction{0x3D, 0xCB, Register, 2, 8, false}
	m["SRL (HL)"] = Instruction{0x3E, 0xCB, RegisterIndirect, 2, 15, false}
	m["SRL A"] = Instruction{0x3F, 0xCB, Register, 2, 8, false}

	// BIT instructions - Test bit b in register r
	for bit := 0; bit < 8; bit++ {
		base := 0x40 + (bit << 3)
		m[fmt.Sprintf("BIT %d,B", bit)] = Instruction{byte(base), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,C", bit)] = Instruction{byte(base + 1), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,D", bit)] = Instruction{byte(base + 2), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,E", bit)] = Instruction{byte(base + 3), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,H", bit)] = Instruction{byte(base + 4), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,L", bit)] = Instruction{byte(base + 5), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("BIT %d,(HL)", bit)] = Instruction{byte(base + 6), 0xCB, BitIndex, 2, 12, false}
		m[fmt.Sprintf("BIT %d,A", bit)] = Instruction{byte(base + 7), 0xCB, BitIndex, 2, 8, false}
	}

	// RES instructions - Reset bit b in register r
	for bit := 0; bit < 8; bit++ {
		base := 0x80 + (bit << 3)
		m[fmt.Sprintf("RES %d,B", bit)] = Instruction{byte(base), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,C", bit)] = Instruction{byte(base + 1), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,D", bit)] = Instruction{byte(base + 2), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,E", bit)] = Instruction{byte(base + 3), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,H", bit)] = Instruction{byte(base + 4), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,L", bit)] = Instruction{byte(base + 5), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("RES %d,(HL)", bit)] = Instruction{byte(base + 6), 0xCB, BitIndex, 2, 15, false}
		m[fmt.Sprintf("RES %d,A", bit)] = Instruction{byte(base + 7), 0xCB, BitIndex, 2, 8, false}
	}

	// SET instructions - Set bit b in register r
	for bit := 0; bit < 8; bit++ {
		base := 0xC0 + (bit << 3)
		m[fmt.Sprintf("SET %d,B", bit)] = Instruction{byte(base), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,C", bit)] = Instruction{byte(base + 1), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,D", bit)] = Instruction{byte(base + 2), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,E", bit)] = Instruction{byte(base + 3), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,H", bit)] = Instruction{byte(base + 4), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,L", bit)] = Instruction{byte(base + 5), 0xCB, BitIndex, 2, 8, false}
		m[fmt.Sprintf("SET %d,(HL)", bit)] = Instruction{byte(base + 6), 0xCB, BitIndex, 2, 15, false}
		m[fmt.Sprintf("SET %d,A", bit)] = Instruction{byte(base + 7), 0xCB, BitIndex, 2, 8, false}
	}
}
