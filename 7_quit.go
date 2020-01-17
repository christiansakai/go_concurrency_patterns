package main

import (
  "fmt"
  "time"
  "math/rand"
)

func main() {
  quit := make(chan bool)
  c := boring("Joe", quit)

  for i := rand.Intn(10); i >= 0; i-- {
    fmt.Println(<- c)
  }

  quit <- true
}

func boring(msg string, quit <- chan bool) <- chan string {
  c := make(chan string)

  go func() {
    for i := 0; ;i++ {
      select {
        case c <- fmt.Sprintf("%s %d", msg, i):
          time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

        case <- quit:
          return
      }
    }
  }()

  return c
}