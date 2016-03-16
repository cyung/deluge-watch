package main

import (
  "time"
  "net/http"
  "encoding/json"
  "os"
)

const BASE_URL string = "http://localhost:3000"

func GetMagnets() {
  url := BASE_URL + "/magnets"

  for {
    res, err := http.Get(url)

    if err != nil {
      panic(err)
    }

    var magnets map[string]bool
    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&magnets)
    if err != nil {
      panic(err)
    }

    if len(magnets) > 0 {
      createMagnets(&magnets)
      deleteMagnets(&magnets)
    }

    time.Sleep(5 * time.Second)
  }
}

func createMagnets(magnets *map[string]bool) {
  file, err := os.Create("magnets.magnet")
  if err != nil {
    panic(err)
  }


  for magnet, _ := range *magnets {
    _, err := file.WriteString(magnet + "\n")
    if err != nil {
      panic(err)
    }
  }

  err = file.Close()
  if err != nil {
    panic(err)
  }
}

func deleteMagnets(magnets *map[string]bool) {
  for magnet, _ := range *magnets {
    go sendDelete(magnet)
  }
}

func sendDelete(magnet string) {
  url := BASE_URL + "/magnets?magnet=" + magnet

  req, err := http.NewRequest("DELETE", url, nil)
  if err != nil {
    panic(err)
  }

  client := &http.Client{}

  _, err = client.Do(req)
  if err != nil {
    panic(err)
  }
}