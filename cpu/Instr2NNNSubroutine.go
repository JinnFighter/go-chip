package cpu

import "fmt"

type Instr2NNNSubroutine struct {
	value uint16
}

func (instr *Instr2NNNSubroutine) SetupValues(instructionBytes uint16) {
	instr.value = instructionBytes & 0x0FFF
}

func (instr *Instr2NNNSubroutine) Execute(cpu *CpuInstance) {
	cpu.addressStack.Push(cpu.programCounter)
	cpu.programCounter = instr.value
	fmt.Printf("2NNN_Subroutine, value: %b, stackCount: %d \n", instr.value, cpu.addressStack.Count())
}
