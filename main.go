package main

import (
	"fmt"
	"go-chip/extensions"
	"os"
	"os/exec"
	"time"
)

const timerDecreaseSpeed = 60.0
const instructionExecutionSpeed = 700.0

const width = 64
const height = 32
const spriteWidth = 8

var memory [4096]byte
var display [width][height]bool
var vRegisters [16]uint8
var indexRegister uint16
var programCounter uint16
var addressStack extensions.Stack
var delayTimer uint8
var soundTimer uint8
var isRunning bool
var ticker *time.Ticker
var tickerChannel chan bool

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

	startLoop()
	for isRunning {

	}
}

func startLoop() {
	if isRunning {
		return
	}

	isRunning = true
	var execSpeed = 1 / instructionExecutionSpeed * float64(time.Second)
	ticker = time.NewTicker(time.Duration(execSpeed))
	tickerChannel = make(chan bool)

	go loop()
}

func stopLoop() {
	if !isRunning {
		return
	}

	isRunning = false
	ticker.Stop()
	tickerChannel <- true
}

func loop() {
	fmt.Println("Enter loop")
	for {
		select {
		case <-tickerChannel:
			return
		case <-ticker.C:

			var nextInstruction = (uint16(memory[programCounter]) << 8) | uint16(memory[programCounter+1])
			programCounter += 2
			decodeInstruction(nextInstruction)

			for i := range height {
				var str = ""
				for j := range width {
					if display[j][i] {
						str += "*"
					} else {
						str += " "
					}
				}
				fmt.Println(str)
			}

			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
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
			Shift_8XY6(xIdx, yIdx)
		case 0x0007:
			Subtract_8XY7(xIdx, yIdx)
		case 0x000E:
			Shift_8XYE(xIdx, yIdx)
		}
	case 0x9000:
		var xIdx = int((instructionBytes & 0x0F00) >> 8)
		var yIdx = int((instructionBytes & 0x00F0) >> 4)
		Skip_conditionally_9XY0(xIdx, yIdx)
	case 0xA000:
		var value = instructionBytes & 0x0FFF
		SetIndex_ANNN(value)
	case 0xD000:
		var xRegister = int((instructionBytes & 0x0F00) >> 8)
		var yRegister = int((instructionBytes & 0x00F0) >> 4)
		var height = int(instructionBytes & 0x000F)
		Display_DXYN(xRegister, yRegister, height)
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
	//fmt.Printf("00E0_ClearScreen\n")
}

func Jump_1NNN(jumpAddress uint16) {
	//var oldAddress = programCounter
	programCounter = jumpAddress
	//fmt.Printf("1NNN_Jump, old address: %d, new address: %d \n", oldAddress, jumpAddress)
}

func Set_6XNN(idx int, value uint8) {
	//var oldValue = vRegisters[idx]
	var newValue = value
	vRegisters[idx] = newValue
	//fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d \n", idx, oldValue, newValue)
}

func Add_7XNN(idx int, value uint8) {
	var oldValue = vRegisters[idx]
	var newValue = oldValue + value

	vRegisters[idx] = newValue
	//fmt.Printf("7XNN_Add, idx %d, oldValue %d, new value %d \n", idx, oldValue, newValue)
}

func SetIndex_ANNN(value uint16) {
	//var oldValue = indexRegister
	var newValue = value
	indexRegister = newValue
	//fmt.Printf("ANNN_SetIndex, oldValue %d, new value %d \n", oldValue, newValue)
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

	//fmt.Printf("DXYN_Display at xReg %d, yReg %d, height %d \n", xRegister, yRegister, spriteHeight)
}

func Subroutine_2NNN(value uint16) {
	addressStack.Push(programCounter)
	programCounter = value
}

func Subroutine_00EE() {
	var address = addressStack.Pop()
	programCounter = address
}

func Skip_conditionally_3XNN(idx int, value uint8) {
	var registerValue = vRegisters[idx]
	if registerValue == value {
		programCounter += 2
	}
}

func Skip_conditionally_4XNN(idx int, value uint8) {
	var registerValue = vRegisters[idx]
	if registerValue != value {
		programCounter += 2
	}
}

func Skip_conditionally_5XY0(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	if xValue == yValue {
		programCounter += 2
	}
}

func Skip_conditionally_9XY0(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var yValue = vRegisters[yIdx]
	if xValue != yValue {
		programCounter += 2
	}
}

func Set_8XY0(xIdx int, yIdx int) {
	var yValue = vRegisters[yIdx]
	vRegisters[xIdx] = yValue
}

func Binary_OR_8XY1(xIdx int, yIdx int) {
	var yValue = vRegisters[xIdx] | vRegisters[yIdx]
	vRegisters[xIdx] = yValue
}

func Binary_AND_8XY2(xIdx int, yIdx int) {
	var yValue = vRegisters[xIdx] & vRegisters[yIdx]
	vRegisters[xIdx] = yValue
}

func Binary_XOR_8XY3(xIdx int, yIdx int) {
	var yValue = vRegisters[xIdx] ^ vRegisters[yIdx]
	vRegisters[xIdx] = yValue
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

func Shift_8XY6(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var newValue = xValue >> 2
	var shiftedBit = xValue & (1 << uint(7))
	vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}

func Shift_8XYE(xIdx int, yIdx int) {
	var xValue = vRegisters[xIdx]
	var newValue = xValue << 2
	var shiftedBit = xValue & 0x00F
	vRegisters[xIdx] = newValue
	var isCarryFlagSet = shiftedBit > 0

	if isCarryFlagSet {
		vRegisters[15] = 1
	} else {
		vRegisters[15] = 0
	}
}
