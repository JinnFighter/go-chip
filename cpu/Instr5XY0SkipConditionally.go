package cpu

import "fmt"

type Instr5XY0SkipConditionally struct {
	xIdx int
	yIdx int
}

func (instr *Instr5XY0SkipConditionally) SetupValues(instructionBytes uint16) {
	instr.xIdx = int((instructionBytes & 0x0F00) >> 8)
	instr.yIdx = int((instructionBytes & 0x00F0) >> 4)
}

func (instr *Instr5XY0SkipConditionally) Execute(cpu *CpuInstance) {
	var xValue = cpu.vRegisters[instr.xIdx]
	var yValue = cpu.vRegisters[instr.yIdx]
	if xValue == yValue {
		cpu.programCounter += 2
	}

	fmt.Printf("5XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", instr.xIdx, instr.yIdx, xValue, yValue)
}
