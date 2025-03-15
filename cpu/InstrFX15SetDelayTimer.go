package cpu

import (
	"fmt"
)

type InstrFX15SetDelayTimer struct {
	InstrParamsX
}

func (instr *InstrFX15SetDelayTimer) Execute(cpu *CpuInstance) {
	var oldValue = cpu.delayTimer
	var newValue = cpu.vRegisters[instr.x]
	cpu.delayTimer = newValue
	fmt.Printf("FX15_Set_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, cpu.delayTimer, newValue)
}
