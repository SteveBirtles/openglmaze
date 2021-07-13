package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 1280
const windowHeight = 720
const fullScreen = false

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "OpenGL Maze Experiment"
	window            *glfw.Window
	cursorX           int
	cursorY           int
	cursorZ           int
	cursorTexture     int
	cursorWall        int
	cursorIndex       int
	cursorLastR       float32
	cursorLastG       float32
	cursorLastB       float32
)

func main() {

	makeMaze()

	initiateOpenGL()
	prepareTextures()
	prepareShaders()

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs()
		cursorX, cursorY, cursorZ, cursorWall, cursorTexture = evaluateCursor()

		for i, record := range vertexRecords {
			if record.x == cursorX && record.y == cursorY && record.z == cursorZ && record.wall == cursorWall {

				cursorIndex = i
				cursorLastR = vertices[record.texture][record.index+4]
				cursorLastG = vertices[record.texture][record.index+5]
				cursorLastB = vertices[record.texture][record.index+6]

				vertices[record.texture][record.index+4] = 1.0
				vertices[record.texture][record.index+5] = 1.0
				vertices[record.texture][record.index+6] = 1.0

				break
			}
		}

		renderWorld()

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			updateWorld()
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			fmt.Printf("FPS: %d\tPlayer x: %v, y: %v, z: %v\n", frames, myX, myY, myZ)
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}
