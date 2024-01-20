// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The context package example.
//
// The selected Go code is setting up a context example with goroutines and
// timeouts.
//
// It starts by importing the necessary packages including the context package.
//
// The main() function is the entry point. It first configures the logger to
// include microseconds in the log output.
//
// It then defines constants for the wait time before canceling the context and
// the number of goroutines to create.
//
// The code creates a parent context using context.WithCancel() which can be
// canceled by calling the returned cancel function.
//
// It initializes two WaitGroup variables to coordinate the goroutines. wg will
// wait for all goroutines to finish. wgStarted will wait for all goroutines to
// start before beginning the work.
//
// It starts a loop to launch the goroutines. Each goroutine gets a context
// with a randomized timeout. The goroutines call printHello() passing the
// context, waitgroup, and goroutine number.
//
// The printHello() function does the main work. It prints starting and done
// messages, decrements the waitgroup, and runs a select loop to wait on the
// context Done() channel or a 1 second timer. When Done() fires, it prints a
// message about the context being canceled.
//
// After launching the goroutines, a timer is created to cancel the parent
// context after the specified wait time.
//
// Finally, the code waits for the goroutine waitgroup to finish and stops the
// timer.
//
// So in summary, it launches goroutines with individual timeouts, cancels the
// parent context after a wait time, and orchestrates the goroutine lifetimes
// using waitgroups. The printHello() function shows how a goroutine can monitor
// a context's Done() channel to know when to exit.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {

	// Set log with microseconds
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Number of seconds befor cancel context
	const wait = 5 * time.Second

	// Number of goroutines
	const goroutines = 3

	// Parent context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to wait for all goroutines finished
	wg := &sync.WaitGroup{}

	// Create a WaitGroup to wait for all goroutines started to start goroutines
	// processing in one time after all goroutines started
	wgStarted := &sync.WaitGroup{}
	wgStarted.Add(goroutines)
	fmt.Println("Start goroutines:")

	// Call a function as a goroutine
	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		// Get random number of seconds from 1 to 12 for timeout
		r := time.Duration(1+rand.Intn(11)) * time.Second

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(ctx, r)
		fmt.Printf("Wait timeout %d %v, %v\n", i, r, ctx)

		go func(i int) { wgStarted.Wait(); printHello(ctx, wg, i); cancel() }(i)
		wgStarted.Done()
	}

	// Create timer which will wait for wait time and than cancel context
	t := time.AfterFunc(wait, cancel)

	// Wait for all goroutines
	wg.Wait()

	// Stop the timer
	t.Stop()

	log.Println("Hello from main, all done")
}

func printHello(ctx context.Context, wg *sync.WaitGroup, i int) {

	fmt.Println("Start", i)
	defer fmt.Println("Done", i)

	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// Calculate delay between deadline and current time
			var delay = func() (delay string) {
				deadline, _ := ctx.Deadline()
				if d := time.Since(deadline); d > 0 {
					delay = fmt.Sprintf(", deadline delay: %v", d)
				}
				return
			}

			// Print message
			log.Printf("Hello %d receive context done, %v%s\n",
				i, ctx.Err(), delay(),
			)

			return

		case <-time.After(1 * time.Second):
			fmt.Println("Hello from printHello", i)
		}
	}
}
