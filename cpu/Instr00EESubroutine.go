package cpu

import (
	"fmt"
)

type Instr00EESubroutine struct {
}

func (instr *Instr00EESubroutine) SetupValues(instructionBytes uint16) {

}

func (instr *Instr00EESubroutine) Execute(cpu *CpuInstance) {
	var address = cpu.addressStack.Pop()
	cpu.programCounter = address

	fmt.Printf("00EE_Subroutine, stackCount: %d \n", cpu.addressStack.Count())
}
