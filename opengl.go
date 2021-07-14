package main

import (
	"log"
	"runtime"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	triangleVertexBuffer uint32
	triangleVertexArray  uint32
	lineVertexBuffer     uint32
	lineVertexArray      uint32
)

func init() {

	runtime.LockOSThread()

}

func initiateOpenGL() {

	var err error
	if err = glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 4)

	if fullScreen {
		window, err = glfw.CreateWindow(windowWidth, windowHeight, windowTitlePrefix, glfw.GetPrimaryMonitor(), nil)
	} else {
		window, err = glfw.CreateWindow(windowWidth, windowHeight, windowTitlePrefix, nil, nil)
	}
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPos(windowWidth/2, windowHeight/2)

	if glfw.RawMouseMotionSupported() {
		window.SetInputMode(glfw.RawMouseMotion, glfw.True)
	}

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	gl.GenVertexArrays(1, &triangleVertexArray)
	gl.GenBuffers(1, &triangleVertexBuffer)

	gl.GenVertexArrays(1, &lineVertexArray)
	gl.GenBuffers(1, &lineVertexBuffer)

}
