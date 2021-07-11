package main

import (
	"math"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

func renderWorld() {

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	position := mgl32.Vec3{float32(myX), float32(myY), float32(myZ)}
	focus := mgl32.Vec3{float32(myX + 10*math.Cos(bearing)*math.Cos(pitch)), float32(myY + 10*math.Sin(pitch)), float32(myZ + 10*math.Sin(bearing)*math.Cos(pitch))}
	up := mgl32.Vec3{0, 1, 0}
	camera := mgl32.LookAtV(position, focus, up)

	cameraUniform := gl.GetUniformLocation(shaderProgram, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	positionUniform := gl.GetUniformLocation(shaderProgram, gl.Str("position\x00"))
	gl.Uniform3f(positionUniform, position.X(), position.Y(), position.Z())

	for i := 0; i < textureCount; i++ {

		if len(vertices[i]) > 0 {

			gl.BufferData(gl.ARRAY_BUFFER, len(vertices[i])*4, gl.Ptr(vertices[i]), gl.STATIC_DRAW)

			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, texture[i])
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices[i])))

		}

	}

	window.SwapBuffers()

}
