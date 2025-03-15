package cpu

import (
	"fmt"
)

type InstrFX07GetDelayTimer struct {
	InstrParamsX
}

func (instr *InstrFX07GetDelayTimer) Execute(cpu *CpuInstance) {
	var oldValue = cpu.vRegisters[instr.x]
	cpu.vRegisters[instr.x] = cpu.delayTimer

	fmt.Printf("FX07_Get_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, cpu.delayTimer, cpu.vRegisters[instr.x])
}
