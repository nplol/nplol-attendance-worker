package main

import (
  "net/http"
  "fmt"
  "os"
  "io/ioutil"
  "encoding/json"
  "time"
)

type Attendance struct {
    Username  string        `json:"username"`
    LastSeen  time.Time     `json:"lastSeen"`
}

func main() {
  args := os.Args[1:]

  url := args[0]

  resp, err := http.Get(url)
  if (err != nil) {
    os.Exit(1)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if (err != nil) {
    os.Exit(1)
  }

  var dat []Attendance

  if err := json.Unmarshal(body, &dat); err != nil {
    fmt.Printf("Could not parse data")
    panic(err)
  }

  fmt.Println(dat)
}
