package main

import (
  "fmt"
)

func main() {
  fmt.Println("Starting watch server")
  // go GetMagnets()
  go GetTorrents()

  for { }
}