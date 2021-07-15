package main

import (
	"fmt"
	_ "image/png"
	"strings"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

const terminator = "\x00"

var triangleShaderProgram uint32
var lineShaderProgram uint32

const triangleVertexShader = `#version 300 es

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
uniform vec3 position;
uniform float drawDistance;

in vec3 vertexIn;
in vec2 texCoordIn;
in vec3 colourIn;

out mediump vec2 texCoordForFrag;
out mediump vec3 colourForFrag;

void main() {
	
	float illumination = min(0.666667, 1.0 - length(vertexIn - position) / drawDistance);

    texCoordForFrag = texCoordIn;
    colourForFrag = colourIn * illumination;

    gl_Position = projection * camera * model * vec4(vertexIn, 1);

}`

const triangleFragmentShader = `#version 300 es

uniform sampler2D tex;

in mediump vec2 texCoordForFrag;
in mediump vec3 colourForFrag;

out mediump vec4 colourOut;

void main() {

    colourOut = texture(tex, texCoordForFrag);

    colourOut.x *= colourForFrag.x;
    colourOut.y *= colourForFrag.y;
    colourOut.z *= colourForFrag.z;

}`

const lineVertexShader = `#version 300 es

uniform mat4 projection;
uniform mat4 matrix;

in vec2 vertex;

void main() {

	gl_Position = projection * matrix * vec4(vertex.x, vertex.y, 0.0, 1.0);

}`

const lineFragmentShader = `#version 300 es

uniform mediump vec3 colour;

out mediump vec4 FragColour;

void main() {
	
   FragColour = vec4(colour, 1.0f);

}`

func prepareShaders() {

	var err error

	triangleShaderProgram, err = newShaderProgram(triangleVertexShader, triangleFragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(triangleShaderProgram)
	gl.BindVertexArray(triangleVertexArray)
	gl.BindBuffer(gl.ARRAY_BUFFER, triangleVertexBuffer)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 5000.0)
	projectionUniform := gl.GetUniformLocation(triangleShaderProgram, gl.Str("projection"+terminator))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(triangleShaderProgram, gl.Str("model"+terminator))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	tex := gl.GetUniformLocation(triangleShaderProgram, gl.Str("tex"+terminator))
	gl.Uniform1i(tex, 0)

	vertexIn := uint32(gl.GetAttribLocation(triangleShaderProgram, gl.Str("vertexIn"+terminator)))
	gl.EnableVertexAttribArray(vertexIn)
	gl.VertexAttribPointer(vertexIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	texCoordIn := uint32(gl.GetAttribLocation(triangleShaderProgram, gl.Str("texCoordIn"+terminator)))
	gl.EnableVertexAttribArray(texCoordIn)
	gl.VertexAttribPointer(texCoordIn, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	colourIn := uint32(gl.GetAttribLocation(triangleShaderProgram, gl.Str("colourIn"+terminator)))
	gl.EnableVertexAttribArray(colourIn)
	gl.VertexAttribPointer(colourIn, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(5*4))

	lineShaderProgram, err = newShaderProgram(lineVertexShader, lineFragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(lineShaderProgram)
	gl.BindVertexArray(lineVertexArray)
	gl.BindBuffer(gl.ARRAY_BUFFER, lineVertexBuffer)

	orthoProjection := mgl32.Ortho(-float32(windowWidth)/2, float32(windowWidth)/2, -float32(windowHeight)/2, float32(windowHeight)/2, 0, 100)
	orthoProjectionUniform := gl.GetUniformLocation(lineShaderProgram, gl.Str("projection"+terminator))
	gl.UniformMatrix4fv(orthoProjectionUniform, 1, false, &orthoProjection[0])

	position := mgl32.Vec3{0, 0, 0}
	focus := mgl32.Vec3{0, 0, -100}
	up := mgl32.Vec3{0, 100, 0}
	orthoCamera := mgl32.LookAtV(position, focus, up)

	matrixUniform := gl.GetUniformLocation(lineShaderProgram, gl.Str("matrix"+terminator))
	gl.UniformMatrix4fv(matrixUniform, 1, false, &orthoCamera[0])

	lineVertexIn := uint32(gl.GetAttribLocation(lineShaderProgram, gl.Str("vertex"+terminator)))
	gl.EnableVertexAttribArray(lineVertexIn)
	gl.VertexAttribPointer(lineVertexIn, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

}

func newShaderProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {

	vertexShader, err := compileShader(vertexShaderSource+terminator, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource+terminator, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	shaderProgram := gl.CreateProgram()

	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)

	var status int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat(terminator, int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link triangle shader program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram, nil

}

func compileShader(source string, shaderType uint32) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat(terminator, int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil

}
