package main

import (
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

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("go-tetris", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	err = renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	running := true
	game := m.NewGame()
	game.LoadConfig()
	game.InitRandomly()

	for running {
		command := ""
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED && t.Repeat == 0 {
					keyCode := t.Keysym.Sym
					switch keyCode {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_LEFT:
						command = "left"
					case sdl.K_RIGHT:
						command = "right"
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
		render(renderer, window, game)
		time.Sleep((1000 / FPS) * time.Millisecond)
	}
}

func render(renderer *sdl.Renderer, window *sdl.Window, game *m.Game) {
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.Clear()

	// render piles
	for i := 0; i < len(game.Piles.Pattern); i++ {
		for j := 0; j < len(game.Piles.Pattern[i]); j++ {
			if game.Piles.Pattern[i][j] >= 1 {
				color := getColor(game.Piles.Pattern[i][j])
				renderer.SetDrawColor(color.R, color.G, color.B, color.A)
				var x int32 = (m.LEFT_WALL_X + int32(j)) * CELL_SIZE_PX
				var y int32 = int32(i) * CELL_SIZE_PX
				rect := sdl.Rect{X: x, Y: y, W: CELL_SIZE_PX, H: CELL_SIZE_PX}
				renderer.FillRect(&rect)
			}
		}

	}
	// render block
	renderBlock(renderer, &game.Block, m.LEFT_WALL_X+game.Block.X, game.Block.Y)

	// render next block
	renderBlock(renderer, &game.NextBlock, 21, 0)

	if game.IsOver {
		renderer.SetDrawColor(0, 0, 0, 128)
		renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: SCREEEN_WIDTH, H: SCREEN_HEIGHT})
	}

	renderer.Present()
}

func renderBlock(r *sdl.Renderer, block *m.Block, xInCell int32, yInCell int32) {
	color := getColor(block.Color + 2)
	r.SetDrawColor(color.R, color.G, color.B, color.A)
	pattern := block.GetPattern()
	for j := 0; j < len(pattern); j++ {
		for i := 0; i < len(pattern[j]); i++ {
			if pattern[j][i] == 1 {
				x := (xInCell + int32(i)) * CELL_SIZE_PX
				y := (yInCell + int32(j)) * CELL_SIZE_PX
				r.FillRect(&sdl.Rect{X: x, Y: y, W: CELL_SIZE_PX, H: CELL_SIZE_PX})
			}
		}
	}
}

func getColor(colorNum uint8) sdl.Color {
	switch colorNum {
	case 0:
		return sdl.Color{R: 0, G: 0, B: 0}
	case 1:
		return sdl.Color{R: 128, G: 128, B: 128, A: 255}
	case 2:
		return sdl.Color{R: 255, G: 128, B: 128, A: 255}
	case 3:
		return sdl.Color{R: 128, G: 255, B: 128, A: 255}
	case 4:
		return sdl.Color{R: 128, G: 128, B: 255, A: 255}
	default:
		return sdl.Color{R: 255, G: 255, B: 255, A: 255}
	}
}
