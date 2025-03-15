package cpu

import (
	"fmt"
)

type InstrFX18SetSoundTimer struct {
	InstrParamsX
}

func (instr *InstrFX18SetSoundTimer) Execute(cpu *CpuInstance) {
	var oldValue = cpu.soundTimer
	var newValue = cpu.vRegisters[instr.x]
	cpu.soundTimer = newValue
	fmt.Printf("FX18_Set_Sound_Timer, old value: %d, currentSoundTimer: %d, newValue: %d\n", oldValue, cpu.soundTimer, newValue)
}
