package cpu

import (
	"fmt"
)

type Instr00E0ClearScreen struct {
	InstrParamsEmpty
}

func (instr *Instr00E0ClearScreen) Execute(cpu *CpuInstance) {
	for i := range DisplayWidth {
		for j := range DisplayHeight {
			cpu.display[i][j] = false
		}
	}
	fmt.Printf("00E0_ClearScreen \n")
}
