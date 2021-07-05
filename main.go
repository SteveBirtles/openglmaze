package main

import (
	"fmt"
	_ "image/png"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const windowWidth = 1920
const windowHeight = 1080

var (
	frames            = 0
	second            = time.Tick(time.Second)
	frameLength       float64
	windowTitlePrefix = "OpenGL Maze Experiment"
	window            *glfw.Window
)

func main() {

	loadMap()

	initiateOpenGL()
	prepareVertices()
	prepareTextures()
	prepareShaders()

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs()
		renderWorld()

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}
