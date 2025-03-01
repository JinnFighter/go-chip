package main

import (
	"fmt"
	"go-chip/extensions"
	"time"
)

const timerDecreaseSpeed = 60.0
const instructionExecutionSpeed = 700.0

var memory [4096]byte
var display [64][32]bool
var vRegisters [16]uint8
var indexRegister uint16
var programCounter uint16
var addressStack extensions.Stack
var delayTimer uint8
var soundTimer uint8
var isRunning bool
var ticker *time.Ticker
var tickerChannel chan bool

func main() {
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
	fmt.Println("Tick duration: ", time.Duration(execSpeed))
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
		case t := <-ticker.C:
			fmt.Println("Tick at ", t)

			var nextInstruction = programCounter
			fmt.Println("Next instruction counter: ", nextInstruction)
			programCounter += 2
			decodeInstruction(nextInstruction)
		}
	}
}

func decodeInstruction(instructionBytes uint16) {
	var firstByte = instructionBytes & 0xF000
	fmt.Println("Current First Byte is ", firstByte)
	fmt.Print("Current command is: ")
	switch firstByte {
	case 0x0000:
		ClearScreen_00E0()
	case 0x1000:
		var jumpAddress = instructionBytes & 0x0FFF
		Jump_1NNN(jumpAddress)
	case 0x6000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Set_6XNN(idx, value)
	case 0x7000:
		var idx = int((instructionBytes & 0x0F00) >> 8)
		var value = uint8(instructionBytes & 0x00FF)
		Add_7XNN(idx, value)
	case 0xA000:
		var value = instructionBytes & 0x0FFF
		SetIndex_ANNN(value)
	case 0xD000:
		var xRegister = int((instructionBytes & 0x0F00) >> 8)
		var yRegister = int((instructionBytes & 0x00F0) >> 4)
		var height = int(instructionBytes & 0x000F)
		Display_DXYN(xRegister, yRegister, height)
	default:
		fmt.Println("Unknown Command")
	}
}

func ClearScreen_00E0() {
	fmt.Println("00E0_ClearScreen")
	for i := range 64 {
		for j := range 32 {
			display[i][j] = false
		}
	}
}

func Jump_1NNN(jumpAddress uint16) {
	fmt.Println("1NNN_Jump, address: ", jumpAddress)
	programCounter = jumpAddress
}

func Set_6XNN(idx int, value uint8) {
	var oldValue = vRegisters[idx]
	var newValue = value
	vRegisters[idx] = newValue
	fmt.Printf("6XNN_Set, idx %d, OldValue %d, new value %d", idx, oldValue, newValue)
}

func Add_7XNN(idx int, value uint8) {
	var oldValue = vRegisters[idx]
	var newValue = oldValue + value

	vRegisters[idx] = newValue
	fmt.Printf("7XNN_Add, idx %d, oldValue %d, new value %d", idx, oldValue, newValue)
}

func SetIndex_ANNN(value uint16) {
	var oldValue = indexRegister
	var newValue = value
	indexRegister = value
	fmt.Printf("ANNN_SetIndex, oldValue %d, new value %d", oldValue, newValue)
}

func Display_DXYN(xRegister int, yRegister int, pixelCount int) {
	var x = vRegisters[xRegister] % 64
	var y = vRegisters[yRegister] % 32
	vRegisters[15] = 0

	var pixelStartAddress = indexRegister
	for i := range pixelCount {
		var spriteByte = memory[pixelStartAddress+uint16(i)]
		for j := range 8 {
			var spritePixel = spriteByte & (1 << j)
			if spritePixel > 1 && display[x][y] {
				display[x][y] = false
				vRegisters[15] = 1
			} else if spritePixel > 1 && !display[x][y] {
				display[x][y] = true
			}
			x++
			if x >= 64 {
				break
			}
		}
		y++
		if y >= 32 {
			break
		}
	}

	fmt.Printf("DXYN_Display at xReg %d, yReg %d, height %d", xRegister, yRegister, pixelCount)
}
