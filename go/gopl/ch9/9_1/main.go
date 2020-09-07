// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 261.
//!+

package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go Deposit(1000)
	go Deposit(1000)
	go Deposit(1000)
	go Deposit(1000)
	fmt.Println(Balance())
	go func() {
		fmt.Println(Withdraw(2000))
		wg.Done()
	}()
	fmt.Println(Balance())
	go Withdraw(2000)
	go Deposit(1000)
	fmt.Println(Balance())
	go func() {
		fmt.Println(Withdraw(8000))
		wg.Done()
	}()
	fmt.Println(Balance())
	wg.Wait()
	fmt.Println(Balance())
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdrawAttempt = make(chan int)
var withdrawResult = make(chan bool)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	withdrawAttempt <- amount
	return <-withdrawResult
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case balances <- balance:
		case amount := <-withdrawAttempt:
			if amount > balance {
				withdrawResult <- false
				continue
			}
			balance -= amount
			withdrawResult <- true
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
