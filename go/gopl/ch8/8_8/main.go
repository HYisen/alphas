package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

func echo(c net.Conn, shout string, delay time.Duration) {
	_, _ = fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	_, _ = fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	_, _ = fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func HandleConn(c net.Conn) {
	reset := make(chan struct{})

	go func(reset chan<- struct{}) {
		input := bufio.NewScanner(c)
		for input.Scan() {
			reset <- struct{}{}
			go echo(c, input.Text(), 1*time.Second)
		}
	}(reset)

	delay := 10 * time.Second
	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			fmt.Println("close after silence of " + delay.String())
			log.Print("CloseWrite " + c.RemoteAddr().String())
			err := c.(*net.TCPConn).CloseWrite()
			if err != nil {
				log.Print(err)
			}
			return
		case <-reset:
			timer.Reset(delay)
		}
	}
}
