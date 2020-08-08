package utility

import (
	"bufio"
	"log"
	"strings"
)

// ExtractTitleAndBody returns (title, body, isBodyExist)
func ExtractTitleAndBody(source string) (string, string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(source))
	var title, body string
	hasMeetFirstEmptyLine := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") {
			continue
		}

		if !hasMeetFirstEmptyLine {
			// title
			if len(line) == 0 {
				hasMeetFirstEmptyLine = true
			} else {
				title = concatLine(title, line)
			}
		} else {
			// body
			body = concatLine(body, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}
	if !hasMeetFirstEmptyLine {
		return title, body, false
	}
	return title, body, true
}

// concatLine return neoContent, which is lines' contacted result with \n.
func concatLine(oldContent string, neoLine string) string {
	if len(oldContent) != 0 {
		oldContent += "\n"
	}
	oldContent += neoLine
	return oldContent
}
