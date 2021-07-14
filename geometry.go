package main

import (
	"math"
)

const texZero = 0.0
const texOne = 1.0
const drawDistance float64 = 16

type vertexRecord struct {
	texture int
	index   int
	x       int
	y       int
	z       int
	wall    int
}

var (
	vertices     [][]float32
	verticesTemp [][]float32

	vertexRecords     []vertexRecord
	vertexRecordsTemp []vertexRecord

	vertexMutex = false

	cubeBottom = []float32{ //blue
		1.0, 0.0, 0.0, texOne, texZero,
		0.0, 0.0, 0.0, texZero, texZero,
		0.0, 0.0, 1.0, texZero, texOne,

		1.0, 0.0, 1.0, texOne, texOne,
		1.0, 0.0, 0.0, texOne, texZero,
		0.0, 0.0, 1.0, texZero, texOne,
	}
	cubeFlippedBottom = []float32{ //magenta
		1.0, 0.0, 0.0, texZero, texZero,
		0.0, 0.0, 1.0, texOne, texOne,
		0.0, 0.0, 0.0, texOne, texZero,

		1.0, 0.0, 1.0, texZero, texOne,
		0.0, 0.0, 1.0, texOne, texOne,
		1.0, 0.0, 0.0, texZero, texZero,
	}
	cubeDarkSide = []float32{ //cyan
		0.0, 0.0, 1.0, texZero, texOne,
		1.0, 0.0, 1.0, texOne, texOne,
		0.0, 1.0, 1.0, texZero, texZero,

		1.0, 0.0, 1.0, texOne, texOne,
		1.0, 1.0, 1.0, texOne, texZero,
		0.0, 1.0, 1.0, texZero, texZero,
	}
	cubeLightSide = []float32{ //green
		1.0, 0.0, 0.0, texZero, texOne,
		0.0, 1.0, 0.0, texOne, texZero,
		1.0, 1.0, 0.0, texZero, texZero,

		0.0, 0.0, 0.0, texOne, texOne,
		0.0, 1.0, 0.0, texOne, texZero,
		1.0, 0.0, 0.0, texZero, texOne,
	}
	cubeLeft = []float32{ //red
		0.0, 0.0, 1.0, texOne, texOne,
		0.0, 1.0, 0.0, texZero, texZero,
		0.0, 0.0, 0.0, texZero, texOne,

		0.0, 0.0, 1.0, texOne, texOne,
		0.0, 1.0, 1.0, texOne, texZero,
		0.0, 1.0, 0.0, texZero, texZero,
	}
	cubeRight = []float32{ //yellow
		1.0, 0.0, 1.0, texZero, texOne,
		1.0, 0.0, 0.0, texOne, texOne,
		1.0, 1.0, 0.0, texOne, texZero,

		1.0, 0.0, 1.0, texZero, texOne,
		1.0, 1.0, 0.0, texOne, texZero,
		1.0, 1.0, 1.0, texZero, texZero,
	}
)

func processVertex(v float32, index int, coords []int, texture int, rgb []float32, wall int) {

	if index%5 < 3 {
		v += float32(coords[index%5])
	}
	verticesTemp[texture] = append(verticesTemp[texture], v)
	if index%5 == 4 {

		verticesTemp[texture] = append(verticesTemp[texture], rgb...)

		vertexRecordsTemp = append(vertexRecordsTemp, vertexRecord{
			index:   len(verticesTemp[texture]) - 8,
			texture: texture,
			x:       coords[0],
			y:       coords[1],
			z:       coords[2],
			wall:    wall,
		})

	}

}

func updateWorld() {

	vertexMutex = true

	verticesTemp = make([][]float32, textureCount)
	for i := 0; i < textureCount; i++ {
		verticesTemp[i] = make([]float32, 0)
	}

	vertexRecordsTemp = make([]vertexRecord, 0)

	for x := int(math.Floor(myX) - drawDistance); x < int(math.Floor(myX)+drawDistance); x++ {
		for z := int(math.Floor(myZ) - drawDistance); z < int(math.Floor(myZ)+drawDistance); z++ {

			if x < 0 || z < 0 || x > MAP_SIZE-1 || z > MAP_SIZE-1 {
				continue
			}

			for y := 0; y < MAP_HEIGHT; y++ {

				coords := []int{x, y, z}

				var flatBit uint
				var wallBit uint
				switch y {
				case 0:
					flatBit = GROUND_BIT
					wallBit = WALL_BIT
				case 1:
					flatBit = WALL_BIT
					wallBit = LOW_BIT
				case 2:
					flatBit = LOW_BIT
					wallBit = HIGH_BIT
				case 3:
					flatBit = HIGH_BIT
					wallBit = CEILING_BIT
				}

				ambient := []float32{0.666, 0.666, 0.666}

				flatTexture := int(grid[x][z].flats[y]) - 1

				if flatTexture != -1 && grid[x][z].cellType&flatBit > 0 {
					if y < MAP_HEIGHT-1 && grid[x][z].cellType&wallBit == 0 {
						for i, v := range cubeBottom {
							processVertex(v, i, coords, flatTexture, ambient, -1)
						}
					}
				}

				if flatTexture != -1 && grid[x][z].cellType&flatBit == 0 && grid[x][z].cellType&wallBit > 0 {
					for i, v := range cubeFlippedBottom {
						processVertex(v, i, coords, flatTexture, ambient, -1)
					}
				}

				if y >= MAP_HEIGHT-1 || grid[x][z].cellType&wallBit == 0 {
					continue
				}

				if (x == 0 || x > 0 && grid[x-1][z].cellType&wallBit == 0) && int(grid[x][z].walls[y][0]) > 0 {
					for i, v := range cubeLeft {
						processVertex(v, i, coords, int(grid[x][z].walls[y][0])-1, ambient, 0)
					}
				}
				if (x == MAP_SIZE-1 || x < MAP_SIZE-1 && grid[x+1][z].cellType&wallBit == 0) && int(grid[x][z].walls[y][1]) > 0 {
					for i, v := range cubeRight {
						processVertex(v, i, coords, int(grid[x][z].walls[y][1])-1, ambient, 1)
					}
				}
				if (z == 0 || z > 0 && grid[x][z-1].cellType&wallBit == 0) && int(grid[x][z].walls[y][2]) > 0 {
					for i, v := range cubeLightSide {
						processVertex(v, i, coords, int(grid[x][z].walls[y][2])-1, ambient, 2)
					}
				}
				if (z == MAP_SIZE-1 || z < MAP_SIZE-1 && grid[x][z+1].cellType&wallBit == 0) && int(grid[x][z].walls[y][3]) > 0 {
					for i, v := range cubeDarkSide {
						processVertex(v, i, coords, int(grid[x][z].walls[y][3])-1, ambient, 3)
					}
				}

			}

		}
	}

	vertexMutex = false

}
