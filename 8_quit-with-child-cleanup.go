package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	quit := make(chan string)
	c := boring("Joe", quit)

	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}

	quit <- "Bye!"

	fmt.Printf("Joe says: %q\n", <-quit)
}

func boring(msg string, quit chan string) <-chan string {
	c := make(chan string)

	go func() {
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

			case <-quit:
				cleanup()

				quit <- "See you!"

				return
			}
		}
	}()

	return c
}

func cleanup() {}
