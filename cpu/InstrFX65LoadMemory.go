package cpu

import (
	"fmt"
)

type InstrFX65LoadMemory struct {
	InstrParamsX
}

func (instr *InstrFX65LoadMemory) Execute(cpu *CpuInstance) {
	fmt.Print("FX65_Load_Memory, ")
	for i := range instr.x + 1 {
		cpu.vRegisters[i] = cpu.memory[cpu.indexRegister+uint16(i)]
		fmt.Printf("idx: %d, memloc: %d, value %d;\n", instr.x, cpu.indexRegister+uint16(i), cpu.memory[cpu.indexRegister+uint16(i)])
	}
}
