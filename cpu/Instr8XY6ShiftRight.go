package cpu

import "fmt"

type Instr8XY6ShiftRight struct {
	InstrParamsXY
}

func (instr *Instr8XY6ShiftRight) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var newValue = xValue >> 1
	var shiftedBit = xValue & 0x1
	cpu.vRegisters[instr.x] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}

	fmt.Printf("8XY6_Shift_Right, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", instr.x, instr.y, xValue, newValue, shiftedBit, isCarryFlagSet)
}
