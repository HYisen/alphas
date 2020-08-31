package utility

import (
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func Fetch(method string, url string, jsonStringPayload io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, jsonStringPayload)
	if err != nil {
		return nil, err
	}
	if jsonStringPayload != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return http.DefaultClient.Do(req)
}

func GetHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer CloseAndLogError(resp.Body)

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func TraverseHTML(root *html.Node, interceptor func(*html.Node)) {
	candidates := []*html.Node{root}
	// Recursion would be a burden on performance.
	for len(candidates) != 0 {
		node := candidates[0]
		candidates = candidates[1:] // pop first
		if node.Type == html.ElementNode {
			interceptor(node)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			candidates = append(candidates, child) // push last
		}
	}
}

func SaveToFile(doc *html.Node, filepath string) error {
	err := os.MkdirAll(path.Dir(filepath), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer CloseAndLogError(file)

	err = html.Render(file, doc)
	if err != nil {
		return err
	}

	return nil
}

func ExtractFilename(url string) string {
	// remove protocol
	if strings.HasPrefix(url, "http") {
		// e.g., https://a.com -> a.com
		url = url[strings.Index(url, ":")+3:]
	}

	if index := strings.LastIndex(url, "/"); index != -1 {
		// e.g., aa/bb/cc -> cc
		url = url[index+1:]
	}

	return url
}

func ExtractNameAndExt(filename string) (name, ext string) {
	index := strings.Index(filename, ".")

	if index == -1 {
		name = filename
		ext = ""
		return
	}

	name = filename[:index]
	ext = filename[index+1:]
	return
}

func URLHasExt(url string) bool {
	_, ext := ExtractNameAndExt(ExtractFilename(url))
	return len(ext) != 0
}
