package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

func main() {
	const (
		xMin, yMin, xMax, yMax = -2, -2, 2, 2
		width, height          = 1024, 1024
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin
			z := complex(x, y)
			img.Set(px, py, mandelbrot(z))
		}
	}

	f, _ := os.Create("a.jpg")
	_ = png.Encode(f, img)
	_ = f.Close()
}

func max(l int16, r uint8) uint8 {
	if l > int16(r) {
		return uint8(l)
	}
	return r
}

func mandelbrot(z complex128) color.Color {
	const iterations = 200
	const contrast = 10

	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			rgba := color.RGBA{
				R: max(int16(255*1)-int16(contrast)*3*int16(n), 0),
				G: max(int16(255*2)-int16(contrast)*3*int16(n), 0),
				B: max(int16(255*3)-int16(contrast)*3*int16(n), 0),
				A: 255,
			}
			//fmt.Println(rgba)
			return rgba
		}
	}
	return color.Black

}
