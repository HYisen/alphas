package main

import (
	"alphas/go/gopl/utility"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// mirroredQuery

func main() {
	flag.Parse()
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(flag.Args()))
	for _, url := range flag.Args() {
		go func(url string) {
			defer func() { wg.Done() }()

			resp, err := get(url, done)
			if err != nil {
				fmt.Printf("failed %s : %v\n", url, err)
				return
			}
			fmt.Printf("succeed %s status %v\n", url, resp.Status)

			var sb strings.Builder
			_, err = io.Copy(&sb, resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(sb.String())

			utility.CloseAndLogError(resp.Body)
			close(done)
		}(url)
	}
	wg.Wait()
}

// do not forget to close the response returned
func get(url string, done <-chan struct{}) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// inject cancel stub
	cx, cancel := context.WithCancel(context.Background())
	request = request.WithContext(cx)
	go func() {
		select {
		case <-cx.Done(): // already done
		case <-done:
			cancel()
			fmt.Println("canceled " + url)
		}
	}()

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
