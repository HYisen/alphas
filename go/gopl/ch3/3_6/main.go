package main

import (
	"image"
	"image/png"
	"os"
)

func main() {
	fp, _ := os.Open("a.jpg")
	old, _ := png.Decode(fp)
	bounds := old.Bounds()
	neo := image.NewRGBA(image.Rect(0, 0, bounds.Dx()*2, bounds.Dy()*2))
	for i := 0; i < bounds.Dx(); i++ {
		for j := 0; j < bounds.Dy(); j++ {
			color := old.At(i, j)
			neo.Set(i*2, j*2, color)
			neo.Set(i*2+1, j*2, color)
			neo.Set(i*2, j*2+1, color)
			neo.Set(i*2+1, j*2+1, color)
		}
	}
	neoFp, _ := os.Create("b.png")
	_ = png.Encode(neoFp, neo)
	_ = neoFp.Close()
	_ = fp.Close()
}
