package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// du

var verbose = flag.Bool("v", false, "show progress")
var threads = flag.Int("threads", 1, "concurrency")
var semaphore chan struct{}

func main() {
	flag.Parse()
	semaphore = make(chan struct{}, *threads)

	roots := flag.Args()
	if len(roots) == 0 {
		roots = append(roots, ".")
	}

	var wg sync.WaitGroup
	wg.Add(len(roots) + 1) // each du and total reporter

	totalFile := make(chan int, len(roots))
	totalByte := make(chan int, len(roots))
	// total reporter
	go func() {
		fileCnt := 0
		byteCnt := 0
		// Length shall be exactly match.
		for i := 0; i < 2*len(roots); i++ {
			select {
			case f := <-totalFile:
				fileCnt += f
			case b := <-totalByte:
				byteCnt += b
			}
		}
		printDiskUsage("total", fileCnt, byteCnt)
		wg.Done()
	}()

	// counters
	for _, root := range roots {
		root := root
		go func() {
			fileSizes := make(chan int)
			var tick <-chan time.Time
			if *verbose {
				tick = time.Tick(1 * time.Second)
			}

			go func() {
				walkDir(root, fileSizes)
				close(fileSizes)
			}()

			fileCnt := 0
			byteCnt := 0
			for {
				select {
				case <-tick:
					fmt.Println("tick")
					printDiskUsage(root, fileCnt, byteCnt)
				case size, ok := <-fileSizes:
					if !ok {
						printDiskUsage(root, fileCnt, byteCnt)
						totalFile <- fileCnt
						totalByte <- byteCnt
						wg.Done()
						return
					}
					fileCnt++
					byteCnt += size
				}
			}
		}()
	}

	wg.Wait()
}

func printDiskUsage(name string, fileCnt, byteCnt int) {
	fmt.Printf("[%s] %d files %d MiB\n", name, fileCnt, byteCnt>>20)
}

func walkDir(dir string, fileSizes chan<- int) {
	for _, entry := range access(dir) {
		if entry.IsDir() {
			childDir := filepath.Join(dir, entry.Name())
			walkDir(childDir, fileSizes)
		} else {
			fileSizes <- int(entry.Size())
		}
	}
}

func access(dir string) []os.FileInfo {
	semaphore <- struct{}{}
	defer func() { _ = <-semaphore }()

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "at [%s]:%v\n", dir, err)
		return []os.FileInfo{}
	}
	return entries
}
