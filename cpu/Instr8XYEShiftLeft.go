package cpu

import "fmt"

type Instr8XYEShiftLeft struct {
	InstrParamsXY
}

func (instr *Instr8XYEShiftLeft) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var newValue = xValue << 1
	var shiftedBit = xValue >> 7
	cpu.vRegisters[instr.x] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}

	fmt.Printf("8XYE_Shift_Left, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", instr.x, instr.y, xValue, newValue, shiftedBit, isCarryFlagSet)
}
