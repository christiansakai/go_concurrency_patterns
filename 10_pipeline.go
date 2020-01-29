package main

import "fmt"

func main() {
	for n := range sq(gen(2, 3)) {
		fmt.Println(n)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for n := range in {
			out <- n * n
		}
	}()

	return out
}
