package cpu

import "fmt"

type Instr8XY0Set struct {
	InstrParamsXY
}

func (instr *Instr8XY0Set) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	cpu.vRegisters[instr.xIdx] = yValue

	fmt.Printf("8XY0_Set, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b, new xRegisterValue: %b \n", instr.xIdx, instr.yIdx, xValue, yValue, cpu.vRegisters[instr.xIdx])
}
