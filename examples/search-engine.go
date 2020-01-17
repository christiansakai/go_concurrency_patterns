package main

import (
  "fmt"
  "time"
  "math/rand"
)

func main() {
  rand.Seed(time.Now().UnixNano())

  start := time.Now()
  results := Google("golang")
  elapsed := time.Since(start)

  fmt.Println(results)
  fmt.Println(elapsed)
}

var (
  Web = fakeSearch("web")
  Image = fakeSearch("image")
  Video = fakeSearch("video")
)

type Search func(query string) Result

type Result string

func fakeSearch(kind string) Search {
  return func(query string) Result {
    time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

    return Result(fmt.Sprintf("%s result for %q\n", kind, query))
  }
}

func Google(query string) (results []Result) {
  c := make(chan Result)

  // Send query to separate backends concurrently
  go func() {
    // Backends have replicas so if one fails
    // there are still others
    c <- First(query, Web, Web)
  }()

  // Send query to separate backends concurrently
  go func() {
    // Backends have replicas so if one fails
    // there are still others
    c <- First(query, Image, Image)
  }()

  // Send query to separate backends concurrently
  go func() {
    // Backends have replicas so if one fails
    // there are still others
    c <- First(query, Video, Video)
  }()

  // Timeout for those backends
  timeout := time.After(80 * time.Millisecond)

  for i := 0; i < 3; i++ {
    select {
    case result := <- c:
      results = append(results, result)

    case <- timeout:
      fmt.Println("timed out")
      return
    }
  }

  return
}

// Function to run replicas of Search but
// only return the first one
func First(query string, replicas ...Search) Result {
  c := make(chan Result)

  searchReplica := func(i int) {
    c <- replicas[i](query)
  }

  for i := range replicas {
    go searchReplica(i)
  }

  return <- c
}
