package cpu

import (
	"fmt"
	"go-chip/extensions"
	"math/rand"
)

const DisplayWidth = 64
const DisplayHeight = 32
const memorySize = 4096
const width = 64
const height = 32
const spriteWidth = 8
const registersSize = 16
const keysCount = 16
const programStartAddress = 0x200
const fontStartAddress = 0x50

var font = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type CpuInstance struct {
	memory         [memorySize]byte
	display        [width][height]bool
	vRegisters     [registersSize]uint8
	indexRegister  uint16
	programCounter uint16
	addressStack   extensions.Stack
	delayTimer     uint8
	soundTimer     uint8
	keyPressed     [keysCount]bool
	isRunning      bool
	instructions   map[uint16]IInstruction
}

func (cpu *CpuInstance) Init(romData []byte) {
	var instructions = CreateInstructions()
	cpu.instructions = instructions
	for i := range height {
		for j := range width {
			cpu.display[j][i] = false
		}
	}

	for index, fontByte := range font {
		cpu.memory[fontStartAddress+index] = fontByte
	}
	for index, b := range romData {
		cpu.memory[programStartAddress+index] = b
	}
	cpu.programCounter = programStartAddress
}

func (cpu *CpuInstance) Terminate() {
	cpu.stopLoopInner()

	cpu.programCounter = 0
	cpu.indexRegister = 0
	cpu.delayTimer = 0
	cpu.soundTimer = 0

	for i := range keysCount {
		cpu.keyPressed[i] = false
	}

	for i := range registersSize {
		cpu.vRegisters[i] = 0
	}

	for i := range memorySize {
		cpu.memory[i] = 0
	}
}

func (cpu *CpuInstance) StartLoop() {
	if cpu.isRunning {
		return
	}

	cpu.startLoopInner()
}

func (cpu *CpuInstance) StopLoop() {
	if !cpu.isRunning {
		return
	}

	cpu.stopLoopInner()
}

func (cpu *CpuInstance) startLoopInner() {
	cpu.isRunning = true
}

func (cpu *CpuInstance) stopLoopInner() {
	cpu.isRunning = false
}

func (cpu *CpuInstance) SetInputPressed(idx uint8, isPressed bool) {
	cpu.keyPressed[idx] = isPressed
}

func (cpu *CpuInstance) ExecuteTimersUpdate() {
	if cpu.delayTimer > 0 {
		cpu.delayTimer -= 1
	}

	if cpu.soundTimer > 0 {
		cpu.soundTimer -= 1
	}
}

func (cpu *CpuInstance) ExecuteLoopStep() {
	var nextInstruction = (uint16(cpu.memory[cpu.programCounter]) << 8) | uint16(cpu.memory[cpu.programCounter+1])
	cpu.programCounter += 2
	cpu.decodeInstruction(nextInstruction)
}

func (cpu *CpuInstance) GetDisplay(i int, j int) bool {
	return cpu.display[i][j]
}

