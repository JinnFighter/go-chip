package cpu

import (
	"fmt"
)

type Instr00E0ClearScreen struct {
}

func (instr *Instr00E0ClearScreen) SetupValues(instructionBytes uint16) {

}

func (instr *Instr00E0ClearScreen) Execute(cpu *CpuInstance) {
	for i := range DisplayWidth {
		for j := range DisplayHeight {
			cpu.display[i][j] = false
		}
	}
	fmt.Printf("00E0_ClearScreen \n")
}
