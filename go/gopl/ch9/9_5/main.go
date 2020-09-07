package main

import (
	"fmt"
	"time"
)

func main() {
	ping := make(chan int)
	pong := make(chan int)
	go func() {
		for {
			num := <-ping
			num++
			pong <- num
		}
	}()
	go func() {
		cnt := 0
		for {
			num := <-pong
			if (num>>20)&1 == 1 {
				num = 0
				cnt++
				fmt.Printf("%dMi\n", cnt)
			}
			ping <- num
		}
	}()

	ping <- 0
	time.Sleep(10 * time.Second)
}
