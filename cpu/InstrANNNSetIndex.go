package cpu

import "fmt"

type InstrANNNSetIndex struct {
	InstrParamsNNN
}

func (instr *InstrANNNSetIndex) Execute(cpu *CpuInstance) {
	var oldValue = cpu.indexRegister
	var newValue = instr.value
	cpu.indexRegister = newValue
	fmt.Printf("ANNN_SetIndex, oldValue %d, new value %d \n", oldValue, newValue)
}
