package display

import (
	"go-chip/cpu"

	"github.com/veandco/go-sdl2/sdl"
)

var blackColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
var whiteColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}

const pixelSizeScaleX = 10
const pixelSizeScaleY = 10

const pixelWidth = 10
const pixelHeight = 10

type DisplayInstance struct {
	cpu           *cpu.CpuInstance
	width         int
	height        int
	displayRects  [][]sdl.Rect
	window        *sdl.Window
	windowSurface *sdl.Surface
}

func (display *DisplayInstance) Init(cpuInstance *cpu.CpuInstance) {
	display.cpu = cpuInstance
	display.width = cpu.DisplayWidth * pixelSizeScaleX
	display.height = cpu.DisplayHeight * pixelSizeScaleY

	display.window, _ = sdl.CreateWindow("Go-Chip by JinnFighter", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, int32(display.width), int32(display.height), sdl.WINDOW_SHOWN)
	display.windowSurface, _ = display.window.GetSurface()
	display.windowSurface.FillRect(nil, 0)

	display.displayRects = make([][]sdl.Rect, cpu.DisplayWidth)
	for i := range cpu.DisplayWidth {
		display.displayRects[i] = make([]sdl.Rect, cpu.DisplayHeight)
		for j := range cpu.DisplayHeight {
			display.displayRects[i][j] = sdl.Rect{X: int32(i * pixelSizeScaleX), Y: int32(j * pixelSizeScaleY), W: pixelWidth, H: pixelHeight}
		}
	}
}

func (display *DisplayInstance) Terminate() {
	display.window.Destroy()
}

func (display *DisplayInstance) Draw() {
	for i := range cpu.DisplayHeight {
		for j := range cpu.DisplayWidth {
			if display.cpu.GetDisplay(j, i) {
				display.windowSurface.FillRect(&display.displayRects[j][i], whiteColor.Uint32())
			} else {
				display.windowSurface.FillRect(&display.displayRects[j][i], blackColor.Uint32())
			}
		}
	}

	display.window.UpdateSurface()
}
