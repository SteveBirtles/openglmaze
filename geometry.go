package main

import (
	"math"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

var (
	vertices []float32

	cubeBottom = []float32{
		0.5, -0.5, -0.5, 1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
	}
	cubeFlippedBottom = []float32{
		0.5, -0.5, -0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
	}
	cubeTop = []float32{
		-0.5, 0.5, -0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
	}
	cubeDarkSide = []float32{
		-0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, 0.5, 0.5, 1.0, 1.0,
	}
	cubeLightSide = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
	}
	cubeLeft = []float32{
		-0.5, -0.5, 0.5, 0.0, 1.0,
		-0.5, 0.5, -0.5, 1.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
		-0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 1.0, 0.0,
	}
	cubeRight = []float32{
		0.5, -0.5, 0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 0.0, 1.0,
	}
)

func processVertex(v float32, i, x, y, z, b int, rgb []float32) {

	if i%5 == 0 {
		v += float32(x)
	} else if i%5 == 1 {
		v += float32(y)
	} else if i%5 == 2 {
		v += float32(z)
	} else if i%5 == 3 {
		v = (v + float32(b%8)) / 8
	} else if i%5 == 4 {
		v = float32(int(v+float32(b/8))) / 5
	}
	vertices = append(vertices, v)
	if i%5 == 4 {
		vertices = append(vertices, rgb...)
	}

}

func prepareVertices() {

	const drawDistance float64 = 20

	vertices = make([]float32, 0)

	for x := int(math.Floor(myX) - drawDistance); x < int(math.Floor(myX)+drawDistance); x++ {
		for z := int(math.Floor(myZ) - drawDistance); z < int(math.Floor(myZ)+drawDistance); z++ {

			if x < 0 || z < 0 || x >= MAP_SIZE || z >= MAP_SIZE {
				continue
			}

			for y := 0; y < MAP_HEIGHT; y++ {

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

				d := dist(float64(x), float64(z), myX, myZ)
				illumination := float32(math.Min(1, 1-d/drawDistance))

				ambient := []float32{illumination, illumination, illumination}

				flatTexture := int(grid[x][z].flats[y]) - 1
				wallTexture := -1
				if y < MAP_HEIGHT-1 {
					wallTexture = int(grid[x][z].walls[y]) - 1
				}

				if flatTexture != -1 && grid[x][z].cellType&flatBit > 0 {
					if y < MAP_HEIGHT-1 && grid[x][z].cellType&wallBit == 0 {
						for i, v := range cubeBottom {
							processVertex(v, i, x, y, z, flatTexture, []float32{0.0, 0.0, ambient[2]})
						}
					}
				}

				if flatTexture != -1 && grid[x][z].cellType&flatBit == 0 && grid[x][z].cellType&wallBit > 0 {
					for i, v := range cubeFlippedBottom {
						processVertex(v, i, x, y, z, flatTexture, []float32{ambient[0], 0.0, ambient[2]})
					}
				}

				if wallTexture == -1 || grid[x][z].cellType&wallBit == 0 {
					continue
				}

				if x == 0 || x > 0 && grid[x-1][z].cellType&wallBit == 0 {
					for i, v := range cubeLeft {
						processVertex(v, i, x, y, z, wallTexture, []float32{ambient[0], 0.0, 0.0})
					}
				}
				if x == MAP_SIZE-2 || x < MAP_SIZE-1 && grid[x+1][z].cellType&wallBit == 0 {
					for i, v := range cubeRight {
						processVertex(v, i, x, y, z, wallTexture, []float32{ambient[0], ambient[1], 0.0})
					}
				}
				if z == 0 || z > 0 && grid[x][z-1].cellType&wallBit == 0 {
					for i, v := range cubeLightSide {
						processVertex(v, i, x, y, z, wallTexture, []float32{0.0, ambient[1], 0.0})
					}
				}
				if z == MAP_SIZE-2 || z < MAP_SIZE-1 && grid[x][z+1].cellType&wallBit == 0 {
					for i, v := range cubeDarkSide {
						processVertex(v, i, x, y, z, wallTexture, []float32{0.0, ambient[1], ambient[2]})
					}
				}

			}

		}
	}

	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

}
