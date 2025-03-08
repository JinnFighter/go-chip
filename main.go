package main

import (
	"fmt"
	"go-chip/extensions"
	"math/rand"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const timerDecreaseSpeed = 60.0
const instructionExecutionSpeed = 700.0

const width = 64
const height = 32
const spriteWidth = 8

var memory [4096]byte
var display [width][height]bool
var displayRects [width][height]sdl.Rect
var blackColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var whiteColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
var vRegisters [16]uint8
var indexRegister uint16
var programCounter uint16
var addressStack extensions.Stack
var delayTimer uint8
var soundTimer uint8
var keyPressed [16]bool
var window *sdl.Window
var windowSurface *sdl.Surface
var isRunning bool
var keysBytes = map[sdl.Keycode]uint8{
	sdl.K_1: 0x1,
	sdl.K_2: 0x2,
	sdl.K_3: 0x3,
	sdl.K_4: 0xC,
	sdl.K_q: 0x4,
	sdl.K_w: 0x5,
	sdl.K_e: 0x6,
	sdl.K_r: 0xD,
	sdl.K_a: 0x7,
	sdl.K_s: 0x8,
	sdl.K_d: 0x9,
	sdl.K_f: 0xE,
	sdl.K_z: 0xA,
	sdl.K_x: 0x0,
	sdl.K_c: 0xB,
	sdl.K_v: 0xF,
}

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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	var args = os.Args[1:]
	data, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	for index, fontByte := range font {
		memory[0x050+index] = fontByte
	}

	for index, b := range data {
		memory[index+0x200] = b
	}
	programCounter = 0x200

	window, _ = sdl.CreateWindow("Go-Chip by JinnFighter", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 640, 320, sdl.WINDOW_SHOWN)
	windowSurface, _ = window.GetSurface()
	windowSurface.FillRect(nil, 0)

	for i := range height {
		for j := range width {
			displayRects[j][i] = sdl.Rect{X: int32(j * 10), Y: int32(i * 10), W: 7, H: 7}
		}
	}

	startLoop()

	window.Destroy()
}

func startLoop() {
	if isRunning {
		return
	}

	isRunning = true

	loop()
}

func stopLoop() {
	if !isRunning {
		return
	}

	isRunning = false
}

func loop() {
	fmt.Println("Enter loop")

	for {
		inputLoop()
		if !isRunning {
			return
		}
		var nextInstruction = (uint16(memory[programCounter]) << 8) | uint16(memory[programCounter+1])
		programCounter += 2
		decodeInstruction(nextInstruction)

		for i := range height {
			for j := range width {
				if display[j][i] {
					windowSurface.FillRect(&displayRects[j][i], whiteColor.Uint32())
				} else {
					windowSurface.FillRect(&displayRects[j][i], blackColor.Uint32())
				}
			}
		}

		window.UpdateSurface()
		sdl.Delay(1000 / 60)
	}
}

func inputLoop() {
	if !isRunning {
		return
	}
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch et := event.(type) {
		case *sdl.QuitEvent:
			stopLoop()
			return
		case *sdl.KeyboardEvent:
			var keyByte, isPresent = keysBytes[et.Keysym.Sym]
			if isPresent {
				if et.Type == sdl.KEYDOWN {
					keyPressed[keyByte] = true
				} else if et.Type == sdl.KEYUP {
					keyPressed[keyByte] = false
				}
			}
		}
	}
}

