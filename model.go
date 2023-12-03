package main

import "fmt"

type Shape int

const (
	SHAPE_MAX = 6
)

type Block struct {
	x     int32
	y     int32
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
	block             Block
	nextBlock         Block
	blockCreatedCount uint32
}

func (g Game) Update(command string) {
	fmt.Println("command=", command)
	sage()
}
