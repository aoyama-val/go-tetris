package model

import "fmt"

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

func (g Game) Update(command string) {
	if command != "" {
		fmt.Println("command=", command)
	}
}
