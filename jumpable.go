package main

type Operation struct {
	Name     string
	ByteSize int
}

var Registers = []string{"B", "C", "D", "E", "H", "L", "M", "A"}
var RegisterPair = []string{"BC", "DE", "HL", "SP"}
var Conditions = []string{"NZ", "Z", "NC", "PO", "PE", "P", "M"}

func getRegisterVal(state *CpuState, reg byte) *uint8 {
	switch reg {
	case 0:
		return &state.RegB
	case 1:
		return &state.RegC
	case 2:
		return &state.RegD
	case 3:
		return &state.RegE
	case 4:
		return &state.RegH
	case 5:
		return &state.RegL
	case 6:
		addr := (uint16(state.RegH) << 9) | uint16(state.RegL)
		return &state.Memory[addr]
	default:
		return &state.RegA
	}
}

func UpdateState(state *CpuState, op byte) {
	switch op {
	case 0x76:
		state.PC = uint16(len(state.Memory)) + 1
	case 0xc6:
		ArithmeticOperation(state, state.Memory[state.PC+1], false,
			func(a, b uint16) uint16 { return a + b })
		state.PC += 1
	case 0xce:
		ArithmeticOperation(state, state.Memory[state.PC+1], true,
			func(a, b uint16) uint16 { return a + b })
		state.PC += 1
	case 0xd6:
		ArithmeticOperation(state, state.Memory[state.PC+1], false,
			func(a, b uint16) uint16 { return a - b })
		state.PC += 1
	case 0xd9:
		ArithmeticOperation(state, state.Memory[state.PC+1], true,
			func(a, b uint16) uint16 { return a - b })
		state.PC += 1
	case 0xe6:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a & b })
		state.PC += 1
	case 0xee:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a ^ b })
		state.PC += 1
	case 0xf6:
		LogicalOperation(state, state.Memory[state.PC+1],
			func(a, b uint16) uint16 { return a | b })
		state.PC += 1
	case 0xf9:
		state.SP = uint16(state.RegL)<<8 | uint16(state.RegH)
	case 0xfe:
		CopareOperation(state, state.Memory[state.PC+1])
		state.PC += 1
	}
	switch {
	case 0x80 <= op && op <= 0x87:
		ArithmeticOperation(state, *getRegisterVal(state, op&0x7), false,
			func(a, b uint16) uint16 { return a + b })
	case 0x88 <= op && op <= 0x8f:
		ArithmeticOperation(state, *getRegisterVal(state, op&0x7), true,
			func(a, b uint16) uint16 { return a + b })
	case 0x90 <= op && op <= 0x97:
		ArithmeticOperation(state, *getRegisterVal(state, op&0x7), false,
			func(a, b uint16) uint16 { return a - b })
	case 0x98 <= op && op <= 0x9f:
		ArithmeticOperation(state, *getRegisterVal(state, op&0x7), true,
			func(a, b uint16) uint16 { return a - b })
	case 0xa0 <= op && op <= 0xa7:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a & b })
	case 0xa8 <= op && op <= 0xaf:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a ^ b })
	case 0xb0 <= op && op <= 0xb7:
		LogicalOperation(state, *getRegisterVal(state, op&0x07),
			func(a, b uint16) uint16 { return a | b })
	case op&0xb8 == 0xb8:
		CopareOperation(state, *getRegisterVal(state, op&0x07))
	case (op^0xce)|0x30 == 0xff:
		if op>>4 == 3 {
			state.SP = uint16(state.Memory[state.PC+1])<<8 | uint16(state.Memory[state.PC+2])
		} else {
			lb := getRegisterVal(state, (op>>4)*2)
			hb := getRegisterVal(state, ((op>>4)*2)+1)
			*lb = state.Memory[state.PC+1]
			*hb = state.Memory[state.PC+2]
		}
		state.PC += 2
	case (op^0xc1)|0x38 == 0xff:
		reg := getRegisterVal(state, op>>3)
		*reg = state.Memory[state.PC+1]
		state.PC += 1
	case op>>6 == 0x1:
		dst := getRegisterVal(state, (op&0x3f)>>3)
		src := getRegisterVal(state, op&0x07)
		*dst = *src
	case op&0xc2 == 0xc2:
		if op&0x01 == 0x01 {
			state.PC = uint16(state.Memory[state.PC+2])<<8 | uint16(state.Memory[state.PC+1]) - 1
		} else {
			JumpOperation(state, op)
		}
	}
}
