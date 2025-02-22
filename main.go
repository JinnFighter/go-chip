package main

import (
	"fmt"
	"go-chip/extensions"
)

const timerDecreaseSpeed = 60

var memory [4096]byte
var display [64][32]bool
var addressStack extensions.Stack
var delayTimer byte
var soundTimer byte

func main() {
	fmt.Println("hello, world!")
}
