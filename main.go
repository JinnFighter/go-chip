package main

import (
	"fmt"
	"go-chip/cpu"
	"go-chip/display"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

var cpuInstance cpu.CpuInstance
var displayInstance display.DisplayInstance
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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	var args = os.Args[1:]
	data, err := os.ReadFile(args[0])
	if err != nil {
		panic(err)
	}

	cpuInstance.Init(data)

	displayInstance.Init(&cpuInstance)

	startLoop()
	stopLoop()
	displayInstance.Terminate()
	cpuInstance.Terminate()
}

func startLoop() {
	cpuInstance.StartLoop()
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
	cpuInstance.StopLoop()
}

func loop() {
	fmt.Println("Enter loop")
	for range time.Tick(time.Duration(1000 / 60 * time.Millisecond)) {
		cpuInstance.ExecuteTimersUpdate()
		inputLoop()
		if !isRunning {
			return
		}
		cpuInstance.ExecuteLoopStep()

		displayInstance.Draw()
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
					cpuInstance.SetInputPressed(keyByte, true)
				} else if et.Type == sdl.KEYUP {
					cpuInstance.SetInputPressed(keyByte, false)
				}
			}
		}
	}
}
