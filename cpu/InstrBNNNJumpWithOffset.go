package cpu

import "fmt"

type InstrBNNNJumpWithOffset struct {
	InstrParamsNNN
}

func (instr *InstrBNNNJumpWithOffset) Execute(cpu *CpuInstance) {
	var oldValue = cpu.programCounter
	var newValue = instr.value + uint16(cpu.vRegisters[0])
	cpu.programCounter = newValue
	fmt.Printf("BNNN_JumpWithOffset, oldValue %d, new value %d \n", oldValue, newValue)
}
