package main

import (
	"math"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

var (
	vertices []float32

	cubeBottom = []float32{
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
	}
	cubeFlippedBottom = []float32{
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
	}
	cubeTop = []float32{
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,
	}
	cubeDarkSide = []float32{
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
	}
	cubeLightSide = []float32{
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,
	}
	cubeLeft = []float32{
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
	}
	cubeRight = []float32{
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}

	shadow []float32
)

func processVertex(v float32, i, x, y, z int, s bool, b int, rgb []float32) {

	if i%5 == 0 {
		v += float32(2 * x)
	} else if i%5 == 1 {
		v += float32(2 * y)
	} else if i%5 == 2 {
		v += float32(2 * z)
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

	for x := -MAP_CENTRE; x <= MAP_CENTRE; x++ {
		for z := -MAP_CENTRE; z <= MAP_CENTRE; z++ {
			for y := 0; y < MAP_HEIGHT; y++ {

				ambient := []float32{float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y))),
					float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y))),
					float32(32 / math.Hypot(math.Hypot(float64(x), float64(z)), float64(32-y)))}

				flatTexture := int(grid[x+MAP_CENTRE][z+MAP_CENTRE].flats[y]) - 1
				wallTexture := -1
				if y < MAP_HEIGHT-1 {
					wallTexture = int(grid[x+MAP_CENTRE][z+MAP_CENTRE].walls[y]) - 1
				}
				shadow = []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]}

				inShadow := calculateMapShadow(float64(x), float64(y), float64(z), uint16(wallTexture+1))

				if flatTexture != -1 {

					/*if wallTexture == -1 {
						if y == 0 || y > 0 && grid[x+MAP_CENTRE][z+MAP_CENTRE].flats[y-1] == 0 {
							for i, v := range cubeBottom {
								processVertex(v, i, x, y, z, inShadow, flatTexture, []float32{1 * ambient[0], 1 * ambient[1], 1 * ambient[2]})
							}
						}
						if y == 0 || y > 0 && grid[x+MAP_CENTRE][z+MAP_CENTRE].flats[y-1] == 0 {
							for i, v := range cubeFlippedBottom {
								processVertex(v, i, x, y, z, false, flatTexture, []float32{0.333 * ambient[0], 0.333 * ambient[1], 0.333 * ambient[2]})
							}
						}
					} else {*/
					if y == MAP_HEIGHT-1 || y < MAP_HEIGHT-2 && grid[x+MAP_CENTRE][z+MAP_CENTRE].flats[y+1] == 0 {
						for i, v := range cubeBottom {
							processVertex(v, i, x, y, z, inShadow, flatTexture, []float32{1 * ambient[0], 1 * ambient[1], 1 * ambient[2]})
						}
					}
					/*if y == 0 || y > 0 && grid[x+MAP_CENTRE][z+MAP_CENTRE].flats[y-1] == 0 {
						for i, v := range cubeFlippedBottom {
							processVertex(v, i, x, y, z, false, flatTexture, []float32{0.333 * ambient[0], 0.333 * ambient[1], 0.333 * ambient[2]})
						}
					}*/
					//}

				}

				if wallTexture == -1 {
					continue
				}

				if x == -MAP_CENTRE || x > -MAP_CENTRE && grid[x+MAP_CENTRE-1][z+MAP_CENTRE].walls[y] == 0 {
					for i, v := range cubeLeft {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]})
					}
				}
				if x == MAP_CENTRE || x < MAP_CENTRE && grid[x+MAP_CENTRE+1][z+MAP_CENTRE].walls[y] == 0 {
					for i, v := range cubeRight {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{0.5 * ambient[0], 0.5 * ambient[1], 0.5 * ambient[2]})
					}
				}
				if z == -MAP_CENTRE || z > -MAP_CENTRE && grid[x+MAP_CENTRE][z+MAP_CENTRE-1].walls[y] == 0 {
					for i, v := range cubeLightSide {
						processVertex(v, i, x, y, z, inShadow, wallTexture, []float32{0.75 * ambient[0], 0.75 * ambient[1], 0.75 * ambient[2]})
					}
				}
				if z == MAP_CENTRE || z < MAP_CENTRE && grid[x+MAP_CENTRE][z+MAP_CENTRE+1].walls[y] == 0 {
					for i, v := range cubeDarkSide {
						processVertex(v, i, x, y, z, false, wallTexture, []float32{0.333 * ambient[0], 0.333 * ambient[1], 0.333 * ambient[2]})
					}
				}

			}

		}
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

}
