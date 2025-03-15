package cpu

import "fmt"

type Instr9XY0SkipConditionally struct {
	InstrParamsXY
}

func (instr *Instr9XY0SkipConditionally) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	if xValue != yValue {
		cpu.programCounter += 2
	}

	fmt.Printf("9XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", instr.x, instr.y, xValue, yValue)
}
