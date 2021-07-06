package main

import (
	"encoding/gob"
	"os"
)

const (
	MAP_SIZE     = 81
	MAP_CENTRE   = 40
	MAP_HEIGHT   = 4
	DEFAULT_FLAT = 25
	DEFAULT_WALL = 8

	CEILING_BIT = 0b10000
	HIGH_BIT    = 0b01000
	LOW_BIT     = 0b00100
	WALL_BIT    = 0b00010
	GROUND_BIT  = 0b00001

	SKY                      = 0b00001
	SKY_SINGLE_BLOCK         = 0b00011
	SKY_FLOATING_BLOCK       = 0b00101
	SKY_DOUBLE_BLOCK         = 0b00111
	LOW_ROOM                 = 0b11001
	LOW_ROOM_SINGLE_BLOCK    = 0b11011
	CORRIDOR                 = 0b11101
	WALL                     = 0b11111
	HIGH_ROOM                = 0b10001
	HIGH_ROOM_SINGLE_BLOCK   = 0b10011
	HIGH_ROOM_FLOATING_BLOCK = 0b10101
	HIGH_ROOM_DOUBLE_BLOCK   = 0b10111

	GENERATOR_WALL = WALL
	GENERATOR_PATH = CORRIDOR
	GENERATOR_ROOM = HIGH_ROOM
)

type mapCell struct {
	cellType int
	flats    [4]int
	walls    [3]int
}

var grid [MAP_SIZE][MAP_SIZE]mapCell

func makeMap() {

	for x := -MAP_CENTRE; x <= MAP_CENTRE; x++ {
		for z := -MAP_CENTRE; z <= MAP_CENTRE; z++ {

			if x == 0 {

				grid[x+MAP_CENTRE][z+MAP_CENTRE] = mapCell{
					SKY_SINGLE_BLOCK,
					[4]int{0, 2, 0, 0},
					[3]int{3, 0, 0},
				}

			} else {

				grid[x+MAP_CENTRE][z+MAP_CENTRE] = mapCell{
					SKY,
					[4]int{1, 0, 0, 0},
					[3]int{0, 0, 0},
				}

			}

		}
	}

}

/*
func kruskalStep()  {

  oneSet := true;
  for i := 0; i < mazeWidth; i {
    for j := 0; j < mazeHeight; j {
      if kruskalMaze[i][j].set != kruskalMaze[0][0].set {
        oneSet = false;
      }
    }
  }
  if oneSet {
		return true;
	}

	x := 0
	y := 0
	a := 0
	b := 0
	horizontal := 0
	vertical := 0;

  for {
    x = rand() % mazeWidth;
    y = rand() % mazeHeight;

    if (rand() % 2 == 0) {
      horizontal = 1;
      vertical = 0;
    } else {
      horizontal = 0;
      vertical = 1;
    }

    if (horizontal > 0 && (kruskalMaze[x][y].right || x == mazeWidth - 1)) {
			continue
		}
    if (vertical > 0 && (kruskalMaze[x][y].down || y == mazeHeight - 1)) {
			continue
	}

    a = kruskalMaze[x][y].set
    b = kruskalMaze[x + horizontal][y + vertical].set

    if a == b {
			continue
		}

    if (vertical > 0) {
      kruskalMaze[x][y].down = true
    } else {
      kruskalMaze[x][y].right = true
    }
    for i = 0; i < mazeWidth; i {
      for j = 0; j < mazeHeight; j++) {
        if (kruskalMaze[i][j].set == b) {
          kruskalMaze[i][j].set = a
        }
      }
    }

    return false
  }
}

func makeMaze() {

		n := 0;
		for i := 0; i < mazeWidth; i++ {
			for j := 0; j < mazeHeight; j++ {
				n++;
				kruskalMaze[i][j].set = n;
				kruskalMaze[i][j].right = false;
				kruskalMaze[i][j].down = false;
			}
		}

		mazeDone := false;
		for !mazeDone {
			mazeDone = kruskalStep();
		}

		for i = -mazeWidth; i <= mazeWidth; i++ {
			for j = -mazeHeight; j <= mazeHeight; j++ {
				map[i + mazeWidth][j + mazeHeight].type =
					(i + mazeWidth) % 2 == 1 && (j + mazeHeight) % 2 == 1
						? GENERATOR_PATH
						: GENERATOR_WALL;
				for f = 0; f < 4; f++ {
					map[i + mazeWidth][j + mazeHeight].flat[f] = DEFAULT_FLAT;
				}
				for f = 0; f < 3; f++ {
					for d = 0; d < 4; d++ {
						map[i + mazeWidth][j + mazeHeight].wall[f][d] = DEFAULT_WALL;
					}
				}
			}
		}

		for i = 0; i < mazeWidth; i++ {
			for j = 0; j < mazeHeight; j++ {
				if (kruskalMaze[i][j].right) {
					map[i * 2 + 2][j * 2 + 1].type = GENERATOR_PATH;
					for f = 0; f < 4; f++ {
						map[i * 2 + 2][j * 2 + 1].flat[f] = DEFAULT_FLAT;
					}
					for f = 0; f < 3; f++ {
						for d = 0; d < 4; d++ {
							map[i * 2 + 2][j * 2 + 1].wall[f][d] = DEFAULT_WALL;
						}
					}
				}
				if (kruskalMaze[i][j].down) {
					map[i * 2 + 1][j * 2 + 2].type = GENERATOR_PATH;
					for f = 0; f < 4; f++ {
						map[i * 2 + 1][j * 2 + 2].flat[f] = DEFAULT_FLAT;
					}
					for f = 0; f < 3; f++ {
						for d = 0; d < 4; d++ {
							map[i * 2 + 1][j * 2 + 2].wall[f][d] = DEFAULT_WALL;
						}
					}
				}
			}
		}

		for r = 0; r < 10; r++ {
			let rw = rand() % 10 + 5;
			let rh = rand() % 10 + 5;
			let x = rand() % (mazeWidth * 2 - rw);
			let y = rand() % (mazeHeight * 2 - rh);
			if (r === 0) {
				playerX = (x + rw / 2 - mazeWidth) * unit;
				playerY = (y + rh / 2 - mazeHeight) * unit;
			}
			for i = x; i <= x + rw; i++ {
				for j = y; j <= y + rh; j++ {
					map[i][j].type = GENERATOR_ROOM;
					for f = 0; f < 4; f++ {
						map[i][j].flat[f] = DEFAULT_FLAT;
					}
					for f = 0; f < 3; f++ {
						for d = 0; d < 4; d++ {
							map[i][j].wall[f][d] = DEFAULT_WALL;
						}
					}
				}
			}
		}
	}
*/

func loadMap(filename string) {

	f1, err1 := os.Open(filename)
	if err1 == nil {
		decoder1 := gob.NewDecoder(f1)
		err := decoder1.Decode(&grid)
		if err != nil {
			panic(err)
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
