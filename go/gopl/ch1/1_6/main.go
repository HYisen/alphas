package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)

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
		step   = 16
	)

	var palette []color.Color

	for k := 0; k < step; k++ {
		grey := uint8(k * 256 / step)
		palette = append(palette, color.RGBA{
			R: grey,
			G: grey,
			B: grey,
			A: 0xff,
		})
	}

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
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), uint8(int((x*step+step)/2)))
		}
		//fmt.Fprintf(os.Stderr, "phase=%v\n", phase)
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	if err := gif.EncodeAll(out, &anim); err != nil {
		fmt.Println(err)
	}
}
