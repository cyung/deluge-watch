package main

import (
  "io/ioutil"
  "log"
  "encoding/json"
)

type Configuration struct {
  ChrisKey string `json:"CHRIS_KEY"`
}

var key string
const base_url string = "http://localhost:3000"

func init() {
  file, err := ioutil.ReadFile("./config.json")
  if err != nil {
    log.Fatal(err)
  }

  var config Configuration
  err = json.Unmarshal(file, &config)
  if err != nil {
    log.Fatal(err)
  }

  key = config.ChrisKey
}

func GetKey() string {
  return key
}

func GetBaseUrl() string {
  return base_url
}