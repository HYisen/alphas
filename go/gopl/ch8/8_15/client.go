package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	var requestPrint = make(chan struct{})
	var respondPrint = make(chan struct{})
	var triggerPrint = make(chan bool)

	// audience
	go func() {
		canPrint := false
		cnt := 0
		for {
			select {
			case neoState := <-triggerPrint:
				canPrint = neoState
				if canPrint {
					for ; cnt > 0; cnt-- {
						respondPrint <- struct{}{}
					}
				}
			case <-requestPrint:
				if canPrint {
					respondPrint <- struct{}{}
				} else {
					cnt++
				}
			}
		}
	}()

	go func() {
		iStream := bufio.NewScanner(conn)
		for iStream.Scan() {
			requestPrint <- struct{}{}
			<-respondPrint
			fmt.Println(iStream.Text())
		}
	}()

	oStream := bufio.NewScanner(os.Stdin)
	triggerPrint <- true
	for oStream.Scan() {
		line := oStream.Text()
		switch line {
		case ":printOn":
			triggerPrint <- true
		case ":printOff":
			triggerPrint <- false
		default:
			_, _ = conn.Write([]byte(line + "\n"))
		}
	}
}
