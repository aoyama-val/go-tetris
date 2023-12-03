package model

import "fmt"

type Shape int

const (
	SHAPE_MAX = 6
)

type Block struct {
	X     int32
	Y     int32
	shape Shape
	rot   int8
	color uint8
}

type Piles struct {
}

type Game struct {
	isOver            bool
	frame             int32
	settleWait        uint32
	piles             Piles
	Block             Block
	nextBlock         Block
	blockCreatedCount uint32
}

func (g Game) Update(command string) {
	if command != "" {
		fmt.Println("command=", command)
	}
}
