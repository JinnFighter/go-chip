package cpu

import (
	"fmt"
	"go-chip/extensions"
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
const vfIndex = 15

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
