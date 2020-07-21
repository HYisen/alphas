package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		prefix := "http://"
		if !strings.HasPrefix(url, prefix) {
			url = prefix + url
			fmt.Println("add prefix " + prefix)
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get "+url+" : %v\n", err)
			continue
		}

		fmt.Printf("%s return %d\n", url, resp.StatusCode)
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to copy body of "+url+" : %v\n", err)
			continue
		}
		resp.Body.Close()
	}
}
