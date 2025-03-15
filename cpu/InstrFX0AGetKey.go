package cpu

import (
	"fmt"
)

type InstrFX0AGetKey struct {
	InstrParamsX
}

func (instr *InstrFX0AGetKey) Execute(cpu *CpuInstance) {
	var isPressed = false
	var currentPressed = -1
	for i := range len(cpu.keyPressed) {
		if cpu.keyPressed[i] {
			isPressed = true
			currentPressed = i
			break
		}
	}

	if isPressed {
		cpu.vRegisters[instr.x] = uint8(currentPressed)
	} else {
		cpu.programCounter -= 2
	}

	fmt.Printf("FX0A_Get_Key, idx: %d, keyVal: %d, isPressed: %t\n", instr.x, currentPressed, isPressed)
}
