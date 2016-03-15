package main

import (
  "time"
  "net/http"
  "encoding/json"
  "os"
)

type Torrent struct {
  Magnet string `json:"magnet"`
}

const base_url string = "http://localhost:3000"

func main() {
  url := base_url + "/torrents"

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
      createFile(&magnets)
      deleteFromServer(&magnets)
    }

    time.Sleep(5 * time.Second)
  }
}

func createFile(magnets *map[string]bool) {
  file, err := os.Create("torrents.magnet")
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

func deleteFromServer(magnets *map[string]bool) {
  for magnet, _ := range *magnets {
    go sendDelete(magnet)
  }
}

func sendDelete(magnet string) {
  url := base_url + "/torrents?magnet=" + magnet

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