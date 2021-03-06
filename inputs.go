package main

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	myX             float64 = float64(MAP_CENTRE)
	myY             float64 = 0.5
	myZ             float64 = float64(MAP_CENTRE)
	pitch           float64 = 0
	bearing         float64 = 0
	unit                    = 1.0
	selectedTexture         = 1
	lastSpecial     bool
	lastNumberKey   bool
	drawDistance    float64 = 16
)

func dist(x0, y0, x1, y1 float64) float64 {
	return math.Sqrt(math.Pow(x0-x1, 2) + math.Pow(y0-y1, 2))
}

func processInputs() {

	potentialX := myX
	potentialZ := myZ

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		potentialX += 10 * frameLength * math.Cos(bearing) //* math.Cos(pitch)
		//myY += 10 * frameLength * math.Sin(pitch)
		potentialZ += 10 * frameLength * math.Sin(bearing) //* math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyS) == glfw.Press {
		potentialX -= 10 * frameLength * math.Cos(bearing) //* math.Cos(pitch)
		//myY -= 10 * frameLength * math.Sin(pitch)
		potentialZ -= 10 * frameLength * math.Sin(bearing) //* math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyA) == glfw.Press {
		potentialX += 10 * frameLength * math.Sin(bearing)
		potentialZ -= 10 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyD) == glfw.Press {
		potentialX -= 10 * frameLength * math.Sin(bearing)
		potentialZ += 10 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyF) == glfw.Press {
		//potentialX += 10 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY -= 10 * frameLength //* math.Cos(pitch)
		//potentialZ += 10 * frameLength * math.Sin(bearing) * math.Sin(pitch)

	}

	if window.GetKey(glfw.KeyR) == glfw.Press {
		//potentialX -= 10 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY += 10 * frameLength //* math.Cos(pitch)
		//potentialZ -= 10 * frameLength * math.Sin(bearing) * math.Sin(pitch)
	}

	if window.GetKey(glfw.KeySpace) == glfw.Press {
		if cursorX >= 0 && cursorY >= 0 && cursorZ >= 0 &&
			cursorX < MAP_SIZE && cursorY < MAP_HEIGHT && cursorZ < MAP_SIZE {
			if cursorWall == -1 {
				grid[cursorX][cursorZ].flats[cursorY] = selectedTexture
			} else if cursorWall >= 0 && cursorWall < 4 && cursorY < len(grid[cursorX][cursorZ].walls) {
				grid[cursorX][cursorZ].walls[cursorY][cursorWall] = selectedTexture
			}
		}
	}

	if window.GetKey(glfw.KeyQ) == glfw.Press {
		if cursorTexture >= 0 && cursorTexture < textureCount {
			selectedTexture = cursorTexture
		}
	}

	if window.GetKey(glfw.KeyRightBracket) == glfw.Press {
		if !lastSpecial {
			selectedTexture++
			if selectedTexture > textureCount {
				selectedTexture = 1
			}
			lastSpecial = true
		}
	} else if window.GetKey(glfw.KeyLeftBracket) == glfw.Press {
		if !lastSpecial {
			selectedTexture--
			if selectedTexture < 1 {
				selectedTexture = textureCount
			}
			lastSpecial = true
		}
	} else if window.GetKey(glfw.KeyPageUp) == glfw.Press {
		if !lastSpecial {
			drawDistance++
			if drawDistance > float64(MAP_SIZE) {
				drawDistance = float64(MAP_SIZE)
			}
			lastSpecial = true
		}
	} else if window.GetKey(glfw.KeyPageDown) == glfw.Press {
		if !lastSpecial {
			drawDistance--
			if drawDistance < 1 {
				drawDistance = 1
			}
			lastSpecial = true
		}
	} else {
		lastSpecial = false
	}

	if cursorX >= 0 && cursorY >= 0 && cursorZ >= 0 &&
		cursorX < MAP_SIZE && cursorY < MAP_HEIGHT && cursorZ < MAP_SIZE {
		if window.GetKey(glfw.Key1) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = WALL
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key2) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = CORRIDOR
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key3) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = LOW_ROOM
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key4) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = HIGH_ROOM
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key5) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = HIGH_ROOM_SINGLE_BLOCK
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key6) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = HIGH_ROOM_DOUBLE_BLOCK
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key8) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = SKY_DOUBLE_BLOCK
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key9) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = SKY_SINGLE_BLOCK
				lastNumberKey = true
			}
		} else if window.GetKey(glfw.Key0) == glfw.Press {
			if !lastNumberKey {
				grid[cursorX][cursorZ].cellType = SKY
				lastNumberKey = true
			}
		} else {
			lastNumberKey = false
		}
	}

	if myY < 0.5 {
		myY = 0.5
	}

	if potentialX < 0 {
		potentialX = 0
	}
	if potentialZ < 0 {
		potentialZ = 0
	}
	if potentialX > float64(MAP_SIZE)-1 {
		potentialX = float64(MAP_SIZE) - 1
	}
	if potentialZ > float64(MAP_SIZE)-1 {
		potentialZ = float64(MAP_SIZE) - 1
	}

	distanceTravelled := dist(potentialX, potentialZ, myX, myZ)
	if distanceTravelled > unit/2 {
		dx := (potentialX - myX) / distanceTravelled
		dz := (potentialZ - myZ) / distanceTravelled
		potentialX = myX + dx*(unit/2)
		potentialZ = myZ + dz*(unit/2)
	}

	mouseX, mouseY := window.GetCursorPos()

	bearing += (mouseX - windowWidth/2) * 0.005
	pitch += (windowHeight/2 - mouseY) * 0.005

	window.SetCursorPos(windowWidth/2, windowHeight/2)

	if bearing > math.Pi {
		bearing -= 2 * math.Pi
	}
	if bearing < -math.Pi {
		bearing += 2 * math.Pi
	}
	if pitch > 0.5*math.Pi-0.001 {
		pitch = 0.5*math.Pi - 0.001
	}
	if pitch < -0.5*math.Pi+0.001 {
		pitch = -0.5*math.Pi + 0.001
	}

	mapXstart := int(math.Floor(myX / unit))
	mapZstart := int(math.Floor(myZ / unit))

	mapXtarget := int(math.Floor(potentialX / unit))
	mapZtarget := int(math.Floor(potentialZ / unit))

	lowestX := int(math.Floor(math.Min(float64(mapXstart), float64(mapXtarget))) - 1)
	highestX := int(math.Floor(math.Max(float64(mapXstart), float64(mapXtarget))) + 1)
	lowestZ := int(math.Floor(math.Min(float64(mapZstart), float64(mapZtarget))) - 1)
	highestZ := int(math.Floor(math.Max(float64(mapZstart), float64(mapZtarget))) + 1)

	interactionBit := uint(0)
	switch math.Floor(myY) {
	case 0:
		interactionBit = WALL_BIT
	case 1:
		interactionBit = LOW_BIT
	case 2:
		interactionBit = HIGH_BIT
	}

	if interactionBit != 0 &&
		mapXstart >= 0 && mapZstart >= 0 &&
		mapXstart < MAP_SIZE && mapZstart < MAP_SIZE {

		for i := lowestX; i <= highestX; i++ {
			for j := lowestZ; j <= highestZ; j++ {

				if i >= 0 && j >= 0 &&
					i < MAP_SIZE && j < MAP_SIZE {

					if grid[i][j].cellType&interactionBit > 0 {

						nearestPointX := math.Max(float64(i)*unit, math.Min(potentialX, float64(i+1)*unit))
						nearestPointZ := math.Max(float64(j)*unit, math.Min(potentialZ, float64(j+1)*unit))

						rayX := nearestPointX - potentialX
						rayZ := nearestPointZ - potentialZ

						rayLength := dist(0, 0, rayX, rayZ)
						if rayLength == 0 {
							continue
						}

						rayOverlap := unit*0.5 - rayLength

						if rayOverlap > 0 {

							potentialX -= rayOverlap * rayX / rayLength
							potentialZ -= rayOverlap * rayZ / rayLength

						}

					}

				}
			}

		}

	}

	myX = potentialX
	myZ = potentialZ

}
