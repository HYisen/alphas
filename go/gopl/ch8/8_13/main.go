package main

import (
	"alphas/go/gopl/utility"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

type client struct {
	name string
	ch   chan<- string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)

	requireUsers = make(chan struct{})
	currentUsers = make(chan []string)

	timeout = flag.Int("timeout", 5, "close client after silence of such long seconds")
)

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func extractUsers(orig map[client]bool) []string {
	var ret []string
	for cli, isOnline := range orig {
		if isOnline {
			ret = append(ret, cli.name)
		}
	}
	return ret
}

func broadcaster() {
	clients := make(map[client]bool) // whether online
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}
		case cli := <-entering:
			clients[cli] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		case <-requireUsers:
			currentUsers <- extractUsers(clients)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who

	requireUsers <- struct{}{}
	users := <-currentUsers

	ch <- fmt.Sprintf("current users: %v", users)
	messages <- who + " has arrived"

	cli := client{
		name: who,
		ch:   ch,
	}
	entering <- cli

	respawn := make(chan struct{})
	go utility.SummonDestroyer(conn, time.Duration((*timeout)*int(time.Second)), func() {
		fmt.Println(who + " is kicking off because of silence")
	}, respawn)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
		respawn <- struct{}{}
	}

	leaving <- cli
	messages <- who + " has left"
	utility.CloseAndLogError(conn)
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, _ = fmt.Fprintln(conn, msg)
	}
}
