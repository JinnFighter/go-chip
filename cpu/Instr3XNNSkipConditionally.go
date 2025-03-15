package cpu

import "fmt"

type Instr3XNNSkipConditionally struct {
	idx   int
	value uint8
}

func (instr *Instr3XNNSkipConditionally) SetupValues(instructionBytes uint16) {
	instr.idx = int((instructionBytes & 0x0F00) >> 8)
	instr.value = uint8(instructionBytes & 0x00FF)
}

func (instr *Instr3XNNSkipConditionally) Execute(cpu *CpuInstance) {
	var registerValue = cpu.vRegisters[instr.idx]
	if registerValue == instr.value {
		cpu.programCounter += 2
	}

	fmt.Printf("3XNN_Skip_conditionally, idx: %d, registerValue: %b, value: %b \n", instr.idx, registerValue, instr.value)
}
