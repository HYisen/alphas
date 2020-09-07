// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 278.

package main

import (
	"fmt"
	"sync"
	"time"
)

//!+Func

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result   // the client wants a single result
	done     <-chan struct{} // close to cancel
}

type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string, done <-chan struct{}) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response, done}
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

//!+monitor

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	precursor := make(map[string]<-chan struct{})
	for req := range memo.requests {
		e := cache[req.key]
		var resend chan<- request = memo.requests
		var canResend chan struct{}
		var done <-chan struct{}

		if e == nil {
			fmt.Printf("calc done=%v\n", req.done)

			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key, req.done) // call f(key)

			resend = nil    // initiator shall not resend
			done = req.done // only initiator can be cancelled
			canResend = make(chan struct{})
			precursor[req.key] = canResend

			// clean up on cancel
			go func() {
				select {
				case <-e.ready: // succeed
					fmt.Println("already succeeded")
				case <-req.done: // cancel
					//e = nil // too young too simple, sometime naive. I thought it's reference.
					delete(cache, req.key)
					close(canResend)
				}
			}()
		}
		go e.deliver(req.response, done, resend, precursor[req.key], req)
	}
}

func (e *entry) call(f Func, key string, done <-chan struct{}) {
	fmt.Println("go " + key)
	job := make(chan result)
	go func() {
		value, err := f(key)
		job <- result{
			value: value,
			err:   err,
		}
	}()

	select {
	case r := <-job:
		// Evaluate the function.
		// Broadcast the ready condition.
		e.res.value, e.res.err = r.value, r.err
		close(e.ready)
	case <-done:
		fmt.Println("catch cancel signal")
		<-job // silently drop result
	}
}

func isClosed(channel <-chan struct{}) bool {
	select {
	case <-channel:
		return true
	default:
		return false
	}
}

func (e *entry) deliver(response chan<- result, done <-chan struct{},
	resendNullable chan<- request, canResend <-chan struct{}, req request) {
	fmt.Printf("delivering %v\n", resendNullable)
	for {
		select {
		case <-e.ready: // Wait for the ready condition.
			response <- e.res // Send the result to the client.
			return
		case <-done:
			fmt.Printf("0done %v\n", resendNullable)
			// resend if it is not the cancel origin
			if resendNullable != nil {
				_ = <-canResend
				resendNullable <- req
			} else {
				response <- result{
					value: "cancelled0",
					err:   nil,
				} // Send the result to the client.
			}
			fmt.Printf("1done %v\n", resendNullable)
			return
		case <-canResend:
			if isClosed(done) {
				response <- result{
					value: "cancelled1",
					err:   nil,
				}
				return
			}
			fmt.Println("resending")
			resendNullable <- req
			return
		}
	}

}

//!-monitor

var wg sync.WaitGroup

func main() {
	m := New(func(key string) (interface{}, error) {
		time.Sleep(2 * time.Second)
		return key, nil
	})

	wg.Add(2)

	done := make(chan struct{})

	go func() {
		defer wg.Done()
		fmt.Println(m.Get("one", nil))
	}()
	//time.Sleep(200 * time.Millisecond)
	go func() {
		defer wg.Done()
		fmt.Println(m.Get("one", done))
	}()
	time.Sleep(1 * time.Second)
	close(done)
	wg.Wait()
	//fmt.Println(m.Get("one", nil))
	//fmt.Println(m.Get("one", nil))
	//fmt.Println(m.Get("two", nil))
}
