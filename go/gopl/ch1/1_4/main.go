package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	data := make(map[string][]string)
	filenames := os.Args[1:]

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open %s : %v\n", filename, err)
			return
		}

		extract(file, data)
		file.Close()
	}

	for line, origins := range data {
		if len(origins) > 1 {
			fmt.Printf("%s occures in %v\n", line, origins)
		}
	}
}

func extract(f *os.File, data map[string][]string) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		// I do not like the functional style of const data.
		data[input.Text()] = append(data[input.Text()], f.Name())
	}
}
