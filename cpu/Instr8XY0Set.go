package cpu

import "fmt"

type Instr8XY0Set struct {
	InstrParamsXY
}

func (instr *Instr8XY0Set) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	cpu.vRegisters[instr.x] = yValue

	fmt.Printf("8XY0_Set, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b, new xRegisterValue: %b \n", instr.x, instr.y, xValue, yValue, cpu.vRegisters[instr.x])
}
