package main

import (
	"fmt"
	"math"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

func renderWorld() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(triangleShaderProgram)
	gl.BindVertexArray(triangleVertexArray)
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexBuffer)

	position := mgl32.Vec3{float32(myX), float32(myY), float32(myZ)}
	focus := mgl32.Vec3{
		float32(myX + 10*math.Cos(bearing)*math.Cos(pitch)),
		float32(myY + 10*math.Sin(pitch)),
		float32(myZ + 10*math.Sin(bearing)*math.Cos(pitch)),
	}
	up := mgl32.Vec3{0, 1, 0}
	camera := mgl32.LookAtV(position, focus, up)

	cameraUniform := gl.GetUniformLocation(triangleShaderProgram, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	positionUniform := gl.GetUniformLocation(triangleShaderProgram, gl.Str("position\x00"))
	gl.Uniform3f(positionUniform, position.X(), position.Y(), position.Z())

	drawDistanceUniform := gl.GetUniformLocation(triangleShaderProgram, gl.Str("drawDistance\x00"))
	gl.Uniform1f(drawDistanceUniform, float32(drawDistance))

	for i := 0; i < len(vertices); i++ {

		if len(vertices[i]) > 0 {

			gl.BufferData(gl.ARRAY_BUFFER, len(vertices[i])*4, gl.Ptr(vertices[i]), gl.STATIC_DRAW)

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, texture[i])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices[i]))/8)

		}

	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.CONSTANT_COLOR, gl.ONE_MINUS_CONSTANT_COLOR)
	gl.BlendColor(0.5, 0.5, 0.5, 1.0)

	if len(cursorVertices) > 0 {

		gl.BufferData(gl.ARRAY_BUFFER, len(cursorVertices)*4, gl.Ptr(cursorVertices), gl.STATIC_DRAW)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture[selectedTexture-1])
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cursorVertices))/8)

	}

	gl.Disable(gl.BLEND)

	gl.UseProgram(lineShaderProgram)
	gl.BindVertexArray(lineVertexArray)
	gl.BindBuffer(gl.ARRAY_BUFFER, lineVertexBuffer)

	lineColourUniform := gl.GetUniformLocation(lineShaderProgram, gl.Str("colour"+terminator))

	gl.Uniform3f(lineColourUniform, 0, 1, 0)

	cursorLineVertices := []float32{
		0, -30, 0, -15,
		-30, 0, -15, 0,
		0, 30, 0, 15,
		30, 0, 15, 0,
	}
	gl.BufferData(gl.ARRAY_BUFFER, len(cursorLineVertices)*4, gl.Ptr(cursorLineVertices), gl.STATIC_DRAW)
	gl.DrawArrays(gl.LINES, 0, int32(len(cursorLineVertices))/2)

	gl.Uniform3f(lineColourUniform, 1, 1, 1)

	drawString(fpsString, -windowWidth/2+20, windowHeight/2-20, 10)
	drawString(fmt.Sprintf("player %.1f %.1f %.1f  %.2f mark %.2f pi", myX, myY, myZ, bearing/math.Pi, pitch/math.Pi), -windowWidth/2+20, windowHeight/2-50, 10)
	drawString(fmt.Sprintf("cursor %v %v %v  wall %v tex %v", cursorX, cursorY, cursorZ, cursorWall, cursorTexture), -windowWidth/2+20, windowHeight/2-80, 10)
	drawString(fmt.Sprintf("texture %v", selectedTexture), -windowWidth/2+20, windowHeight/2-110, 10)
	drawString(fmt.Sprintf("distance %v", drawDistance), -windowWidth/2+20, windowHeight/2-140, 10)

	gl.Enable(gl.DEPTH_TEST)

	window.SwapBuffers()

}
