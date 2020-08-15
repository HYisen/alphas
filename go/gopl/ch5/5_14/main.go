// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 139.

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"gopl.io/ch5/links"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string, pred func(string) bool) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				if pred(item) {
					if err := download(item); err != nil {
						log.Fatal(err)
					}
					worklist = append(worklist, f(item)...)
				}
			}
		}
	}
}

//!-breadthFirst

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	param := []string{"https://www.baidu.com/"}
	pred := genSameDomainPredicate(param[0])
	breadthFirst(crawl, param, pred)
}

//!-main

func genSameDomainPredicate(url string) func(orig string) bool {
	keys := extractCenter(url)
	return func(orig string) bool {
		array := extractCenter(orig)
		if array == nil || len(array) < 2 {
			return false
		}
		if !reflect.DeepEqual(keys[len(keys)-2:], array[len(array)-2:]) {
			return false
		}
		return true
	}
}

func extractCenter(url string) []string {
	for _, limb := range strings.Split(url, "/") {
		if strings.Contains(limb, ".") {
			return strings.Split(limb, ".")
		}
	}
	return nil
}

func download(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	file, err := os.Create(strings.Replace(url, "/", "$", -1))
	if err != nil {
		_ = resp.Body.Close()
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		_ = file.Close()
		return err
	}
	_ = resp.Body.Close()
	_ = file.Close()
	return nil
}
