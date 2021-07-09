package main

import (
	"math"
)

const infinity float64 = 100000

func evaluateCursor() (int, int, int, int, int) {

	cursorX := -1
	cursorY := -1
	cursorZ := -1
	cursorWall := -1
	cursorTexture := -1

	dx := math.Sin(bearing) * math.Cos(pitch)
	dy := math.Sin(pitch)
	dz := math.Cos(bearing) * math.Cos(pitch)

	stepX := myX / unit
	stepY := myY / unit
	stepZ := myZ / unit

	for steps := 0; steps < 100; steps++ {

		toXedge := infinity
		toYedge := infinity
		toZedge := infinity

		_, fractX := math.Modf(stepX)
		_, fractY := math.Modf(stepY)
		_, fractZ := math.Modf(stepZ)

		if fractX == 0 {
			toXedge = 1 / math.Abs(dx)
		} else {
			if dx > 0 {
				toXedge = (1 - fractX) / dx
			} else if dx < 0 {
				toXedge = -fractX / dx
			}
		}

		if fractY == 0 {
			toYedge = 1 / math.Abs(dy)
		} else {
			if dy > 0 {
				toYedge = (1 - fractY) / dy
			} else if dy < 0 {
				toYedge = -fractY / dy
			}
		}

		if fractZ == 0 {
			toZedge = 1 / math.Abs(dz)
		} else {
			if dz > 0 {
				toZedge = (1 - fractZ) / dz
			} else if dz < 0 {
				toZedge = -fractZ / dz
			}
		}

		if toXedge == infinity && toYedge == infinity && toZedge == infinity {
			break
		}

		if toXedge <= toYedge && toXedge <= toZedge {

			stepX += dx * toXedge
			stepY += dy * toXedge
			stepZ += dz * toXedge

		} else if toYedge <= toXedge && toYedge <= toZedge {

			stepX += dx * toYedge
			stepY += dy * toYedge
			stepZ += dz * toYedge

		} else if toZedge <= toXedge && toZedge <= toYedge {

			stepX += dx * toZedge
			stepY += dy * toZedge
			stepZ += dz * toZedge

		}

		if stepX < 0 || stepX >= float64(MAP_SIZE) ||
			stepY < 0 || stepY >= float64(MAP_SIZE) {
			break
		}

		gridX := int(math.Floor(stepX + dx*0.001))
		gridY := int(math.Floor(stepY + dy*0.001))
		gridZ := int(math.Floor(stepZ + dz*0.001))

		INTERACTION_BIT := uint(0)
		if gridY <= 0 {
			INTERACTION_BIT = GROUND_BIT
		} else if gridY == 1 {
			INTERACTION_BIT = WALL_BIT
		} else if gridY == 2 {
			INTERACTION_BIT = LOW_BIT
		} else if gridY == 3 {
			INTERACTION_BIT = HIGH_BIT
		} else {
			INTERACTION_BIT = CEILING_BIT
		}

		if grid[gridX][gridZ].cellType&INTERACTION_BIT > 0 {
			cursorX = gridX
			cursorY = gridY
			cursorZ = gridZ

			_, fractX := math.Modf(stepX)
			_, fractY := math.Modf(stepY)
			_, fractZ := math.Modf(stepZ)

			if fractY != 0 {
				if fractX == 0 {
					if dx > 0 {
						cursorWall = 3
					} else if dx < 0 {
						cursorWall = 1
					}
				} else if fractZ == 0 {
					if dz > 0 {
						cursorWall = 0
					} else if dz < 0 {
						cursorWall = 2
					}
				}
			}

			if cursorWall != -1 && gridY >= 0 && gridY <= 2 {
				cursorTexture = grid[gridX][gridZ].walls[gridY][cursorWall]
			} else if gridY >= 0 && gridY <= 3 {
				cursorTexture = grid[gridX][gridZ].flats[gridY]
			}

			return cursorX, cursorY, cursorZ, cursorWall, cursorTexture

		}

	}

	return -1, -1, -1, -1, -1

}
