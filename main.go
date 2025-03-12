package main

import (
	"fmt"
	"go-chip/cpu"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const timerDecreaseSpeed = 60.0
const instructionExecutionSpeed = 700.0

const width = 64
const height = 32

var cpuInstance cpu.CpuInstance
var displayRects [width][height]sdl.Rect
var blackColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var whiteColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}
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

	for {
		cpuInstance.ExecuteTimersUpdate()
		inputLoop()
		if !isRunning {
			return
		}
		cpuInstance.ExecuteLoopStep()

		for i := range height {
			for j := range width {
				if cpuInstance.GetDisplay(j, i) {
					windowSurface.FillRect(&displayRects[j][i], whiteColor.Uint32())
				} else {
					windowSurface.FillRect(&displayRects[j][i], blackColor.Uint32())
				}
			}
		}

		window.UpdateSurface()
		sdl.Delay(1000 / 700)
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
