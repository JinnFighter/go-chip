package cpu

import (
	"fmt"
)

type InstrEX9ESkipIfKey struct {
	InstrParamsX
}

func (instr *InstrEX9ESkipIfKey) Execute(cpu *CpuInstance) {
	var keyVal = cpu.vRegisters[instr.x]
	if cpu.keyPressed[keyVal] {
		cpu.programCounter += 2
	}

	fmt.Printf("EX9E_Skip_If_Key, idx: %d, keyVal: %d, isPressed: %t\n", instr.x, keyVal, cpu.keyPressed[keyVal])
}
