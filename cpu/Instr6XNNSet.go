package cpu

import "fmt"

type Instr6XNNSet struct {
	idx   int
	value uint8
}

func (instr *Instr6XNNSet) SetupValues(instructionBytes uint16) {
	instr.idx = int((instructionBytes & 0x0F00) >> 8)
	instr.value = uint8(instructionBytes & 0x00FF)
}

func (instr *Instr6XNNSet) Execute(cpu *CpuInstance) {
	var oldValue = cpu.vRegisters[instr.idx]
	var newValue = instr.value
	cpu.vRegisters[instr.idx] = newValue
	fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d \n", instr.idx, oldValue, newValue)
}
