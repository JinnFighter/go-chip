package cpu

import "fmt"

type Instr5XY0SkipConditionally struct {
	InstrParamsXY
}

func (instr *Instr5XY0SkipConditionally) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.x]
	var yValue = cpu.vRegisters[instr.y]
	if xValue == yValue {
		cpu.programCounter += 2
	}

	fmt.Printf("5XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", instr.x, instr.y, xValue, yValue)
}
