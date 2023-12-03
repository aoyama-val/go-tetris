package model

import (
	"fmt"
	"math/rand"
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
	base := b.Shape.getBasePattern()
	return base
	// TODO: 回転対応
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

func (g *Game) Update(command string) {
	if command != "" {
		fmt.Println("command=", command)
		switch command {
		case "left":
			g.Block.X -= 1
		case "right":
			g.Block.X += 1
		}
	}
}

func (g *Game) spawnBlock() {
	g.Block = g.NextBlock
	g.NextBlock = createRandomBlock(g.BlockCreatedCount)
	g.BlockCreatedCount += 1
}
