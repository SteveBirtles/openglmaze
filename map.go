package main

import (
	"math/rand"
)

const (
	MAP_SIZE   int = 81
	MAP_CENTRE int = 40
	MAP_HEIGHT int = 4

	DEFAULT_FLAT int = 26
	DEFAULT_WALL int = 10

	CEILING_BIT uint = 0b10000
	HIGH_BIT    uint = 0b01000
	LOW_BIT     uint = 0b00100
	WALL_BIT    uint = 0b00010
	GROUND_BIT  uint = 0b00001

	SKY                      uint = 0b00001
	SKY_SINGLE_BLOCK         uint = 0b00011
	SKY_FLOATING_BLOCK       uint = 0b00101
	SKY_DOUBLE_BLOCK         uint = 0b00111
	LOW_ROOM                 uint = 0b11001
	LOW_ROOM_SINGLE_BLOCK    uint = 0b11011
	CORRIDOR                 uint = 0b11101
	WALL                     uint = 0b11111
	HIGH_ROOM                uint = 0b10001
	HIGH_ROOM_SINGLE_BLOCK   uint = 0b10011
	HIGH_ROOM_FLOATING_BLOCK uint = 0b10101
	HIGH_ROOM_DOUBLE_BLOCK   uint = 0b10111

	GENERATOR_WALL uint = WALL
	GENERATOR_PATH uint = CORRIDOR
	GENERATOR_ROOM uint = HIGH_ROOM
)

type mapCell struct {
	cellType uint
	flats    [4]int
	walls    [3]int
}

type kruskalCell struct {
	set   int
	right bool
	down  bool
}

var grid [MAP_SIZE][MAP_SIZE]mapCell
var kruskalMaze [MAP_CENTRE][MAP_CENTRE]kruskalCell

func clearMap() {

	for x := -MAP_CENTRE; x <= MAP_CENTRE; x++ {
		for z := -MAP_CENTRE; z <= MAP_CENTRE; z++ {
			grid[x+MAP_CENTRE][z+MAP_CENTRE] = mapCell{
				0b0,
				[4]int{0, 0, 0, 0},
				[3]int{0, 0, 0},
			}
		}
	}

}

func kruskalStep() bool {

	oneSet := true
	for i := 0; i < MAP_CENTRE; i++ {
		for j := 0; j < MAP_CENTRE; j++ {
			if kruskalMaze[i][j].set != kruskalMaze[0][0].set {
				oneSet = false
			}
		}
	}

	if oneSet {
		return true
	}

	x := 0
	y := 0
	a := 0
	b := 0
	horizontal := 0
	vertical := 0

	for {
		x = rand.Intn(MAP_CENTRE)
		y = rand.Intn(MAP_CENTRE)

		if rand.Intn(2) == 0 {
			horizontal = 1
			vertical = 0
		} else {
			horizontal = 0
			vertical = 1
		}

		if horizontal > 0 && (kruskalMaze[x][y].right || x == MAP_CENTRE-1) {
			continue
		}
		if vertical > 0 && (kruskalMaze[x][y].down || y == MAP_CENTRE-1) {
			continue
		}

		a = kruskalMaze[x][y].set
		b = kruskalMaze[x+horizontal][y+vertical].set

		if a == b {
			continue
		}

		if vertical > 0 {
			kruskalMaze[x][y].down = true
		} else {
			kruskalMaze[x][y].right = true
		}
		for i := 0; i < MAP_CENTRE; i++ {
			for j := 0; j < MAP_CENTRE; j++ {
				if kruskalMaze[i][j].set == b {
					kruskalMaze[i][j].set = a
				}
			}
		}

		return false
	}
}

func makeMaze() {

	clearMap()

	n := 0
	for i := 0; i < MAP_CENTRE; i++ {
		for j := 0; j < MAP_CENTRE; j++ {
			n++
			kruskalMaze[i][j].set = n
			kruskalMaze[i][j].right = false
			kruskalMaze[i][j].down = false
		}
	}

	mazeDone := false
	for !mazeDone {
		mazeDone = kruskalStep()
	}

	for i := -MAP_CENTRE; i <= MAP_CENTRE; i++ {
		for j := -MAP_CENTRE; j <= MAP_CENTRE; j++ {
			if (i+MAP_CENTRE)%2 == 1 && (j+MAP_CENTRE)%2 == 1 {
				grid[i+MAP_CENTRE][j+MAP_CENTRE].cellType = GENERATOR_PATH
			} else {
				grid[i+MAP_CENTRE][j+MAP_CENTRE].cellType = GENERATOR_WALL
			}
			for f := 0; f < 4; f++ {
				grid[i+MAP_CENTRE][j+MAP_CENTRE].flats[f] = DEFAULT_FLAT
				if f == 3 {
					break
				}
				grid[i+MAP_CENTRE][j+MAP_CENTRE].walls[f] = DEFAULT_WALL
			}
		}
	}

	for i := 0; i < MAP_CENTRE; i++ {
		for j := 0; j < MAP_CENTRE; j++ {
			if kruskalMaze[i][j].right {
				grid[i*2+2][j*2+1].cellType = GENERATOR_PATH
				for f := 0; f < 4; f++ {
					grid[i*2+2][j*2+1].flats[f] = DEFAULT_FLAT
					if f == 3 {
						break
					}
					grid[i*2+2][j*2+1].walls[f] = DEFAULT_WALL
				}
			}
			if kruskalMaze[i][j].down {
				grid[i*2+1][j*2+2].cellType = GENERATOR_PATH
				for f := 0; f < 4; f++ {
					grid[i*2+1][j*2+2].flats[f] = DEFAULT_FLAT
					if f == 3 {
						break
					}
					grid[i*2+1][j*2+2].walls[f] = DEFAULT_WALL
				}
			}
		}
	}

	for r := 0; r < 10; r++ {
		rw := rand.Intn(10) + 5
		rh := rand.Intn(10) + 5
		x := rand.Intn(MAP_CENTRE*2 - rw)
		y := rand.Intn(MAP_CENTRE*2 - rh)
		/*if r == 0 {
			playerX = (x + rw/2 - MAP_CENTRE) * unit
			playerY = (y + rh/2 - MAP_CENTRE) * unit
		}*/
		for i := x; i <= x+rw; i++ {
			for j := y; j <= y+rh; j++ {
				grid[i][j].cellType = GENERATOR_ROOM
				for f := 0; f < 4; f++ {
					grid[i][j].flats[f] = DEFAULT_FLAT
				}
				for f := 0; f < 3; f++ {
					grid[i][j].walls[f] = DEFAULT_WALL
				}
			}
		}
	}
}

func calculateMapShadow(x float64, y float64, z float64, frontTile uint16) bool {

	/*for s := 1.0; y+s < MAP_HEIGHT; s++ {

		if int(z-s) >= -MAP_CENTRE && int(z-s) < MAP_CENTRE {

			if frontTile == 0 &&
				(grid[int(x)+MAP_CENTRE][int(z-s)+MAP_CENTRE][int(y+s-1)][1] > 0 || grid[int(x)+MAP_CENTRE][int(z-s)+MAP_CENTRE][int(y+s)][0] > 0) ||
				frontTile > 0 && s > 1 &&
					(grid[int(x)+MAP_CENTRE][int(z-s)+MAP_CENTRE][int(y+s)][0] > 0 || grid[int(x)+MAP_CENTRE][int(z-s+1)+MAP_CENTRE][int(y+s)][0] > 0) {
				return true
			}

		}
	}*/

	return false

}
