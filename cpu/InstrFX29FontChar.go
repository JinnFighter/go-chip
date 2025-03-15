package cpu

import (
	"fmt"
)

type InstrFX29FontChar struct {
	InstrParamsX
}

func (instr *InstrFX29FontChar) Execute(cpu *CpuInstance) {
	var oldIndexRegister = cpu.indexRegister
	var newAddress = uint16(cpu.vRegisters[instr.x] & 0x00F)
	cpu.indexRegister = newAddress

	fmt.Printf("FX29_Font_Character, idx: %d, old address: %d, new address: %d\n", instr.x, oldIndexRegister, newAddress)
}