func decodeInstruction(instructionBytes uint16) {
	var firstByte = instructionBytes & 0xF000
	switch firstByte {
	case 0x0000:
		var lastByte = instructionBytes & 0x000F
		if lastByte == 0 {
			ClearScreen_00E0()
		} else {
			Subroutine_00EE()
		}
	case 0x1000:
		var jumpAddress = instructionBytes & 0x0FFF
		Jump_1NNN(jumpAddress)
	case 0x2000:
		var address = instructionBytes & 0x0FFF
		Subroutine_2NNN(address)
	case 0x3000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Skip_conditionally_3XNN(idx, value)
	case 0x4000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Skip_conditionally_4XNN(idx, value)
	case 0x5000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		Skip_conditionally_5XY0(xIdx, yIdx)
	case 0x6000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Set_6XNN(idx, value)
	case 0x7000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Add_7XNN(idx, value)
	case 0x8000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		var lastByte = instructionBytes & 0x000F
		switch lastByte {
		case 0x0000:
			Set_8XY0(xIdx, yIdx)
		case 0x0001:
			Binary_OR_8XY1(xIdx, yIdx)
		case 0x0002:
			Binary_AND_8XY2(xIdx, yIdx)
		case 0x0003:
			Binary_XOR_8XY3(xIdx, yIdx)
		case 0x0004:
			Add_8XY4(xIdx, yIdx)
		case 0x0005:
			Subtract_8XY5(xIdx, yIdx)
		case 0x0006:
			Shift_Right_8XY6(xIdx, yIdx)
		case 0x0007:
			Subtract_8XY7(xIdx, yIdx)
		case 0x000E:
			Shift_Left_8XYE(xIdx, yIdx)
		}
	case 0x9000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		Skip_conditionally_9XY0(xIdx, yIdx)
	case 0xA000:
		var value = instructionBytes & 0x0FFF
		SetIndex_ANNN(value)
	case 0xB000:
		var address = instructionBytes & 0x0FFF
		Jump_With_Offset_BNNN(address)
	case 0xC000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Random_CXNN(idx, value)
	case 0xD000:
		var xRegister = int((instructionBytes & 0x0F00) >> 8)
		var yRegister = int((instructionBytes & 0x00F0) >> 4)
		var height = int(instructionBytes & 0x000F)
		Display_DXYN(xRegister, yRegister, height)
	case 0xE000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var checkedByte = (instructionBytes & 0x00F0) >> 4
		switch checkedByte {
		case 0x9:
			Skip_If_Key_EX9E(idx)
		case 0xA:
			Skip_If_Not_Key_EXA1(idx)
		}
	case 0xF000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var lastBytes = instructionBytes & 0x00FF
		switch lastBytes {
		case 0x000A:
			Get_Key_FX0A(idx)
		case 0x0007:
			Get_Value_Of_Delay_Timer_FX07(idx)
		case 0x0015:
			Set_Delay_Timer_FX15(idx)
		case 0x0018:
			Set_Sound_Timer_FX18(idx)
		case 0x001E:
			Add_To_Index_FX1E(idx)
		case 0x0029:
			Font_Character_FX29(idx)
		case 0x0033:
			Binary_Coded_Decimal_Conversion_FX33(idx)
		case 0x0055:
			Store_Memory_FX55(idx)
		case 0x0065:
			Load_Memory_FX65(idx)
		}
	default:
		fmt.Printf("Unknown Command\n")
	}
}

func ClearScreen_00E0() {
	for i := range width {
		for j := range height {
			display[i][j] = false
		}
	}
	fmt.Printf("00E0_ClearScreen\n")
}

func Jump_1NNN(jumpAddress uint16) {
	var oldAddress = programCounter
	programCounter = jumpAddress
	fmt.Printf("1NNN_Jump, old address: %d, new address: %d \n", oldAddress, jumpAddress)
}

func Set_6XNN(idx int, value uint8) {
	var oldValue = vRegisters[idx]
	var newValue = value
	vRegisters[idx] = newValue
	fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d \n", idx, oldValue, newValue)
}

func Add_7XNN(idx int, value uint8) {
	var oldValue = vRegisters[idx]
	var newValue = oldValue + value

	vRegisters[idx] = newValue
	fmt.Printf("7XNN_Add, idx %d, oldValue %d, new value %d \n", idx, oldValue, newValue)
}

func SetIndex_ANNN(value uint16) {
	var oldValue = indexRegister
	var newValue = value
	indexRegister = newValue
	fmt.Printf("ANNN_SetIndex, oldValue %d, new value %d \n", oldValue, newValue)
}

