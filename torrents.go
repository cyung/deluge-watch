package main

import (
  "time"
  "net/http"
  // "encoding/json"
  "os"
  // "fmt"
  "math/rand"
  "io"
)

func GetTorrents() {
  url := GetBaseUrl() + "/torrents"

  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    panic(err)
  }
  req.Header.Add("Authorization", GetKey())

  for {
    res, err := client.Do(req)

    // continue if no zip file received
    if err != nil {
      time.Sleep(5 * time.Second)
      continue
    }

    if res.StatusCode != 200 {
      time.Sleep(5 * time.Second)
      continue
    }

    // save to zip file
    zip_filename := RandomFilename() + ".zip"
    file, err := os.Create("./tmp/" + zip_filename)
    if err != nil {
      panic(err)
    }

    _, err = io.Copy(file, res.Body)
    if err != nil {
      panic(err)
    }

    res.Body.Close()
    file.Close()

    time.Sleep(5 * time.Second)
  }
}

func RandomFilename() string {
  const CHAR_LENGTH = 10
  const chars = "abcdefghijklmnopqrstuvwxyz" +
                "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
                "0123456789"

  rand.Seed(time.Now().UTC().UnixNano())
  result := make([]byte, CHAR_LENGTH)

  for i := 0; i<CHAR_LENGTH; i++ {
    result[i] = chars[rand.Intn(len(chars))]
  }

  return string(result)
}

// func deleteTorrents(magnets *map[string]bool) {
//   for magnet, _ := range *magnets {
//     go sendDelete(magnet)
//   }
// }

// func sendDelete(magnet string) {
//   url := base_url + "/magnets?magnet=" + magnet

//   req, err := http.NewRequest("DELETE", url, nil)
//   if err != nil {
//     panic(err)
//   }

//   client := &http.Client{}

//   _, err = client.Do(req)
//   if err != nil {
//     panic(err)
//   }
// }