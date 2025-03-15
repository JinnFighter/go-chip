package cpu

type IInstruction interface {
	SetupValues(instructionBytes uint16)
	Execute(cpu *CpuInstance)
}

func CreateInstructions() map[uint16]IInstruction {
	var values = map[uint16]IInstruction{
		0x0000: &InstrWrapperLastByte{
			values: map[uint16]IInstruction{
				0x0000: &Instr00E0ClearScreen{},
			},
			defaultInstr: &Instr00EESubroutine{},
		},
		0x1000: &Instr1NNNJump{},
		0x2000: &Instr2NNNSubroutine{},
		0x3000: &Instr3XNNSkipConditionally{},
		0x4000: &Instr4XNNSkipConditionally{},
		0x5000: &Instr5XY0SkipConditionally{},
	}
	return values
}
