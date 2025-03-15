package cpu

import "fmt"

type Instr6XNNSet struct {
	InstrParamsXNN
}

func (instr *Instr6XNNSet) Execute(cpu *CpuInstance) {
	var oldValue = cpu.vRegisters[instr.x]
	var newValue = instr.value
	cpu.vRegisters[instr.x] = newValue
	fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d \n", instr.x, oldValue, newValue)
}
