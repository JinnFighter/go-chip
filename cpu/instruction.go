package cpu

type IInstruction interface {
	SetupValues(instructionBytes uint16)
	Execute(cpu *CpuInstance)
}

type InstrParams interface {
	SetupValues(instrBytes uint16)
}

type InstrParamsEmpty struct {
}

type InstrParamsXY struct {
	x int
	y int
}

type InstrParamsNNN struct {
	value uint16
}

type InstrParamsXNN struct {
	x     int
	value uint8
}

type InstrParamsXYN struct {
	x int
	y int
	n int
}

type InstrParamsX struct {
	x int
}

func (params *InstrParamsEmpty) SetupValues(instrBytes uint16) {
}

func (params *InstrParamsXY) SetupValues(instrBytes uint16) {
	params.x = int((instrBytes & 0x0F00) >> 8)
	params.y = int((instrBytes & 0x00F0) >> 4)
}

func (params *InstrParamsNNN) SetupValues(instrBytes uint16) {
	params.value = instrBytes & 0x0FFF
}

func (params *InstrParamsXNN) SetupValues(instrBytes uint16) {
	params.x = int((instrBytes & 0x0F00) >> 8)
	params.value = uint8(instrBytes & 0x00FF)
}

func (params *InstrParamsXYN) SetupValues(instrBytes uint16) {
	params.x = int((instrBytes & 0x0F00) >> 8)
	params.y = int((instrBytes & 0x00F0) >> 4)
	params.n = int(instrBytes & 0x000F)
}

func (params *InstrParamsX) SetupValues(instrBytes uint16) {
	params.x = int((instrBytes & 0x0F00) >> 8)
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
		0x6000: &Instr6XNNSet{},
		0x7000: &Instr7XNNAdd{},
		0x8000: &InstrWrapperLastByte{values: map[uint16]IInstruction{
			0x0000: &Instr8XY0Set{},
			0x0001: &Instr8XY1BinaryOr{},
			0x0002: &Instr8XY2BinaryAnd{},
			0x0003: &Instr8XY3BinaryXOr{},
			0x0004: &Instr8XY4Add{},
			0x0005: &Instr8XY5Subtract{},
			0x0006: &Instr8XY6ShiftRight{},
			0x0007: &Instr8XY7Subtract{},
			0x000E: &Instr8XYEShiftLeft{},
		}},
		0x9000: &Instr9XY0SkipConditionally{},
		0xA000: &InstrANNNSetIndex{},
		0xB000: &InstrBNNNJumpWithOffset{},
		0xC000: &InstrCXNNRandom{},
		0xD000: &InstrDXYNDisplay{},
		0xE000: &InstrWrapperByte3{values: map[uint16]IInstruction{
			0x9: &InstrEX9ESkipIfKey{},
			0xA: &InstrEXA1SkipIfNotKey{},
		}},
	}
	return values
}
