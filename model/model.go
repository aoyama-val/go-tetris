package model

import (
	"fmt"
	"math/rand"

	"github.com/pelletier/go-toml"
)

const (
	BOARD_X_LEN uint32 = 12
	BOARD_X_MIN uint32 = 0
	BOARD_X_MAX uint32 = BOARD_X_LEN - 1
	BOARD_Y_LEN uint32 = 21
	BOARD_Y_MIN uint32 = 0
	BOARD_Y_MAX uint32 = BOARD_Y_LEN - 1
	LEFT_WALL_X int32  = 6
)

type Pattern [5][5]uint8

type Shape int

const (
	SHAPE_S0  = 0
	SHAPE_S1  = 1
	SHAPE_S2  = 2
	SHAPE_S3  = 3
	SHAPE_S4  = 4
	SHAPE_S5  = 5
	SHAPE_S6  = 6
	SHAPE_MAX = 6
)

func (s *Shape) getBasePattern() Pattern {
	switch *s {
	case 0:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
		}
	case 1:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
		}
	case 2:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
		}
	case 3:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}
	case 4:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 0, 0},
			{0, 1, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}
	case 5:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}
	case 6:
		return Pattern{
			{0, 0, 0, 0, 0},
			{0, 1, 1, 0, 0},
			{0, 0, 1, 1, 0},
			{0, 0, 0, 0, 0},
			{0, 0, 0, 0, 0},
		}
	default:
		msg := fmt.Sprintf("Invalid shape: %v", *s)
		panic(msg)
	}
}

type Block struct {
	X     int32
	Y     int32
	Shape Shape
	Rot   int8
	Color uint8
}

func (b *Block) GetPattern() Pattern {
	result := b.Shape.getBasePattern()
	for i := 0; i < int(b.Rot); i++ {
		result = rotatePattern(result)
	}
	return result
}

func rotatePattern(p Pattern) Pattern {
	result := p
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			result[4-j][i] = p[i][j]
		}
	}
	return result
}

func (b *Block) rotate(dir int32) {
	if dir > 0 {
		b.Rot = (b.Rot + 1) % 4
	} else {
		b.Rot = (b.Rot + 3) % 4
	}
}

func (b *Block) moveByDelta(xDelta int32, yDelta int32) {
	b.X += xDelta
	b.Y += yDelta
}

func createRandomBlock(createCount uint32) Block {
	return Block{
		X:     4,
		Y:     0,
		Shape: Shape(rand.Intn(SHAPE_MAX + 1)),
		Rot:   0,
		Color: uint8(createCount % 3),
	}
}

type Piles struct {
	Pattern [BOARD_Y_LEN][BOARD_X_LEN]uint8
}

func (p *Piles) SetupWallAndFloor() {
	for i := BOARD_Y_MIN; i <= BOARD_Y_MAX; i++ {
		p.Pattern[i][BOARD_X_MIN] = 1
		p.Pattern[i][BOARD_X_MAX] = 1
	}
	for i := BOARD_X_MIN; i <= BOARD_X_MAX; i++ {
		p.Pattern[BOARD_Y_MAX][i] = 1
	}
}

func (p *Piles) isFilled(x uint, y uint) bool {
	return p.Pattern[y][x] >= 1
}

type Game struct {
	IsOver            bool
	Frame             int32
	SettleWait        uint32
	Piles             Piles
	Block             Block
	NextBlock         Block
	BlockCreatedCount uint32
}

func NewGame() *Game {
	g := &Game{}
	g.IsOver = false
	g.Frame = 0
	g.SettleWait = 0
	g.Piles.SetupWallAndFloor()
	g.NextBlock = createRandomBlock(g.BlockCreatedCount)
	g.BlockCreatedCount += 1
	g.spawnBlock()
	return g
}

func (g *Game) LoadConfig() {
	filename := "tetris.toml"
	config, err := toml.LoadFile(filename)
	if err != nil {
		return
	}
	seed, ok := config.Get("seed").(int64)
	if ok {
		fmt.Printf("seed = %d", seed)
		rand.New(rand.NewSource(seed))
	}

	pattern, ok := config.Get("pattern").(string)
	if ok {
		fmt.Printf("pattern:%s", pattern)
		// TODO: set piles
	}
}

func (g *Game) Update(command string) {
	if g.IsOver {
		return
	}

	if g.SettleWait > 0 {
		g.SettleWait -= 1
		if g.SettleWait == 0 {
			// 床に接触した
			if g.isCollide() {
				g.settleBlock()
				g.spawnBlock()
			}
		}
	}

	switch command {
	case "left":
		g.moveByDelta(-1, 0)
	case "right":
		g.moveByDelta(1, 0)
	case "down":
		g.moveByDelta(0, 1)
	case "rotate_left":
		g.rotate(1)
	case "rotate_right":
		g.rotate(-1)
	}

	if g.Frame != 0 && g.Frame%20 == 0 {
		g.moveByDelta(0, 1)
		if g.isCollide() {
			println("Game over!")
			g.IsOver = true
		} else {
			g.SettleWait = 15
		}
	}
	g.checkEraseRow()

	g.Frame += 1
}

func (g *Game) isCollide() bool {
	pattern := g.Block.GetPattern()
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if pattern[i][j] != 0 {
				newX := g.Block.X + int32(j)
				newY := g.Block.Y + int32(i)
				if g.Piles.isFilled(uint(newX), uint(newY)) {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) settleBlock() {
	// TODO: 実装
}

func (g *Game) moveByDelta(xDelta int32, yDelta int32) {
	g.Block.moveByDelta(xDelta, yDelta)
	if g.isCollide() {
		g.Block.moveByDelta(-xDelta, -yDelta)
	}
}

func (g *Game) rotate(dir int32) {
	g.Block.rotate(dir)
	if g.isCollide() {
		g.Block.rotate(-dir)
	}
}

func (g *Game) spawnBlock() {
	g.Block = g.NextBlock
	g.NextBlock = createRandomBlock(g.BlockCreatedCount)
	g.BlockCreatedCount += 1
}

func (g *Game) checkEraseRow() {
	// TODO: 実装
}

func (g *Game) getFilledRows() {

}
