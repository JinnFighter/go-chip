package cpu

import (
	"fmt"
	"go-chip/extensions"
)

const DisplayWidth = 64
const DisplayHeight = 32
const memorySize = 4096
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
	display        [DisplayWidth][DisplayHeight]bool
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
	for i := range DisplayWidth {
		for j := range DisplayHeight {
			cpu.display[i][j] = false
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
	var instruction = cpu.decodeInstruction(nextInstruction)
	if instruction != nil {
		instruction.SetupValues(nextInstruction)
		instruction.Execute(cpu)
	} else {
		fmt.Println("unknown command")
	}
}

func (cpu *CpuInstance) GetDisplay(i int, j int) bool {
	return cpu.display[i][j]
}

func (cpu *CpuInstance) decodeInstruction(instructionBytes uint16) IInstruction {
	var firstByte = instructionBytes & 0xF000
	var instruction = cpu.instructions[firstByte]
	return instruction
}
