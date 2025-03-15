package cpu

type InstrWrapperByte3 struct {
	n            uint16
	values       map[uint16]IInstruction
	defaultInstr IInstruction
}

func (instr *InstrWrapperByte3) SetupValues(instructionBytes uint16) {
	instr.n = (instructionBytes & 0x00F0) >> 4
	for _, value := range instr.values {
		value.SetupValues(instructionBytes)
	}

	if instr.defaultInstr != nil {
		instr.defaultInstr.SetupValues(instructionBytes)
	}
}

func (instr *InstrWrapperByte3) Execute(cpu *CpuInstance) {
	var value = instr.values[instr.n]
	if value != nil {
		value.Execute(cpu)
	} else if instr.defaultInstr != nil {
		instr.defaultInstr.Execute(cpu)
	}
}
