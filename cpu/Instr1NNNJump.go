package cpu

import (
	"fmt"
)

type Instr1NNNJump struct {
	InstrParamsNNN
}

func (instruction *Instr1NNNJump) Execute(cpu *CpuInstance) {
	var oldAddress = cpu.programCounter
	cpu.programCounter = instruction.value
	fmt.Printf("1NNN_Jump, old address: %d, new address: %d \n", oldAddress, instruction.value)
}
