package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"log"
	"os"
)

var SHALength = flag.Int("SHALength", 256, "SHA length, candidates are 256/384/512")

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()

	var str string

	switch *SHALength {
	case 256:
		str = fmt.Sprintf("%x", sha256.Sum256([]byte(text)))
	case 512:
		str = fmt.Sprintf("%x", sha512.Sum512([]byte(text)))
	case 384:
		str = fmt.Sprintf("%x", sha512.Sum384([]byte(text)))
	default:
		log.Fatalf("bad SHA length %d", *SHALength)
	}

	fmt.Println(str)
}
