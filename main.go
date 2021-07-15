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
	tenth             = time.Tick(time.Millisecond * 100)
	frameLength       float64
	fpsString         = "Caclulating fps..."
	windowTitlePrefix = "OpenGL Maze Experiment"
	window            *glfw.Window
)

func main() {

	makeMaze()

	initiateOpenGL()
	prepareTextures()
	prepareShaders()

	for !window.ShouldClose() {

		frameStart := time.Now()

		processInputs()

		select {
		case <-tenth:
			go updateWorld()
		default:
		}

		if !vertexMutex {
			vertices = verticesTemp
			vertexRecords = vertexRecordsTemp
		}
		evaluateCursor()

		renderWorld()

		glfw.PollEvents()
		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			fpsString = fmt.Sprintf("%d FPS", frames)
			frames = 0
		default:
		}
		frameLength = time.Since(frameStart).Seconds()

	}

	glfw.Terminate()
}
