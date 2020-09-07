package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"sync"
	"time"
)

func main() {
	exec(true)
}

// mode
// 0 singular
// 1 split in outer loop, width goroutines
// 2 split in both outer and inner loop, width+width*height goroutine
// 3 split in outer loop, width tasks in threads size queue goroutines
func exec(print bool) [32]byte {
	start := time.Now()

	const (
		xMin, yMin, xMax, yMax = -2, -2, 2, 2
		width, height          = 2048, 2048
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	wg := sync.WaitGroup{}
	wg.Add(height)
	for py := 0; py < height; py++ {
		go func(py int) {
			defer wg.Done()
			y := float64(py)/height*(yMax-yMin) + yMin
			for px := 0; px < width; px++ {
				x := float64(px)/width*(xMax-xMin) + xMin
				z := complex(x, y)
				img.Set(px, py, mandelbrot(z))
			}
		}(py)
	}
	wg.Wait()
	if print {
		fmt.Println("calc cost " + time.Now().Sub(start).String())
	}
	start = time.Now()

	buff := bytes.Buffer{}
	_ = png.Encode(&buff, img)
	if print {
		fmt.Println("code cost " + time.Now().Sub(start).String())
	}
	start = time.Now()

	hash := sha256.Sum256(buff.Bytes())
	if print {
		fmt.Println("hash cost " + time.Now().Sub(start).String())
		fmt.Println(hash)
	}
	return hash
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
