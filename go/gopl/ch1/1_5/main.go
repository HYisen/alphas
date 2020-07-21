package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

var palette = []color.Color{color.Black, color.RGBA{R: 0x39, G: 0xff, B: 0x14, A: 0xff}} // neon green #39ff014

const (
	whiteIndex = 0
	blackIndex = 1
)

func main() {
	lissajous(os.Stdout)
}

func lissajous(out io.Writer) {
	const (
		cycles = 5
		res    = 0.001
		size   = 100
		frames = 64
		delay  = 8
	)

	freq := rand.Float64() * 30
	anim := gif.GIF{LoopCount: frames}
	phase := 0.0
	for i := 0; i < frames; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			//fmt.Fprintf(os.Stderr, "%d %d \n", size+int(x*size+0.5), size+int(y*size+0.5))
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		//fmt.Fprintf(os.Stderr, "phase=%v\n", phase)
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
