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

	"github.com/pbnjay/memory"
)

// Number of goroutines
const goroutines = 5

func main() {

	// Set log with microseconds
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Start time
	var start = time.Now()

	// Number of seconds befor cancel context
	const wait = 5 * time.Second

	// Parent context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a WaitGroup to wait for all goroutines started to start goroutines
	// processing in one time after all goroutines started.
	wgStarted := &sync.WaitGroup{}
	wgStarted.Add(goroutines)

	// Can't start goroutine flag when free memory less than 10 MB
	first := true

	// Started goroutines counter
	var started int

	// WaitGroup to wait for all goroutines finished
	wg := &sync.WaitGroup{}

	// Start goroutines loop
	fmt.Printf("Try to start %d goroutines:\n", goroutines)
	for i := 0; i < goroutines; i++ {

		// If free memory less than 10 MB than skip all other goroutines creation
		if !first || memory.FreeMemory() < 10*1024*1024 { // 10 MB
			if first {
				fmt.Printf("Free memory %d, %d goroutines started. "+
					"Can't start any more.\n", memory.FreeMemory(), i)
				first = false
			}
			wgStarted.Done()
			continue
		}

		// Get random number of seconds from 1 to 12 for timeout
		r := time.Duration(1+rand.Intn(11)) * time.Second

		// Create a context with timeout
		ctx, cancel := context.WithTimeout(ctx, r)
		if goroutines <= 10 {
			fmt.Printf("Wait timeout %d %v, %v\n", i, r, ctx)
		}

		// Start goroutine task
		wg.Add(1)
		go func(i int) {
			started++          // Number of started goroutines
			wgStarted.Done()   // Task is started
			wgStarted.Wait()   // Wait when all goroutines started
			printHello(ctx, i) // Execute task body
			wg.Done()          // This task is done
			cancel()           // Cancel context to avoid memory leaks
		}(i)

	}

	// Get all task started and show time
	wgStarted.Wait()
	fmt.Printf("All %d goroutines started in %v\n", started, time.Since(start))
	fmt.Printf("Total system memory: %d\n", memory.TotalMemory())
	fmt.Printf("Free system memory: %d\n", memory.FreeMemory())

	// Create timer which will wait for wait time and than cancel context
	t := time.AfterFunc(wait, cancel)

	// Wait for all goroutines
	wg.Wait()

	// Stop the timer
	t.Stop()

	log.Println("Hello from main, all done", time.Since(start))
}

// printHello is task body, it prints a message and waits ctx.Done() or 
// for a second.
func printHello(ctx context.Context, i int) {

	if goroutines <= 10 {
		fmt.Println("Start", i)
		defer fmt.Println("Done", i)
	}

	for {
		select {
		case <-ctx.Done():

			// Calculate delay between deadline and current time
			if goroutines <= 10 {
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
			}

			return

		case <-time.After(1 * time.Second):
			if goroutines <= 10 {
				fmt.Println("Hello from printHello", i)
			}
		}
	}
}
