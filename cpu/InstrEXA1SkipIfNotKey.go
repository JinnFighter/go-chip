package cpu

import (
	"fmt"
)

type InstrEXA1SkipIfNotKey struct {
	InstrParamsX
}

func (instr *InstrEXA1SkipIfNotKey) Execute(cpu *CpuInstance) {
	var keyVal = cpu.vRegisters[instr.x]
	if !cpu.keyPressed[keyVal] {
		cpu.programCounter += 2
	}

	fmt.Printf("EXA1_Skip_If_Not_Key, idx: %d, keyVal: %d, isPressed: %t\n", instr.x, keyVal, cpu.keyPressed[keyVal])
}
