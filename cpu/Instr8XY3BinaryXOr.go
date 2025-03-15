package cpu

import "fmt"

type Instr8XY3BinaryXOr struct {
	InstrParamsXY
}

func (instr *Instr8XY3BinaryXOr) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	var newValue = xValue ^ yValue
	cpu.vRegisters[instr.x] = newValue

	fmt.Printf("8XY3_Binary_XOR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", instr.x, instr.y, xValue, yValue, newValue, cpu.vRegisters[instr.x])
}
