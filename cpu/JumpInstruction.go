package cpu

import (
	"fmt"
)

type JumpInstruction struct {
	jumpAddress uint16
}

func (instruction *JumpInstruction) SetupValues(instructionBytes uint16) {
	instruction.jumpAddress = instructionBytes & 0x0FFF
}

func (instruction *JumpInstruction) Execute(cpu *CpuInstance) {
	var oldAddress = cpu.programCounter
	cpu.programCounter = instruction.jumpAddress
	fmt.Printf("1NNN_Jump, old address: %d, new address: %d \n", oldAddress, instruction.jumpAddress)
}
