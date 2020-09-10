package tar

import (
	"alphas/go/gopl/ch10/10_2/expressutil"
	"archive/tar"
	"fmt"
	"io"
)

func init() {
	fmt.Println("init tar")
	expressutil.RegisterFormat("tar", isTar, depress)
}

func isTar(data []byte) bool {
	return string(data[257:262]) == "ustar"
}

func depress(r io.Reader, w io.Writer) error {
	reader := tar.NewReader(r)
	for {
		head, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "header %v\n", head)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			return err
		}
	}
	return nil
}
