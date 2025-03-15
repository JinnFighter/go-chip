package cpu

import (
	"fmt"
)

type Instr00EESubroutine struct {
	InstrParamsEmpty
}

func (instr *Instr00EESubroutine) Execute(cpu *CpuInstance) {
	var address = cpu.addressStack.Pop()
	cpu.programCounter = address

	fmt.Printf("00EE_Subroutine, stackCount: %d \n", cpu.addressStack.Count())
}
