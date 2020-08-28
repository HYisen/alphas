package main

import (
	"alphas/go/gopl/utility"
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// reverb

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleConn(conn)
	}
}

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	_, _ = fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	_, _ = fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	_, _ = fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func HandleConn(c net.Conn) {
	defer utility.CloseAndLogError(c)
	input := bufio.NewScanner(c)
	var wg sync.WaitGroup
	for input.Scan() {
		wg.Add(1)
		go echo(c, input.Text(), 1*time.Second, &wg)
	}
	wg.Wait()
	log.Print("CloseWrite " + c.RemoteAddr().String())
	err := c.(*net.TCPConn).CloseWrite()
	if err != nil {
		log.Print(err)
	}
}
