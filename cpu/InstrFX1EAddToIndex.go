package cpu

import (
	"fmt"
)

type InstrFX1EAddToIndex struct {
	InstrParamsX
}

func (instr *InstrFX1EAddToIndex) Execute(cpu *CpuInstance) {
	var val = uint16(cpu.vRegisters[instr.x])
	var oldValue = cpu.indexRegister
	var isOverFlow = val > cpu.indexRegister
	var newValue = cpu.indexRegister + val
	cpu.indexRegister = newValue
	if isOverFlow {
		cpu.vRegisters[vfIndex] = 1
	} else {
		cpu.vRegisters[vfIndex] = 0
	}
	fmt.Printf("FX1E_Add_To_Index, old value: %d, newValue: %d, isOverflow: %t \n", oldValue, newValue, isOverFlow)
}
