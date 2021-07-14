package main

import (
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

	lineVertices := []float32{
		0, -30, 0, -15,
		-30, 0, -15, 0,
		0, 30, 0, 15,
		30, 0, 15, 0,
	}
	lineColour := []float32{0, 1, 0}

	lineColourUniform := gl.GetUniformLocation(lineShaderProgram, gl.Str("colour"+terminator))
	gl.Uniform3f(lineColourUniform, lineColour[0], lineColour[1], lineColour[2])

	gl.BufferData(gl.ARRAY_BUFFER, len(lineVertices)*4, gl.Ptr(lineVertices), gl.STATIC_DRAW)

	gl.DrawArrays(gl.LINES, 0, int32(len(lineVertices))/2)

	gl.Enable(gl.DEPTH_TEST)

	window.SwapBuffers()

}
