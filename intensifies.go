package yeller

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"

	"github.com/soniakeys/quant/median"
)

func createImagePalette(srcImage image.Image, maxColors int) color.Palette {
	bounds := srcImage.Bounds()
	quantizer := median.Quantizer(maxColors)

	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, srcImage, image.Point{}, draw.Src)
	palette := quantizer.Quantize(make(color.Palette, 0, maxColors - 1), img)
	palette = append(palette, color.RGBA{0, 0, 0, 0})
	return palette
}

func Intensifies(target image.Image) *bytes.Buffer {
	var frames []*image.Paletted
	var delays []int
	var disposal []byte

	target = scaleDown(target)
	bounds := target.Bounds()

	for i := 0; i < 15; i++ {
		frame := image.NewRGBA(bounds)
		offsetX := rand.Intn(30) - 15
		offsetY := rand.Intn(30) - 15

		draw.Draw(frame, bounds, &image.Alpha{}, image.Point{}, draw.Src)
		draw.Draw(frame, bounds, target, image.Point{offsetX, offsetY}, draw.Over)

		framePaletted := image.NewPaletted(bounds, createImagePalette(target, 256))
		draw.FloydSteinberg.Draw(framePaletted, bounds, frame, image.Point{})

		frames = append(frames, framePaletted)
		delays = append(delays, 5)
		disposal = append(disposal, byte(gif.DisposalPrevious))
	}

	buf := new(bytes.Buffer)

	enc := base64.NewEncoder(base64.StdEncoding, buf)

	if err := gif.EncodeAll(enc, &gif.GIF{
		Image: frames,
		Delay: delays,
		Disposal: disposal,
	}); err != nil {
		panic(err)
	}
	
	if err := enc.Close(); err != nil {
		panic(err)
	}

	return buf
}
