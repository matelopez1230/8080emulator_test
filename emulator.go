package main

type CpuFlags struct {
	Parity   bool
	Zero     bool
	Sign     bool
	Carry    bool
	AuxCarry bool
}

type CpuState struct {
	RegA      uint8
	RegB      uint8
	RegC      uint8
	RegD      uint8
	RegE      uint8
	RegH      uint8
	RegL      uint8
	SP        uint16
	PC        uint16
	Memory    []byte
	IntEnable bool
	Condition CpuFlags
}

type Operator func(a uint16, b uint16) uint16

func bitParity(byte uint16) bool {
	var acc uint16 = 0
	for i := 0; i < 8; i++ {
		acc += (byte >> i) & 0x01
	}
	return acc%2 == 0
}

func ArithmeticOperation(state *CpuState, value uint8, useCarry bool, op Operator) {
	var answer uint16 = op(uint16(state.RegA), uint16(value))

	if state.Condition.Carry && useCarry {
		answer = op(answer, 1)
	}

	//cpu flags
	state.Condition.Zero = ((answer & 0xff) == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = (answer > 0xff)
	state.Condition.Parity = bitParity(answer & 0xff)
	state.Condition.AuxCarry = (((state.RegA & value) & 0x04) == 0x04)

	state.RegA = uint8(answer)
}

func LogicalOperation(state *CpuState, value uint8, op Operator) {
	var answer uint16 = uint16(state.RegA) - uint16(value)

	// CpuFlags
	state.Condition.Zero = (answer == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = false
	state.Condition.Parity = bitParity(answer & 0xff)
}

func CopareOperation(state *CpuState, value uint8) {
	var answer uint16 = uint16(state.RegA) - uint16(value)

	//CpuFlags
	state.Condition.Zero = ((answer & 0xff) == 0)
	state.Condition.Sign = ((answer & 0x80) != 0)
	state.Condition.Carry = (answer > 0xff)
	state.Condition.Parity = bitParity(answer & 0xff)
	state.Condition.AuxCarry = (((state.RegA & value) & 0x04) == 0x04)
}

func JumpOperation(state *CpuState, op byte) {
	var answer bool

	switch (op & 0x38) >> 3 {
	case 0:
		answer = state.Condition.Zero
	case 1:
		answer = !state.Condition.Zero
	case 2:
		answer = !state.Condition.Carry
	case 3:
		answer = state.Condition.Carry
	case 4:
		answer = !bitParity(uint16(op))
	case 5:
		answer = bitParity(uint16(op))
	case 6:
		answer = !state.Condition.Sign
	case 7:
		answer = state.Condition.Sign
	}

	if answer {
		state.PC = uint16(state.Memory[state.PC+2])<<8 | uint16(state.Memory[state.PC+1]) - 1
	} else {
		state.PC += 2
	}
}
