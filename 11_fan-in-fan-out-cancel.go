package main

import (
	"fmt"
	"sync"
)

func main() {
	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline  exists, as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	defer close(done)

	in := gen(done, 2, 3)

	// Distribute the sq work across two goroutines that both read from in
	c1 := sq(done, in)
	c2 := sq(done, in)

	for n := range merge(done, c1, c2) {
		fmt.Println(n)
	}
}

func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			select {
			case out <- n:

			case <-done:
				return
			}
		}
	}()

	return out
}

func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			select {
			case out <- n * n:

			case <-done:
				return
			}
		}
	}()

	return out
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.
	// Output copies values from c to out until c is closed, then calls wg.Done
	output := func(c <-chan int) {
		defer wg.Done()

		for n := range c {
			select {
			case out <- n:

			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call
	go func() {
		defer close(out)

		wg.Wait()
	}()

	return out
}
