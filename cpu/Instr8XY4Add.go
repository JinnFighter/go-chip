package cpu

import "fmt"

type Instr8XY4Add struct {
	InstrParamsXY
}

func (instr *Instr8XY4Add) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	var isCarryFlagSet = (int(xValue) + int(yValue)) > 255
	var oldValue = cpu.vRegisters[instr.x]
	var newValue = xValue + yValue
	cpu.vRegisters[instr.x] = newValue

	if isCarryFlagSet {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}

	fmt.Printf("8XY4_Add, x: %d, y: %d, oldValue: %b, newValue: %b, isCarryFlagSet: %t \n", xValue, yValue, oldValue, newValue, isCarryFlagSet)
}
