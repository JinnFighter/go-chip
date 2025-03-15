package cpu

import "fmt"

type Instr8XY3BinaryXOr struct {
	InstrParamsXY
}

func (instr *Instr8XY3BinaryXOr) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	var newValue = xValue ^ yValue
	cpu.vRegisters[instr.xIdx] = newValue

	fmt.Printf("8XY3_Binary_XOR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", instr.xIdx, instr.yIdx, xValue, yValue, newValue, cpu.vRegisters[instr.xIdx])
}
