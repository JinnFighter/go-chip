package cpu

import "fmt"

type Instr7XNNAdd struct {
	idx   int
	value uint8
}

func (instr *Instr7XNNAdd) SetupValues(instructionBytes uint16) {
	instr.idx = int((instructionBytes & 0x0F00) >> 8)
	instr.value = uint8(instructionBytes & 0x00FF)
}

func (instr *Instr7XNNAdd) Execute(cpu *CpuInstance) {
	var oldValue = cpu.vRegisters[instr.idx]
	var newValue = oldValue + instr.value

	cpu.vRegisters[instr.idx] = newValue
	fmt.Printf("7XNN_Add, idx %d, oldValue %d, new value %d \n", instr.idx, oldValue, newValue)
}
