package main

import (
  "fmt"
  "time"
  "math/rand"
)

func main() {
  c := fanIn(
    boring("Joe"),
    boring("Ann"),
  )
  
  for i := 0; i < 5; i++ {
    msg1 := <-c
    fmt.Println(msg1.str)

    msg2 := <-c
    fmt.Println(msg2.str)

    // Tell Joe to run
    msg1.wait <- true

    // Tell Ann to run
    msg2.wait <- true
  }

  fmt.Println("You're both boring; I'm leaving.")
}

type Message struct {
  str string
  wait chan bool
}

func boring(msg string) <- chan Message {
  waitForIt := make(chan bool)
  c := make(chan Message)

  go func() {
    for i := 0; ;i++ {
      c <- Message{
        str: fmt.Sprintf("%s %d", msg, i),
        wait: waitForIt,
      }

      time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)

      // Block until main tells us to go again
      <- waitForIt
    }
  }()

  return c
}

func fanIn(input1, input2 <- chan Message) <- chan Message {
  c := make(chan Message)

  go func() {
    for {
      c <- <- input1
    }
  }()

  go func() {
    for {
      c <- <- input2
    }
  }()

  return c
}
