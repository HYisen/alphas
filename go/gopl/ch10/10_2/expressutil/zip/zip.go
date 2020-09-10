package zip

import (
	"alphas/go/gopl/ch10/10_2/expressutil"
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

func init() {
	fmt.Println("init zip")
	expressutil.RegisterFormat("zip", isZip, depress)
}

func isZip(data []byte) bool {
	return binary.LittleEndian.Uint32(data[0:4]) == 0x04034b50
}

func depress(r io.Reader, w io.Writer) error {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	reader, err := zip.NewReader(bytes.NewReader(all), int64(len(all)))
	if err != nil {
		return err
	}
	for _, file := range reader.File {
		_, err = fmt.Fprintf(w, "file %v\n", file.Name)
		if err != nil {
			return err
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(w, open)
		if err != nil {
			return err
		}
	}
	return nil
}
