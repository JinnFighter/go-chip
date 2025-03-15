package cpu

import "fmt"

type Instr8XY5Subtract struct {
	InstrParamsXY
}

func (instr *Instr8XY5Subtract) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	var isCarryFlagSet = xValue >= yValue
	var oldValue = cpu.vRegisters[instr.xIdx]
	var newValue = xValue - yValue
	cpu.vRegisters[instr.xIdx] = newValue

	if isCarryFlagSet {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}

	fmt.Printf("8XY5_Subtract, x: %d, y: %d, oldValue: %b, newValue: %b, isCarryFlagSet: %t \n", xValue, yValue, oldValue, newValue, isCarryFlagSet)
}
