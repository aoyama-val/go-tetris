package main

import (
	"fmt"
	"time"

	m "github.com/aoyama-val/go-tetris/model"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCREEEN_WIDTH = 640
	SCREEN_HEIGHT = 420
	CELL_SIZE_PX  = 20
	FPS           = 32
)

func sage() {
	println("sage")
}

func main() {
	fmt.Println("Hello world!")

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("go-tetris", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	running := true
	var game m.Game
	var command string

	for running {
		command = ""
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED && t.Repeat == 0 {
					keyCode := t.Keysym.Sym
					switch keyCode {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_LEFT:
						command = "left"
						game.Block.X -= 1
					case sdl.K_RIGHT:
						command = "right"
						game.Block.X += 1
					case sdl.K_DOWN:
						command = "down"
					case sdl.K_z:
						command = "rotate_left"
					case sdl.K_x:
						command = "rotate_right"
					}
				}
			}
		}
		game.Update(command)
		if command != "" {
			fmt.Println("command=", command)
		}
		render(surface, window, &game)
		time.Sleep((1000 / FPS) * time.Millisecond)
	}
}

func render(surface *sdl.Surface, window *sdl.Window, game *m.Game) {
	surface.FillRect(nil, 0)
	rect := sdl.Rect{game.Block.X, game.Block.Y, 200, 200}
	colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
	window.UpdateSurface()
}