func (cpu *CpuInstance) decodeInstruction(instructionBytes uint16) {
	var firstByte = instructionBytes & 0xF000
	var instruction = cpu.instructions[firstByte]
	if instruction != nil {
		instruction.SetupValues(instructionBytes)
		instruction.Execute(cpu)
		return
	}
	switch firstByte {
	case 0x6000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		cpu.Set_6XNN(idx, value)
	case 0x7000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		cpu.Add_7XNN(idx, value)
	case 0x8000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		var lastByte = instructionBytes & 0x000F
		switch lastByte {
		case 0x0000:
			cpu.Set_8XY0(xIdx, yIdx)
		case 0x0001:
			cpu.Binary_OR_8XY1(xIdx, yIdx)
		case 0x0002:
			cpu.Binary_AND_8XY2(xIdx, yIdx)
		case 0x0003:
			cpu.Binary_XOR_8XY3(xIdx, yIdx)
		case 0x0004:
			cpu.Add_8XY4(xIdx, yIdx)
		case 0x0005:
			cpu.Subtract_8XY5(xIdx, yIdx)
		case 0x0006:
			cpu.Shift_Right_8XY6(xIdx, yIdx)
		case 0x0007:
			cpu.Subtract_8XY7(xIdx, yIdx)
		case 0x000E:
			cpu.Shift_Left_8XYE(xIdx, yIdx)
		}
	case 0x9000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		cpu.Skip_conditionally_9XY0(xIdx, yIdx)
	case 0xA000:
		var value = instructionBytes & 0x0FFF
		cpu.SetIndex_ANNN(value)
	case 0xB000:
		var address = instructionBytes & 0x0FFF
		cpu.Jump_With_Offset_BNNN(address)
	case 0xC000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		cpu.Random_CXNN(idx, value)
	case 0xD000:
		var xRegister = int((instructionBytes & 0x0F00) >> 8)
		var yRegister = int((instructionBytes & 0x00F0) >> 4)
		var height = int(instructionBytes & 0x000F)
		cpu.Display_DXYN(xRegister, yRegister, height)
	case 0xE000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var checkedByte = (instructionBytes & 0x00F0) >> 4
		switch checkedByte {
		case 0x9:
			cpu.Skip_If_Key_EX9E(idx)
		case 0xA:
			cpu.Skip_If_Not_Key_EXA1(idx)
		}
	case 0xF000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var lastBytes = instructionBytes & 0x00FF
		switch lastBytes {
		case 0x000A:
			cpu.Get_Key_FX0A(idx)
		case 0x0007:
			cpu.Get_Value_Of_Delay_Timer_FX07(idx)
		case 0x0015:
			cpu.Set_Delay_Timer_FX15(idx)
		case 0x0018:
			cpu.Set_Sound_Timer_FX18(idx)
		case 0x001E:
			cpu.Add_To_Index_FX1E(idx)
		case 0x0029:
			cpu.Font_Character_FX29(idx)
		case 0x0033:
			cpu.Binary_Coded_Decimal_Conversion_FX33(idx)
		case 0x0055:
			cpu.Store_Memory_FX55(idx)
		case 0x0065:
			cpu.Load_Memory_FX65(idx)
		}
	default:
		fmt.Printf("Unknown Command\n")
	}
}

func (cpu *CpuInstance) Set_6XNN(idx int, value uint8) {
	var oldValue = cpu.vRegisters[idx]
	var newValue = value
	cpu.vRegisters[idx] = newValue
	fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d \n", idx, oldValue, newValue)
}

func (cpu *CpuInstance) Add_7XNN(idx int, value uint8) {
	var oldValue = cpu.vRegisters[idx]
	var newValue = oldValue + value

	cpu.vRegisters[idx] = newValue
	fmt.Printf("7XNN_Add, idx %d, oldValue %d, new value %d \n", idx, oldValue, newValue)
}

func (cpu *CpuInstance) SetIndex_ANNN(value uint16) {
	var oldValue = cpu.indexRegister
	var newValue = value
	cpu.indexRegister = newValue
	fmt.Printf("ANNN_SetIndex, oldValue %d, new value %d \n", oldValue, newValue)
}

func (cpu *CpuInstance) Display_DXYN(xRegister int, yRegister int, spriteHeight int) {
	var x = int(cpu.vRegisters[xRegister] & (width - 1))
	var y = int(cpu.vRegisters[yRegister] & (height - 1))
	cpu.vRegisters[15] = 0

	for col := 0; col < spriteWidth; col++ {
		for row := 0; row < int(spriteHeight); row++ {
			px := int(x) + col
			py := int(y) + row
			bit := (cpu.memory[cpu.indexRegister+uint16(row)] & (1 << uint(spriteWidth-1-col))) != 0
			if px < width && py < height && px >= 0 && py >= 0 {
				src := cpu.display[px][py]
				dst := bit != src // Да, оператор XOR с булевыми значениями не работает
				cpu.display[px][py] = dst
				if src && !dst {
					cpu.vRegisters[15] = 1
				}
			}
		}
	}

	fmt.Printf("DXYN_Display at xReg %d, yReg %d, height %d \n", xRegister, yRegister, spriteHeight)
}

