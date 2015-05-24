package main

import (
  "net/http"
  "fmt"
  "os"
  "io/ioutil"
  "encoding/json"
  "time"
  "strconv"
  "bytes"
)

type Attendance struct {
    Username  string        `json:"username"`
    LastSeen  time.Time     `json:"lastSeen"`
}

type AttendanceDays struct {
  Username    string        `json:"username"`
  DaysSince   int64         `json:"daysSince"`
  AlertLevel  int           `json:"alertLevel"`
}

func NewAttendanceDays(att Attendance) *AttendanceDays {
  a := new(AttendanceDays)
  a.Username = att.Username
  a.AlertLevel = 0
  now := time.Now()
  diff := now.Sub(att.LastSeen)
  days := diff / (24 * time.Hour)
  a.DaysSince = int64(days)
  return a
}

func main() {
  args := os.Args[1:]

  producerUrl := args[0]
  consumerUrl := args[1]

  warningLevel, err := strconv.ParseInt(args[2], 10, 64)
  if (err != nil) {
    fmt.Println("Could not convert " + args[2] + " to int")
    os.Exit(1)
  }

  dangerLevel, err := strconv.ParseInt(args[3], 10, 64)
  if (err != nil) {
    fmt.Println("Could not convert " + args[3] + " to int")
    os.Exit(1)
  }

  resp, err := http.Get(producerUrl)
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

  var dangerZone []AttendanceDays

  for _,att := range dat {
      attDays := *NewAttendanceDays(att)
      if attDays.DaysSince >= warningLevel{
        attDays.AlertLevel = 1

        if attDays.DaysSince >= dangerLevel {attDays.AlertLevel = 2}

        dangerZone = append(dangerZone, attDays)
      }
  }

  marshalled, _ := json.Marshal(dangerZone)
  reader := bytes.NewReader(marshalled)
  http.Post(consumerUrl, "application/json", reader)
  fmt.Println(string(marshalled))
}
