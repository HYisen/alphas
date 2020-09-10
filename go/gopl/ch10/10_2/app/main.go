package main

import (
	"alphas/go/gopl/ch10/10_2/expressutil"
	_ "alphas/go/gopl/ch10/10_2/expressutil/tar"
	_ "alphas/go/gopl/ch10/10_2/expressutil/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var input = flag.String("input", "", "path of input file, empty means read from std::cin")

func getData() []byte {
	var data []byte

	if len(*input) == 0 {
		file, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalln(err)
		}
		data = file
	} else {
		file, err := ioutil.ReadFile(*input)
		if err != nil {
			log.Fatalln(err)
		}
		data = file
	}

	return data
}

func main() {
	flag.Parse()
	data := getData()

	fmt.Println(len(data))
	str, err := expressutil.Depress(data)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(str)
}
