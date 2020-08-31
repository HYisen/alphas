package main

import (
	"alphas/go/gopl/utility"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var depth = flag.Int("depth", 2, "max jump since first page")
var first = flag.String("first", "http://www.fireemblem.net/fe/", "the home page from which crawlers start")
var sleep = flag.Int("sleep", 500, "the sleep time in millisecond of each crawler")
var count = flag.Int("count", 5, "the size of crawlers, a.k.a. concurrency")
var vault = flag.String("vault", "/tmp/web", "the save directory")

var website string

type job struct {
	url   string
	depth int
	doc   *html.Node
}

func addPotentialLostExt(url *string) {
	if !utility.URLHasExt(*url) {
		*url += ".html"
	}
}

func main() {
	flag.Parse()

	// assure output dir
	err := os.RemoveAll(*vault)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(*vault, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// handle website from *first
	website, err = extractWebsite(*first)
	if err != nil {
		log.Fatal(err)
	}

	// only download need concurrency
	results := search()

	for link, node := range results {
		traverseAndHandleHref(node, func(attribute *html.Attribute) {
			if _, ok := results[attribute.Val]; !ok {
				return
			}

			neo := genNeoHref(attribute.Val)
			if len(neo) == 0 {
				return
			}
			attribute.Val = neo
		})

		err := utility.SaveToFile(node, genFilepath(link))
		if err != nil {
			log.Print(err)
			continue
		}
	}
}

func genFilepath(url string) string {
	if len(url) == 0 {
		// slash
		return path.Join(*vault, "index.html")
	}
	addPotentialLostExt(&url)

	// Redundant, as I use absolute url as input. HasPrefix must be true.
	// http://web.com/a.html -> web.com/a.html
	if strings.HasPrefix(url, website) {
		url = url[len(website):]
	}

	return path.Join(append([]string{*vault}, strings.Split(url, "/")...)...)
}

// genNeoHref return empty string if no need to replace the old one.
func genNeoHref(old string) string {
	// self reference
	// tiring to optimize by extract similar check on old[:-1]
	if len(old) == 0 || old == *first || old+"/" == *first {
		return "./index.html"
	}

	// . or .., nice situation
	if old[0] == '.' {
		return ""
	}

	// protocol ones
	if strings.HasPrefix(old, "http") {
		// same site (untested)
		if strings.HasPrefix(old, website) {
			filepath := old[len(website):]
			addPotentialLostExt(&filepath)
			return "./" + filepath
		}

		return ""
	}

	// direct ones
	addPotentialLostExt(&old)
	return "./" + old
}

func extractWebsite(indexURL string) (string, error) {
	// ends with slash
	if indexURL[len(indexURL)-1] == '/' {
		return indexURL, nil
	}

	// ends with index (untested)
	if name, _ := utility.ExtractNameAndExt(utility.ExtractFilename(indexURL)); name == "index" {
		return indexURL[:strings.LastIndex(indexURL, "/")+1], nil
	}

	return "", fmt.Errorf("can not extract website from %s", indexURL)
}

func traverseAndHandleHref(root *html.Node, handler func(attribute *html.Attribute)) {
	utility.TraverseHTML(root, func(node *html.Node) {
		for i, attribute := range node.Attr {
			if attribute.Key == "href" {
				handler(&node.Attr[i])
				break
			}
		}
	})
}

func doSearchServant(task <-chan job, page chan<- job, candidate chan<- job) {
	for t := range task {
		// download
		doc, err := utility.GetHTML(t.url)
		if err != nil {
			log.Print(err)
			continue
		}

		// add potential next stage
		if t.depth < *depth {
			neoDepth := t.depth + 1
			base, err := url.Parse(t.url)
			if err != nil {
				log.Print(err)
				continue
			}
			traverseAndHandleHref(doc, func(attribute *html.Attribute) {
				href := attribute.Val
				if !strings.HasPrefix(href, "http") {
					u, err := url.Parse(href)
					if err != nil {
						log.Print(err)
						return
					}
					href = base.ResolveReference(u).String()
				}
				candidate <- job{
					url:   href,
					depth: neoDepth,
					doc:   nil,
				}
			})
		}

		// send download data
		t.doc = doc
		page <- t

		// rest and clean up
		time.Sleep(time.Duration((*sleep) * int(time.Millisecond)))
	}
}

func search() map[string]*html.Node {
	var wg sync.WaitGroup
	task := make(chan job, 64)
	candidate := make(chan job, 64) // match depth requirement, uncertain for others like prefix, may duplicate
	page := make(chan job, 64)

	// workers
	for i := 0; i < *count; i++ {
		go doSearchServant(task, page, candidate)
	}

	// manager
	go func() {
		seen := make(map[string]bool)

		for c := range candidate {
			if strings.HasPrefix(c.url, website) && !seen[c.url] {
				wg.Add(1)
				task <- c
			}
		}
	}()

	// overseer
	report := make(chan string)
	update := time.Now()
	// I can not use make(chan struct{}, 1) to simulate the lock, which would occur a dead lock.
	// Anyway, how's the importance to guard a pointer? Which is prone to be atomic.
	go func(report <-chan string, update *time.Time) {
		for r := range report {
			now := time.Now()
			fmt.Printf("%v done %s\n", now.Format("15:04:05"), r)

			*update = now
		}
	}(report, &update)
	go func() {
		for {
			time.Sleep(2 * time.Second)
			now := time.Now()
			if now.After(update.Add(5 * time.Second)) {
				fmt.Printf("%v Timeout, release the final wg.\n", now.Format("15:04:05"))
				wg.Done()
			}
		}
	}()

	// collector
	ret := make(map[string]*html.Node)
	go func(ret map[string]*html.Node, report chan<- string) {
		for p := range page {
			// wg also guard this check by its no minus assert
			if _, ok := ret[p.url]; ok {
				panic("duplicate download " + p.url)
			}

			ret[p.url[len(website):]] = p.doc
			report <- p.url
			wg.Done()
		}
	}(ret, report)

	// bootstrap
	candidate <- job{
		url:   *first,
		depth: 0,
		doc:   nil,
	}
	wg.Add(1) // extra one that shall be released by overseer after timeout
	wg.Wait()

	return ret
}
