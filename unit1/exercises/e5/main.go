package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

type Image struct {
	w int
	h int
	v uint8
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	var w, h = 255, 255
	return image.Rect(0, 0, w, h)
}

func (i Image) At(x, y int) color.Color {
	var v uint8 = 220
	return color.RGBA{v, v, 255, 255}
}

func main() {

	m := Image{}
	pic.ShowImage(m)
}
