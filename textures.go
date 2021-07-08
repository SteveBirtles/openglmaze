package main

import (
	"fmt"
	"image"
	"image/draw"
	"log"
	"math"
	"os"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

const textureCount = 40

var texture [textureCount]uint32

func prepareTextures() {

	var err error
	for i := 0; i < textureCount; i++ {
		texture[i], err = newTexture(fmt.Sprintf("textures/%v.png", i+1))
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func newTexture(file string) (uint32, error) {

	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)

	var texID uint32
	gl.GenTextures(1, &texID)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texID)
	gl.TexStorage2D(gl.TEXTURE_2D, int32(math.Log2(math.Min(float64(width), float64(height)))), gl.RGBA8, width, height)
	gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, width, height, gl.BGRA_EXT, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		width,
		height,
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texID, nil

}
