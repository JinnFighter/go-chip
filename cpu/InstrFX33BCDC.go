package cpu

import (
	"fmt"
)

type InstrFX33BCDC struct {
	InstrParamsX
}

func (instr *InstrFX33BCDC) Execute(cpu *CpuInstance) {
	var val = cpu.vRegisters[instr.x]
	var indexRegisterOffset = 0
	var startAddress = cpu.indexRegister
	var oldStartMemoryValue = cpu.memory[cpu.indexRegister]
	var divider = 100
	for divider > 0 {
		cpu.memory[cpu.indexRegister+uint16(indexRegisterOffset)] = val / uint8(divider)
		indexRegisterOffset += 1
		val %= uint8(divider)
		divider /= 10
	}

	var newStartMemoryValue = cpu.memory[cpu.indexRegister]
	fmt.Printf("FX33_Binary_Coded_Decimal_Conversion, idx: %d, startAddress %b, oldStartMemoryValue: %b, newStartMemoryValue: %b \n", instr.x, startAddress, oldStartMemoryValue, newStartMemoryValue)
}
