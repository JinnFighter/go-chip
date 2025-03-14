package cpu

type InstrWrapperLastByte struct {
	lastByte     uint16
	values       map[uint16]IInstruction
	defaultInstr IInstruction
}

func (instr *InstrWrapperLastByte) SetupValues(instructionBytes uint16) {
	instr.lastByte = instructionBytes & 0x000F
	for _, value := range instr.values {
		value.SetupValues(instructionBytes)
	}

	if instr.defaultInstr != nil {
		instr.defaultInstr.SetupValues(instructionBytes)
	}
}

func (instr *InstrWrapperLastByte) Execute(cpu *CpuInstance) {
	var value = instr.values[instr.lastByte]
	if value != nil {
		value.Execute(cpu)
	} else if instr.defaultInstr != nil {
		instr.defaultInstr.Execute(cpu)
	}
}
