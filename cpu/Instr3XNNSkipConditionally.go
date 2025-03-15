package cpu

import "fmt"

type Instr3XNNSkipConditionally struct {
	InstrParamsXNN
}

func (instr *Instr3XNNSkipConditionally) Execute(cpu *CpuInstance) {
	var registerValue = cpu.vRegisters[instr.x]
	if registerValue == instr.value {
		cpu.programCounter += 2
	}

	fmt.Printf("3XNN_Skip_conditionally, idx: %d, registerValue: %b, value: %b \n", instr.x, registerValue, instr.value)
}
