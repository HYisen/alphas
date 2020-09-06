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

	timeout = flag.Int("timeout", 60, "close client after silence of such long seconds")
	bufSize = flag.Int("bufSize", 2, "max most ancient messages in buffer if client does not accept")
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

func contain(key string, list []string) bool {
	for _, one := range list {
		if one == key {
			return true
		}
	}
	return false
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	requireUsers <- struct{}{}
	users := <-currentUsers
	ch <- fmt.Sprintf("current users: %v", users)

	ch <- "Your name?"
	input := bufio.NewScanner(conn)
	input.Scan()
	name := input.Text()
	for contain(name, users) {
		ch <- "409 Conflict, your name?"
		input.Scan()
		name = input.Text()
	}
	fmt.Printf("register %s as %s\n", conn.RemoteAddr().String(), name)
	who := name
	ch <- "So you are " + who

	messages <- who + " has arrived"

	cli := client{
		name: who,
		ch:   ch,
	}
	entering <- cli

	respawn := make(chan struct{})
	go utility.SummonDestroyer(conn, time.Duration((*timeout)*int(time.Second)), func() {
		fmt.Println(who + " is kicking off because of silence")
		ch <- fmt.Sprintf("kicking you because of keeping silence for %d sec", *timeout)
	}, respawn)

	for input.Scan() {
		messages <- who + ": " + input.Text()
		respawn <- struct{}{}
	}

	leaving <- cli
	messages <- who + " has left"
	utility.CloseAndLogError(conn)
}

func clientWriter(conn net.Conn, ch <-chan string) {
	buff := make(chan string, *bufSize)

	go func() {
		for msg := range buff {
			_, _ = fmt.Fprintln(conn, msg)
		}
	}()

	for msg := range ch {
		if len(buff) == cap(buff) {
			fmt.Printf("drop msg \"%s\" to %s\n", msg, conn.RemoteAddr().String())
			continue
		}
		buff <- msg
	}
}