func Display_DXYN(xRegister int, yRegister int, spriteHeight int) {
	var x = int(vRegisters[xRegister] & (width - 1))
	var y = int(vRegisters[yRegister] & (height - 1))
	vRegisters[15] = 0

	for col := 0; col < 8; col++ {
		for row := 0; row < int(spriteHeight); row++ {
			px := int(x) + col
			py := int(y) + row
			bit := (memory[indexRegister+uint16(row)] & (1 << uint(8-1-col))) != 0
			if px < 64 && py < 32 && px >= 0 && py >= 0 {
				src := display[px][py]
				dst := bit != src // Да, оператор XOR с булевыми значениями не работает
				display[px][py] = dst
				if src && !dst {
					vRegisters[0xf] = 1
				}
			}
		}
	}

	fmt.Printf("DXYN_Display at xReg %d, yReg %d, height %d \n", xRegister, yRegister, spriteHeight)
}

func Subroutine_2NNN(value uint16) {
	addressStack.Push(programCounter)
	programCounter = value
	fmt.Printf("2NNN_Subroutine, value: %b, stackCount: %d \n", value, addressStack.Count())
}

func Subroutine_00EE() {
	var address = addressStack.Pop()
	programCounter = address

	fmt.Printf("00EE_Subroutine, stackCount: %d \n", addressStack.Count())
}

func Skip_conditionally_3XNN(idx int, value uint8) {
	var registerValue = vRegisters[idx]
	if registerValue == value {
		programCounter += 2
	}

	fmt.Printf("3XNN_Skip_conditionally, idx: %d, registerValue: %b, value: %b \n", idx, registerValue, value)
}

func Skip_conditionally_4XNN(idx int, value uint8) {
	var registerValue = vRegisters[idx]
	if registerValue != value {
		programCounter += 2
	}

	fmt.Printf("4XNN_Skip_conditionally, idx: %d, registerValue: %b, value: %b \n", idx, registerValue, value)
}

func Skip_conditionally_5XY0(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	if xValue == yValue {
		programCounter += 2
	}

	fmt.Printf("5XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", xIdx, yIdx, xValue, yValue)
}

func Skip_conditionally_9XY0(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	if xValue != yValue {
		programCounter += 2
	}

	fmt.Printf("9XY0_Skip_conditionally, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b \n", xIdx, yIdx, xValue, yValue)
}

