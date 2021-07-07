package main

import (
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

	shadow []float32
)

func processVertex(v float32, i, x, y, z int, s bool, b int, rgb []float32) {

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
		if s {
			rgb = shadow
		}
		vertices = append(vertices, rgb...)
	}

}

func prepareVertices() {

	vertices = make([]float32, 0)

	for x := 0; x < MAP_SIZE; x++ {
		for z := 0; z < MAP_SIZE; z++ {
			for y := 0; y < MAP_HEIGHT; y++ {

				d := dist(float64(x), float64(z), myX, myZ)

				if d > 10 {
					continue
				}

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

				//d := float32(16 / math.Hypot(math.Hypot(myX-float64(x), myY-float64(y)), myZ-float64(z)))

				//ambient := []float32{d, d, d}

				ambient := []float32{1, 1, 1}

				flatTexture := int(grid[x][z].flats[y]) - 1
				wallTexture := -1
				if y < MAP_HEIGHT-1 {
					wallTexture = int(grid[x][z].walls[y]) - 1
				}
				shadow = []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]}

				inShadow := calculateMapShadow(float64(x), float64(y), float64(z), uint16(wallTexture+1))

				if flatTexture != -1 && grid[x][z].cellType&flatBit > 0 {
					if y < MAP_HEIGHT-1 && grid[x][z].cellType&wallBit == 0 {
						for i, v := range cubeBottom {
							processVertex(v, i, x, y, z, inShadow, flatTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
						}
					}
				}

				if flatTexture != -1 && grid[x][z].cellType&flatBit == 0 && grid[x][z].cellType&wallBit > 0 {
					for i, v := range cubeFlippedBottom {
						processVertex(v, i, x, y, z, false, flatTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
					}
				}

				if wallTexture == -1 || grid[x][z].cellType&wallBit == 0 {
					continue
				}

				if x == 0 || x > 0 && grid[x-1][z].cellType&wallBit == 0 {
					for i, v := range cubeLeft {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
					}
				}
				if x == MAP_SIZE-2 || x < MAP_SIZE-1 && grid[x+1][z].cellType&wallBit == 0 {
					for i, v := range cubeRight {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
					}
				}
				if z == 0 || z > 0 && grid[x][z-1].cellType&wallBit == 0 {
					for i, v := range cubeLightSide {
						processVertex(v, i, x, y, z, inShadow, wallTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
					}
				}
				if z == MAP_SIZE-2 || z < MAP_SIZE-1 && grid[x][z+1].cellType&wallBit == 0 {
					for i, v := range cubeDarkSide {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{1.0 * ambient[0], 1.0 * ambient[1], 1.0 * ambient[2]})
					}
				}

			}

		}
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

}