func (cpu *CpuInstance) Skip_conditionally_9XY0(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	if xValue != yValue {
		cpu.programCounter += 2
	}

	fmt.Printf("9XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", xIdx, yIdx, xValue, yValue)
}

func (cpu *CpuInstance) Set_8XY0(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	cpu.vRegisters[xIdx] = yValue

	fmt.Printf("8XY0_Set, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b, new xRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, cpu.vRegisters[xIdx])
}

func (cpu *CpuInstance) Binary_OR_8XY1(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var newValue = xValue | yValue
	cpu.vRegisters[xIdx] = newValue

	fmt.Printf("8XY1_Binary_OR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, cpu.vRegisters[xIdx])
}

func (cpu *CpuInstance) Binary_AND_8XY2(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var newValue = xValue & yValue
	cpu.vRegisters[xIdx] = newValue
	fmt.Printf("8XY2_Binary_AND, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, cpu.vRegisters[xIdx])
}

func (cpu *CpuInstance) Binary_XOR_8XY3(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var newValue = xValue ^ yValue
	cpu.vRegisters[xIdx] = newValue

	fmt.Printf("8XY3_Binary_XOR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, cpu.vRegisters[xIdx])
}

func (cpu *CpuInstance) Add_8XY4(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var isCarryFlagSet = (int(xValue) + int(yValue)) > 255
	cpu.vRegisters[xIdx] = xValue + yValue

	if isCarryFlagSet {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}
}

func (cpu *CpuInstance) Subtract_8XY5(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var isCarryFlagSet = xValue >= yValue
	cpu.vRegisters[xIdx] = xValue - yValue

	if isCarryFlagSet {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}
}

func (cpu *CpuInstance) Subtract_8XY7(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var yValue = cpu.vRegisters[yIdx]
	var isCarryFlagSet = yValue >= xValue
	cpu.vRegisters[xIdx] = yValue - xValue

	if isCarryFlagSet {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}
}

func (cpu *CpuInstance) Shift_Right_8XY6(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var newValue = xValue >> 1
	var shiftedBit = xValue & 0x1
	cpu.vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}

	fmt.Printf("8XY6, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", xIdx, yIdx, xValue, newValue, shiftedBit, isCarryFlagSet)
}

func (cpu *CpuInstance) Shift_Left_8XYE(xIdx int, yIdx int) {
	var xValue = cpu.vRegisters[xIdx]
	var newValue = xValue << 1
	var shiftedBit = xValue >> 7
	cpu.vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}

	fmt.Printf("8XYE, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", xIdx, yIdx, xValue, newValue, shiftedBit, isCarryFlagSet)
}

func (cpu *CpuInstance) Jump_With_Offset_BNNN(address uint16) {
	cpu.programCounter = address + uint16(cpu.vRegisters[0])
}

func (cpu *CpuInstance) Random_CXNN(xIdx int, value uint8) {
	var rand = uint8(rand.Intn(256))
	var newValue = rand & value
	cpu.vRegisters[xIdx] = newValue
}

func (cpu *CpuInstance) Skip_If_Key_EX9E(idx int) {
	var keyVal = cpu.vRegisters[idx]
	if cpu.keyPressed[keyVal] {
		cpu.programCounter += 2
	}

	fmt.Printf("EX9E_Skip_If_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, keyVal, cpu.keyPressed[keyVal])
}

func (cpu *CpuInstance) Skip_If_Not_Key_EXA1(idx int) {
	var keyVal = cpu.vRegisters[idx]
	if !cpu.keyPressed[keyVal] {
		cpu.programCounter += 2
	}

	fmt.Printf("EXA1_Skip_If_Not_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, keyVal, cpu.keyPressed[keyVal])
}

func (cpu *CpuInstance) Get_Key_FX0A(idx int) {
	var isPressed = false
	var currentPressed = -1
	for i := range len(cpu.keyPressed) {
		if cpu.keyPressed[i] {
			isPressed = true
			currentPressed = i
			break
		}
	}

	if isPressed {
		cpu.vRegisters[idx] = uint8(currentPressed)
	} else {
		cpu.programCounter -= 2
	}

	fmt.Printf("FX0A_Get_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, currentPressed, isPressed)
}

func (cpu *CpuInstance) Binary_Coded_Decimal_Conversion_FX33(idx int) {
	var val = cpu.vRegisters[idx]
	var indexRegisterOffset = 0
	var divider = 100
	for divider > 0 {
		cpu.memory[cpu.indexRegister+uint16(indexRegisterOffset)] = val / uint8(divider)
		indexRegisterOffset += 1
		val %= uint8(divider)
		divider /= 10
	}
}

func (cpu *CpuInstance) Store_Memory_FX55(idx int) {
	fmt.Print("FX55_Store_Memory ")
	for i := range idx + 1 {
		cpu.memory[cpu.indexRegister+uint16(i)] = cpu.vRegisters[i]
		fmt.Printf("idx: %d, memloc: %d, value %d;", idx, cpu.indexRegister+uint16(i), cpu.vRegisters[i])
	}
	fmt.Println()
}

func (cpu *CpuInstance) Load_Memory_FX65(idx int) {
	fmt.Print("FX65_Load_Memory ")
	for i := range idx + 1 {
		cpu.vRegisters[i] = cpu.memory[cpu.indexRegister+uint16(i)]
		fmt.Printf("idx: %d, memloc: %d, value %d;", idx, cpu.indexRegister+uint16(i), cpu.memory[cpu.indexRegister+uint16(i)])
	}
	fmt.Println()
}

func (cpu *CpuInstance) Add_To_Index_FX1E(idx int) {
	var val = uint16(cpu.vRegisters[idx])
	var isOverFlow = val > cpu.indexRegister
	cpu.indexRegister += val
	if isOverFlow {
		cpu.vRegisters[15] = 1
	} else {
		cpu.vRegisters[15] = 0
	}
}

func (cpu *CpuInstance) Get_Value_Of_Delay_Timer_FX07(idx int) {
	var oldValue = cpu.vRegisters[idx]
	cpu.vRegisters[idx] = cpu.delayTimer

	fmt.Printf("FX07_Get_Value_Of_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, cpu.delayTimer, cpu.vRegisters[idx])
}

func (cpu *CpuInstance) Set_Delay_Timer_FX15(idx int) {
	var oldValue = cpu.delayTimer
	var newValue = cpu.vRegisters[idx]
	cpu.delayTimer = newValue
	fmt.Printf("FX15_Set_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, cpu.delayTimer, newValue)
}

func (cpu *CpuInstance) Set_Sound_Timer_FX18(idx int) {
	var oldValue = cpu.soundTimer
	var newValue = cpu.vRegisters[idx]
	cpu.soundTimer = newValue
	fmt.Printf("FX18_Set_Sound_Timer, old value: %d, currentSoundTimer: %d, newValue: %d\n", oldValue, cpu.soundTimer, newValue)
}

func (cpu *CpuInstance) Font_Character_FX29(idx int) {
	var oldIndexRegister = cpu.indexRegister
	var newAddress = uint16(cpu.vRegisters[idx] & 0x00F)
	cpu.indexRegister = newAddress
	fmt.Printf("FX29_Font_Character, idx: %d, old address: %d, new address: %d\n", idx, oldIndexRegister, newAddress)
}
