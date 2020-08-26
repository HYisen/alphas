package main

import (
	"alphas/go/gopl/utility"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	for _, param := range os.Args[1:] {
		split := strings.Split(param, "=")
		name := split[0]
		addr := split[1]
		go run(name, addr)
	}
	go screenNextLine(1 * time.Second)
	time.Sleep(1 * time.Hour)
}

func screenNextLine(sleepTime time.Duration) {
	for {
		time.Sleep(sleepTime)
		fmt.Println()
	}
}

func run(name, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer utility.CloseAndLogError(conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		fmt.Printf("%-8s:%-12s\t", name, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Printf("err at %s=%s : %v\n", name, addr, err)
	}
}
