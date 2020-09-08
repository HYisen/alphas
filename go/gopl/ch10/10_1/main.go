package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"

	"image/gif"
	"image/jpeg"
	"image/png"
)

var out = flag.String("out", "png", "output type, jpeg/png/gif")

func main() {
	flag.Parse()
	if err := convert(os.Stdin, os.Stdout); err != nil {
		log.Fatalln("convert: ", err)
	}
}

func convert(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintln(os.Stderr, "Input format = ", kind)
	return GetEncodeFunc()(out, img)
}

func GetEncodeFunc() func(writer io.Writer, image image.Image) error {
	switch *out {
	case "jpeg":
		return func(writer io.Writer, image image.Image) error {
			return jpeg.Encode(writer, image, &jpeg.Options{Quality: 95})
		}
	case "png":
		return png.Encode
	case "gif":
		return func(writer io.Writer, image image.Image) error {
			return gif.Encode(writer, image, &gif.Options{NumColors: 256})
		}
	default:
		return func(_ io.Writer, _ image.Image) error {
			return fmt.Errorf("unsupported output type %s", *out)
		}
	}
}
