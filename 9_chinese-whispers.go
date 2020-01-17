package main

import "fmt"

// Daisy chain long communications
// From tail (right) to head (left)
func main() {
	const n = 10000

	leftmost := make(chan int)
	right := leftmost
	left := leftmost

	for i := 0; i < n; i++ {
		right = make(chan int)

		go f(left, right)

		left = right
	}

	go func(c chan int) {
		c <- 1
	}(right)

	fmt.Println(<-leftmost)
}

func f(left, right chan int) {
	left <- 1 + <-right
}
