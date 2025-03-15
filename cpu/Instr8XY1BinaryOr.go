package cpu

import "fmt"

type Instr8XY1BinaryOr struct {
	InstrParamsXY
}

func (instr *Instr8XY1BinaryOr) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	var newValue = xValue | yValue
	cpu.vRegisters[instr.xIdx] = newValue

	fmt.Printf("8XY1_Binary_OR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", instr.xIdx, instr.yIdx, xValue, yValue, newValue, cpu.vRegisters[instr.xIdx])
}
