package main

import (
  "time"
  "net/http"
  "os"
  "fmt"
  "math/rand"
  "io"
  "io/ioutil"
  "archive/zip"
  "errors"
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
    time.Sleep(5 * time.Second)

    zip_filename, err := saveZipfile(client, req)
    if err != nil {
      continue
    }

    err = unzip(zip_filename)
    if err != nil {
      continue
    }

    filenames, err := ackTorrents()
    if err != nil {
      continue
    }

    err = moveToWatchFolder(filenames)
    if err != nil {
      continue
    }
  }
}

func saveZipfile(client *http.Client, req *http.Request) (string, error) {
  res, err := client.Do(req)
  if err != nil {
    return "", err
  }

  if res.StatusCode != 200 {
    fmt.Println("No torrents on server")
    return "", errors.New("No torrents on server")
  }

  // save zip file locally
  zip_filename := "./tmp/" + RandomFilename() + ".zip"
  file, err := os.Create(zip_filename)
  if err != nil {
    fmt.Println("error creating zip file")
    return "", err
  }
  defer file.Close()

  _, err = io.Copy(file, res.Body)
  if err != nil {
    return "", err
  }
  defer res.Body.Close()

  return zip_filename, nil
}

func unzip(zip_filename string) error {
  r, err := zip.OpenReader(zip_filename)
  if err != nil {
    fmt.Println("error opening zip file")
    return err
  }
  defer r.Close()

  for _, f := range r.File {
    err := createTorrent(f)
    if err != nil {
      return err
    }
  }

  return nil
}


// create torrent
func createTorrent(f *zip.File) error {
  // open file, save locally
  current_file, err := f.Open()
  if err != nil {
    fmt.Println("error opening torrent file")
    return err
  }
  defer current_file.Close()

  torrent, err := os.Create("./tmp/torrents/" + f.Name)
  if err != nil {
    fmt.Println("error creating torrent file")
    return err
  }
  defer torrent.Close()

  _, err = io.Copy(torrent, current_file)
  if err != nil {
    return err
  }

  return nil
}

// acknowledge receival of torrents and delete from server
func ackTorrents() (*[]string, error) {
  files, err := ioutil.ReadDir("./tmp/torrents")
  if err != nil {
    fmt.Println("error reading ./tmp/torrents")
    return nil, err
  }
  var filenames []string

  for _, file := range files {
    if file.Name() == ".DS_Store" {
      continue
    }

    filenames = append(filenames, file.Name())
    err := deleteTorrent(file.Name())
    if err != nil {
      fmt.Println("error deleting torrent")
      return nil, err
    }
  }

  return &filenames, nil
}

func deleteTorrent(filename string) error {
  client := &http.Client{}

  url := GetBaseUrl() + "/torrents?torrent=" + filename 
  fmt.Printf("url = %s\n", url)
  req, err := http.NewRequest("DELETE", url, nil)
  if err != nil {
    return err
  }
  req.Header.Add("Authorization", GetKey())

  res, err := client.Do(req)
  if err != nil {
    fmt.Println("error processing request")
    return err
  }

  if res.StatusCode != 200 {
    fmt.Println("file not on server")
    return errors.New("Could not find file to delete")
  }

  return nil
}

func moveToWatchFolder(filenames *[]string) error {
  for _, filename := range *filenames {
    src := "./tmp/torrents/" + filename
    dest := "./torrents/" + filename
    err := os.Rename(src, dest) 
    if err != nil {
      fmt.Println("Error moving file")
      return err
    }
  }

  return nil
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
