package cpu

import (
	"fmt"
	"math/rand"
)

type InstrCXNNRandom struct {
	InstrParamsXNN
}

func (instr *InstrCXNNRandom) Execute(cpu *CpuInstance) {
	var rand = uint8(rand.Intn(256))
	var oldValue = cpu.vRegisters[instr.x]
	var newValue = rand & instr.value
	cpu.vRegisters[instr.x] = newValue

	fmt.Printf("CXNN_Random, oldValue %d, new value %d \n", oldValue, newValue)
}
