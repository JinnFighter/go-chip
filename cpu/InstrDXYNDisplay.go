package cpu

import (
	"fmt"
)

type InstrDXYNDisplay struct {
	InstrParamsXYN
}

func (instr *InstrDXYNDisplay) Execute(cpu *CpuInstance) {
	var x = int(cpu.vRegisters[instr.x] & (DisplayWidth - 1))
	var y = int(cpu.vRegisters[instr.y] & (DisplayHeight - 1))
	cpu.vRegisters[vfIndex] = 0

	for col := range spriteWidth {
		for row := range instr.n {
			px := int(x) + col
			py := int(y) + row
			bit := (cpu.memory[cpu.indexRegister+uint16(row)] & (1 << uint(spriteWidth-1-col))) != 0
			if px < DisplayWidth && py < DisplayHeight && px >= 0 && py >= 0 {
				src := cpu.display[px][py]
				dst := bit != src
				cpu.display[px][py] = dst
				if src && !dst {
					cpu.vRegisters[vfIndex] = 1
				}
			}
		}
	}

	cpu.IsRedraw = true
	fmt.Printf("DXYN_Display at xReg %d, yReg %d, height %d \n", instr.x, instr.y, instr.n)
}