func Set_8XY0(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	vRegisters[xIdx] = yValue

	fmt.Printf("8XY0_Set, xIdx: %d, yIdx: %d, xRegisterValue: %b, yRegisterValue: %b, new xRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, vRegisters[xIdx])
}

func Binary_OR_8XY1(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var newValue = xValue | yValue
	vRegisters[xIdx] = newValue

	fmt.Printf("8XY1_Binary_OR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, vRegisters[xIdx])
}

func Binary_AND_8XY2(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var newValue = xValue & yValue
	vRegisters[xIdx] = newValue
	fmt.Printf("8XY2_Binary_AND, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, vRegisters[xIdx])
}

func Binary_XOR_8XY3(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var newValue = xValue ^ yValue
	vRegisters[xIdx] = newValue

	fmt.Printf("8XY3_Binary_XOR, xIdx: %d, yIdx: %d, xValue: %b, yValue: %b, newValue: %b, newXRegisterValue: %b \n", xIdx, yIdx, xValue, yValue, newValue, vRegisters[xIdx])
}

func Add_8XY4(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var isCarryFlagSet = (int(xValue) + int(yValue)) > 255
	vRegisters[xIdx] = xValue + yValue

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}

func Subtract_8XY5(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var isCarryFlagSet = xValue >= yValue
	vRegisters[xIdx] = xValue - yValue

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}

func Subtract_8XY7(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	var isCarryFlagSet = yValue >= xValue
	vRegisters[xIdx] = yValue - xValue

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}

func Shift_Right_8XY6(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var newValue = xValue >> 1
	var shiftedBit = xValue & 0x1
	vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}

	fmt.Printf("8XY6, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", xIdx, yIdx, xValue, newValue, shiftedBit, isCarryFlagSet)
}

func Shift_Left_8XYE(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var newValue = xValue << 1
	var shiftedBit = xValue >> 7
	vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}

	fmt.Printf("8XYE, X = %d, Y = %d, xValue: %b, newValue: %b, shiftedBit: %b, isCarryFlagSet: %t, \n", xIdx, yIdx, xValue, newValue, shiftedBit, isCarryFlagSet)
}

func Jump_With_Offset_BNNN(address uint16) {
	programCounter = address + uint16(vRegisters[0])
}

func Random_CXNN(xIdx int, value uint8) {
	var rand = uint8(rand.Intn(256))
	var newValue = rand & value
	vRegisters[xIdx] = newValue
}

func Skip_If_Key_EX9E(idx int) {
	var keyVal = vRegisters[idx]
	if keyPressed[keyVal] {
		programCounter += 2
	}

	fmt.Printf("EX9E_Skip_If_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, keyVal, keyPressed[keyVal])
}

func Skip_If_Not_Key_EXA1(idx int) {
	var keyVal = vRegisters[idx]
	if !keyPressed[keyVal] {
		programCounter += 2
	}

	fmt.Printf("EXA1_Skip_If_Not_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, keyVal, keyPressed[keyVal])
}

func Get_Key_FX0A(idx int) {
	var isPressed = false
	var currentPressed = -1
	for i := range len(keyPressed) {
		if keyPressed[i] {
			isPressed = true
			currentPressed = i
			break
		}
	}

	if isPressed {
		vRegisters[idx] = uint8(currentPressed)
	} else {
		programCounter -= 2
	}

	fmt.Printf("FX0A_Get_Key, idx: %d, keyVal: %d, isPressed: %t\n", idx, currentPressed, isPressed)
}

func Binary_Coded_Decimal_Conversion_FX33(idx int) {
	var val = vRegisters[idx]
	var indexRegisterOffset = 0
	var divider = 100
	for divider > 0 {
		memory[indexRegister+uint16(indexRegisterOffset)] = val / uint8(divider)
		indexRegisterOffset += 1
		val %= uint8(divider)
		divider /= 10
	}
}

func Store_Memory_FX55(idx int) {
	fmt.Print("FX55_Store_Memory ")
	for i := range idx + 1 {
		memory[indexRegister+uint16(i)] = vRegisters[i]
		fmt.Printf("idx: %d, memloc: %d, value %d;", idx, indexRegister+uint16(i), vRegisters[i])
	}
	fmt.Println()
}

func Load_Memory_FX65(idx int) {
	fmt.Print("FX65_Load_Memory ")
	for i := range idx + 1 {
		vRegisters[i] = memory[indexRegister+uint16(i)]
		fmt.Printf("idx: %d, memloc: %d, value %d;", idx, indexRegister+uint16(i), memory[indexRegister+uint16(i)])
	}
	fmt.Println()
}

func Add_To_Index_FX1E(idx int) {
	var val = uint16(vRegisters[idx])
	var isOverFlow = val > indexRegister
	indexRegister += val
	if isOverFlow {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}

func Get_Value_Of_Delay_Timer_FX07(idx int) {
	var oldValue = vRegisters[idx]
	vRegisters[idx] = delayTimer

	fmt.Printf("FX07_Get_Value_Of_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, delayTimer, vRegisters[idx])
}

func Set_Delay_Timer_FX15(idx int) {
	var oldValue = delayTimer
	var newValue = vRegisters[idx]
	delayTimer = newValue
	fmt.Printf("FX15_Set_Delay_Timer, old value: %d, currentDelayTimer: %d, newValue: %d\n", oldValue, delayTimer, newValue)
}

func Set_Sound_Timer_FX18(idx int) {
	var oldValue = soundTimer
	var newValue = vRegisters[idx]
	soundTimer = newValue
	fmt.Printf("FX18_Set_Sound_Timer, old value: %d, currentSoundTimer: %d, newValue: %d\n", oldValue, soundTimer, newValue)
}

func Font_Character_FX29(idx int) {
	var oldIndexRegister = indexRegister
	var newAddress = uint16(vRegisters[idx] & 0x00F)
	indexRegister = newAddress
	fmt.Printf("FX29_Font_Character, idx: %d, old address: %d, new address: %d\n", idx, oldIndexRegister, newAddress)
}
