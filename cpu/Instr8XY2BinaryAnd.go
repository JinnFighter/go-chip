package cpu

import "fmt"

type Instr8XY2BinaryAnd struct {
	InstrParamsXY
}

func (instr *Instr8XY2BinaryAnd) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	var newValue = xValue & yValue
	cpu.vRegisters[instr.x] = newValue
	fmt.Printf("8XY2_Binary_AND, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", instr.x, instr.y, xValue, yValue, newValue, cpu.vRegisters[instr.x])
}
