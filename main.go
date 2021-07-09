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

		prepareVertices()
		renderWorld()

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			fmt.Printf("FPS: %d\tPlayer x: %v, y: %v, z: %v\tCursor x: %v, y: %v, z: %v, w: %v, t: %v\n",
				frames, myX, myY, myZ,
				cursorX, cursorY, cursorZ, cursorWall, cursorTexture)
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}
