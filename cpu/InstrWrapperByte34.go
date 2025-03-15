package cpu

type InstrWrapperByte34 struct {
	nn           uint16
	values       map[uint16]IInstruction
	defaultInstr IInstruction
}

func (instr *InstrWrapperByte34) SetupValues(instructionBytes uint16) {
	instr.nn = instructionBytes & 0x00FF
	for _, value := range instr.values {
		value.SetupValues(instructionBytes)
	}

	if instr.defaultInstr != nil {
		instr.defaultInstr.SetupValues(instructionBytes)
	}
}

func (instr *InstrWrapperByte34) Execute(cpu *CpuInstance) {
	var value = instr.values[instr.nn]
	if value != nil {
		value.Execute(cpu)
	} else if instr.defaultInstr != nil {
		instr.defaultInstr.Execute(cpu)
	}
}
