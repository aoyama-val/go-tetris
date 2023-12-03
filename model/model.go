package model

import "fmt"

const (
	BOARD_X_LEN uint32 = 12
	BOARD_X_MIN uint32 = 0
	BOARD_X_MAX uint32 = BOARD_X_LEN - 1
	BOARD_Y_LEN uint32 = 21
	BOARD_Y_MIN uint32 = 0
	BOARD_Y_MAX uint32 = BOARD_Y_LEN - 1
	LEFT_WALL_X int32  = 6
)

type Shape int

const (
	SHAPE_MAX = 6
)

type Block struct {
	X     int32
	Y     int32
	Shape Shape
	Rot   int8
	Color uint8
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
	// g.NextBlock = createRandomBlock()
	g.BlockCreatedCount = 0
	return g
}

func (g Game) Update(command string) {
	if command != "" {
		fmt.Println("command=", command)
	}
}
