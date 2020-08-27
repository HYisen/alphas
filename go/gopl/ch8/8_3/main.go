package main

import (
	"alphas/go/gopl/utility"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()
	utility.MustCopy(conn, os.Stdin)
	err = conn.(*net.TCPConn).CloseWrite()
	if err != nil {
		log.Println(err)
	}
	utility.CloseAndLogError(conn)
	<-done

}
