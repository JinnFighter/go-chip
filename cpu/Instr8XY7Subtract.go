package cpu

import "fmt"

type Instr8XY7Subtract struct {
	InstrParamsXY
}

func (instr *Instr8XY7Subtract) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	var isCarryFlagSet = yValue >= xValue
	var oldValue = cpu.vRegisters[instr.x]
	var newValue = yValue - xValue
	cpu.vRegisters[instr.x] = newValue

	if isCarryFlagSet {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}

	fmt.Printf("8XY7_Subtract, x: %d, y: %d, oldValue: %b, newValue: %b, isCarryFlagSet: %t \n", xValue, yValue, oldValue, newValue, isCarryFlagSet)
}
