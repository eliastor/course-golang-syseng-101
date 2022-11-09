package main

import (
	"image"
	"image/color"

	"golang.org/x/tour/pic"
)

type Image struct{}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 256, 256)
}

func (i Image) At(x, y int) color.Color {
	var r, g, b, a uint8 = 25, 0, 0, 0
	b = uint8(x) * uint8(y)
	return color.RGBA{r, g, b, a}
}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
