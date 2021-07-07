package main

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	myX     float64 = float64(MAP_CENTRE)
	myY     float64 = 0
	myZ     float64 = float64(MAP_CENTRE)
	lastX   float64 = myX
	lastZ   float64 = myZ
	pitch   float64 = 0
	bearing float64 = 0
	unit            = 1.0
)

func dist(x0, y0, x1, y1 float64) float64 {
	return math.Sqrt(math.Pow(x0-x1, 2) + math.Pow(y0-y1, 2))
}

func processInputs() {

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

	if window.GetKey(glfw.KeyW) == glfw.Press {
		myX += 10 * frameLength * math.Cos(bearing) //* math.Cos(pitch)
		//myY += 10 * frameLength * math.Sin(pitch)
		myZ += 10 * frameLength * math.Sin(bearing) //* math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyS) == glfw.Press {
		myX -= 10 * frameLength * math.Cos(bearing) //* math.Cos(pitch)
		//myY -= 10 * frameLength * math.Sin(pitch)
		myZ -= 10 * frameLength * math.Sin(bearing) //* math.Cos(pitch)
	}

	if window.GetKey(glfw.KeyA) == glfw.Press {
		myX += 10 * frameLength * math.Sin(bearing)
		myZ -= 10 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyD) == glfw.Press {
		myX -= 10 * frameLength * math.Sin(bearing)
		myZ += 10 * frameLength * math.Cos(bearing)
	}

	if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
		//myX += 10 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY -= 10 * frameLength //* math.Cos(pitch)
		//myZ += 10 * frameLength * math.Sin(bearing) * math.Sin(pitch)

	}

	if window.GetKey(glfw.KeySpace) == glfw.Press {
		//myX -= 10 * frameLength * math.Cos(bearing) * math.Sin(pitch)
		myY += 10 * frameLength //* math.Cos(pitch)
		//myZ -= 10 * frameLength * math.Sin(bearing) * math.Sin(pitch)
	}

	distanceTravelled := dist(myX, myZ, lastX, lastZ)
	if distanceTravelled > unit/2 {
		dx := (myX - lastX) / distanceTravelled
		dy := (myZ - lastZ) / distanceTravelled
		myX = lastX + dx*(unit/2)
		myZ = lastZ + dy*(unit/2)
	}

	lastX = myX
	lastZ = myZ

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

	if myY < 0 {
		myY = 0
	}

	if myX < 0 {
		myX = 0
	}
	if myZ < 0 {
		myZ = 0
	}
	if myX > float64(MAP_SIZE) {
		myX = float64(MAP_SIZE)
	}
	if myZ > float64(MAP_SIZE) {
		myZ = float64(MAP_SIZE)
	}

	mapXstart := int(math.Floor(lastX / unit))
	mapZstart := int(math.Floor(lastZ / unit))

	mapXtarget := int(math.Floor(myX / unit))
	mapZtarget := int(math.Floor(myZ / unit))

	lowestX := int(math.Floor(math.Min(float64(mapXstart), float64(mapXtarget))) - 1)
	highestX := int(math.Floor(math.Max(float64(mapXstart), float64(mapXtarget))) + 1)
	lowestZ := int(math.Floor(math.Min(float64(mapZstart), float64(mapZtarget))) - 1)
	highestZ := int(math.Floor(math.Max(float64(mapZstart), float64(mapZtarget))) + 1)

	mapY := math.Floor(myY)

	interactionBit := uint(0)
	switch mapY {
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

						nearestPointX := math.Max(float64(i)*unit, math.Min(myX, float64((i)+1)*unit))
						nearestPointZ := math.Max(float64(j)*unit, math.Min(myZ, float64((j)+1)*unit))

						rayX := nearestPointX - myX
						rayZ := nearestPointZ - myZ

						rayLength := dist(0, 0, rayX, rayZ)
						if rayLength == 0 {
							continue
						}

						rayOverlap := unit*0.5 - rayLength

						if rayOverlap > 0 {

							myX -= rayOverlap * rayX / rayLength
							myZ -= rayOverlap * rayZ / rayLength

						}

					}

				}
			}

		}

	}
}
