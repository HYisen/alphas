package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/l", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	params, err := singular(r.URL.Query())
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		fmt.Println(w, err)
	}

	fmt.Println(params)
	if err := lissajous(w, params); err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, err)
		fmt.Println(w, err)
	}
}

func singular(values url.Values) (map[string]string, error) {
	ret := make(map[string]string)
	for k, v := range values {
		if len(v) != 1 {
			return ret, fmt.Errorf("bad query %s->%v, not injection", k, v)
		}
		ret[k] = v[0]
	}
	return ret, nil
}

func lissajous(out io.Writer, params map[string]string) error {
	var (
		cycles = 5
		res    = 0.001
		size   = 100
		frames = 64
		delay  = 8
		step   = 16
	)

	// With the help of generic type and inject function, I can generalize following similar procedures, do I?
	var err error = nil
	if elem, ok := params["cycles"]; ok {
		cycles, err = strconv.Atoi(elem)
		if err != nil {
			return err
		}
	}
	if elem, ok := params["res"]; ok {
		res, err = strconv.ParseFloat(elem, 32)
		if err != nil {
			return err
		}
	}
	if elem, ok := params["size"]; ok {
		size, err = strconv.Atoi(elem)
		if err != nil {
			return err
		}
	}
	if elem, ok := params["frames"]; ok {
		frames, err = strconv.Atoi(elem)
		if err != nil {
			return err
		}
	}
	if elem, ok := params["delay"]; ok {
		delay, err = strconv.Atoi(elem)
		if err != nil {
			return err
		}
	}
	if elem, ok := params["step"]; ok {
		step, err = strconv.Atoi(elem)
		if err != nil {
			return err
		}
	}

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

		for t := 0.0; t < float64(cycles*2)*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			//fmt.Fprintf(os.Stderr, "%d %d \n", size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5))
			img.SetColorIndex(
				size+int(x*float64(size)+0.5),
				size+int(y*float64(size)+0.5),
				uint8(int((x*float64(step)+float64(step))/2)))
		}
		//fmt.Fprintf(os.Stderr, "phase=%v\n", phase)
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	// lzw: input byte too large for the litWidth
	// ignore it is okay
	//noinspection GoUnhandledErrorResult
	return gif.EncodeAll(out, &anim)
	//return nil
}
