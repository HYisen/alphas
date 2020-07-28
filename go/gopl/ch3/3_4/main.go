package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

var (
	width, height = 800, 480
	cells         = 100
	xyRange       = 4.0
	xyScale       = float64(width) / 2 / xyRange
	zScale        = float64(height) * 0.4
	angle         = math.Pi / 6
	color         = "byDepth"
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

type fType func(float64, float64) float64

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	fName := query["f"][0]

	fmt.Println(fName)

	var f fType

	switch fName {
	case "eggbox":
		f = fEggBox
	case "moguls":
		f = fMoguls
	case "saddle":
		f = fSaddle
	default:
		f = fDefault
	}

	h := query["RLH"]
	if h != nil {
		num, err := strconv.ParseFloat(h[0], 64)
		if err != nil {
			_, _ = fmt.Fprintf(w, "can not parse RLH %v", h)
			w.WriteHeader(400)
			return
		}
		f = addReferenceLine(f, num)
	}

	height, _ = strconv.Atoi(query["height"][0])
	width, _ = strconv.Atoi(query["width"][0])
	color = query["color"][0]

	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = fmt.Fprint(w, fmt.Sprint(gen(f)))
}

func gen(f fType) string {
	var ret string
	ret += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, az := corner(i+1, j, f)
			bx, by, bz := corner(i, j, f)
			cx, cy, cz := corner(i, j+1, f)
			dx, dy, dz := corner(i+1, j+1, f)

			if math.IsNaN(ax) || math.IsNaN(ay) ||
				math.IsNaN(bx) || math.IsNaN(by) ||
				math.IsNaN(cx) || math.IsNaN(cy) ||
				math.IsNaN(dx) || math.IsNaN(dy) {
				fmt.Printf("skip %d %d\n", i, j)
				continue
			}
			colorParam := color
			if colorParam == "byDepth" {
				colorParam = calcColor(az, bz, cz, dz)
			}
			ret += fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s' stroke='black' />\n",
				ax, ay, bx, by, cx, cy, dx, dy, colorParam)
		}
	}
	ret += fmt.Sprint("</svg>")
	return ret
}

func calcColor(az float64, bz float64, cz float64, dz float64) string {
	z := (az + bz + cz + dz) / 4.0
	if z > 1.0 {
		z = 1.0
	} else if z < -1.0 {
		z = -1.0
	}
	color := genColor(z * 5)
	return color
}

func corner(i, j int, f fType) (float64, float64, float64) {
	x := xyRange * (float64(i)/float64(cells) - 0.5)
	y := xyRange * (float64(j)/float64(cells) - 0.5)
	z := f(x, y)
	sx := float64(width/2) + (x-y)*cos30*xyScale
	sy := float64(height/2) + (x+y)*sin30*xyScale - z*zScale
	return sx, sy, z
}

func fDefault(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func fEggBox(x, y float64) float64 {
	//fmt.Printf("%f %f\n", x, y)

	mul := 2.5
	xCenter := math.Floor(math.Abs(x * mul))
	yCenter := math.Floor(math.Abs(y * mul))
	xDiff := math.Abs(x*mul) - xCenter - 0.5
	yDiff := math.Abs(y*mul) - yCenter - 0.5
	//fmt.Printf("%f %f\n", xDiff, yDiff)

	return (xDiff*xDiff + yDiff*yDiff - 0.25) / 4.0
}

func fMoguls(x, y float64) float64 {
	if int(math.Floor(math.Abs(x*10)))%4 == 0 && int(math.Floor(math.Abs(y*10)))%4 == 0 {
		return 0.1
	}
	return -0.1
}

func fSaddle(x, y float64) float64 {
	return (x*x - y*y) / 10
}

func addReferenceLine(orig fType, height float64) fType {
	return func(x float64, y float64) float64 {
		if x == 0 || y == 0 {
			return height
		}
		return orig(x, y)
	}
}

func genColor(z float64) string {
	d := uint16((z + 1.0) * 128.0)
	if d == 256 {
		d--
	}
	// d between [0,255]
	return fmt.Sprintf("#%02x00%02x", d, 255-d)
}
