package cpu

import "fmt"

type Instr8XY2BinaryAnd struct {
	InstrParamsXY
}

func (instr *Instr8XY2BinaryAnd) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	var newValue = xValue & yValue
	cpu.vRegisters[instr.xIdx] = newValue
	fmt.Printf("8XY2_Binary_AND, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", instr.xIdx, instr.yIdx, xValue, yValue, newValue, cpu.vRegisters[instr.xIdx])
}
