package main

import (
	"fmt"
	"go-chip/extensions"
	"time"
)

const timerDecreaseSpeed = 60.0
const instructionExecutionSpeed = 700.0

var memory [4096]uint8
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

}
