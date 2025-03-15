package cpu

import (
	"fmt"
)

type InstrFX55StoreMemory struct {
	InstrParamsX
}

func (instr *InstrFX55StoreMemory) Execute(cpu *CpuInstance) {
	fmt.Print("FX55_Store_Memory, ")
	for i := range instr.x + 1 {
		cpu.memory[cpu.indexRegister+uint16(i)] = cpu.vRegisters[i]
		fmt.Printf("idx: %d, memloc: %d, value %d; \n", instr.x, cpu.indexRegister+uint16(i), cpu.vRegisters[i])
	}
}
