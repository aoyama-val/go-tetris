package model

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

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

func createRandomBlock(rng *rand.Rand, createCount uint32) Block {
	return Block{
		X:     4,
		Y:     0,
		Shape: Shape(rng.Intn(SHAPE_MAX + 1)),
		Rot:   0,
		Color: uint8(createCount % 3),
	}
}

// 壁と床を含めた堆積物を表す構造体
// 壁と床は別にした方が良かったかも
type Piles struct {
	Pattern [BOARD_Y_LEN][BOARD_X_LEN]uint8 // 0:なし 1:壁or床 2〜:ブロック残骸
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
	Rng               *rand.Rand
	IsOver            bool
	Frame             int32
	Piles             Piles
	Block             Block
	NextBlock         Block
	BlockCreatedCount uint32
}

func NewGame() *Game {
	timestamp := time.Now().Unix()
	rng := rand.New(rand.NewSource(timestamp))

	g := &Game{}
	g.Rng = rng
	g.IsOver = false
	g.Frame = 0
	g.Piles.SetupWallAndFloor()
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
		fmt.Printf("seed = %d\n", seed)
		rng := rand.New(rand.NewSource(seed))
		g.Rng = rng
	}

	pattern, ok := config.Get("pattern").(string)
	if ok {
		var p Piles
		scanner := bufio.NewScanner(strings.NewReader(pattern))
		i := 0
		for scanner.Scan() {
			line := scanner.Text()
			if err := scanner.Err(); err != nil {
				msg := fmt.Sprintf("Error occurred: %v\n", err)
				panic(msg)
			}
			cols := strings.Fields(line)
			for j, col := range cols {
				num, err := strconv.Atoi(col)
				if err != nil {
					msg := fmt.Sprintf("Error occurred while parsing line: i=%d: line=%s", i, line)
					panic(msg)
				}
				p.Pattern[i][j] = uint8(num)
			}
			i += 1
		}
		g.Piles = p
		fmt.Printf("pattern:%s", pattern)
	}
}

func (g *Game) InitRandomly() {
	for i := 0; i < 2; i++ {
		g.spawnBlock()
	}
}

func (g *Game) Update(command string) {
	if g.IsOver {
		return
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
		if g.isCollide(0, 1) {
			// すでに床に接触しているなら固定
			g.settleBlock()
			g.checkEraseRow()
			g.spawnBlock()
			if g.isCollide(0, 0) {
				g.IsOver = true
				println("Game over!")
			}
		} else {
			g.moveByDelta(0, 1)
		}
	}

	g.Frame += 1
}

func (g *Game) isCollide(xDelta int32, yDelta int32) bool {
	pattern := g.Block.GetPattern()
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if pattern[i][j] != 0 {
				newX := g.Block.X + int32(j) + xDelta
				newY := g.Block.Y + int32(i) + yDelta
				if g.Piles.isFilled(uint(newX), uint(newY)) {
					return true
				}
			}
		}
	}
	return false
}

func (g *Game) settleBlock() {
	blockPattern := g.Block.GetPattern()
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if blockPattern[i][j] == 1 {
				g.Piles.Pattern[uint(g.Block.Y+int32(i))][uint(g.Block.X+int32(j))] = 2 + g.Block.Color
			}
		}
	}
}

func (g *Game) moveByDelta(xDelta int32, yDelta int32) {
	if !g.isCollide(xDelta, yDelta) {
		g.Block.moveByDelta(xDelta, yDelta)
	}
}

func (g *Game) rotate(dir int32) {
	g.Block.rotate(dir)
	if g.isCollide(0, 0) {
		g.Block.rotate(-dir)
	}
}

func (g *Game) spawnBlock() {
	g.Block = g.NextBlock
	g.NextBlock = createRandomBlock(g.Rng, g.BlockCreatedCount)
	g.BlockCreatedCount += 1
}

func (g *Game) checkEraseRow() {
	filledRows := g.getFilledRows()
	if len(filledRows) > 0 {
		// そろった行を消す
		maxFilledRow := filledRows[len(filledRows)-1]
		for y := int(maxFilledRow); y >= 0; y-- {
			for x := 1; x <= int(BOARD_X_MAX)-1; x++ {
				above := int(y) - len(filledRows)
				if above >= 0 {
					g.Piles.Pattern[y][x] = g.Piles.Pattern[above][x]
				} else {
					g.Piles.Pattern[y][x] = 0
				}
			}
		}
	}
}

func (g *Game) getFilledRows() []uint {
	var result []uint
	for y := BOARD_Y_MIN; y <= BOARD_Y_MAX-1; y++ {
		isFilled := true
		for x := 1; x <= int(BOARD_X_MAX)-1; x++ {
			if !g.Piles.isFilled(uint(x), uint(y)) {
				isFilled = false
			}
		}
		if isFilled {
			result = append(result, uint(y))
		}
	}
	return result
}
