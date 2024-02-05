package main 

import (
  "encoding/json"
  "fmt"
  "io"
  "os"
  "log"
  "net/http"
  "regexp"
)

func main() {
  text := text()

  f, err := os.OpenFile("README.md", os.O_RDWR,0644)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  fInfo, err := f.Stat()

  if err != nil {
    fmt.Println("Error getting file information:", err)
    return
  }

  b_in := make([]byte, fInfo.Size())

  _,err = f.Read(b_in)
  if err != nil {
    log.Fatal(err)
  }

  var re = regexp.MustCompile(`(<!-- <RandomFunFact> -->)[\s\S]*(<!-- </RandomFunFact> -->)`);
  b_out := re.ReplaceAll(b_in, []byte("$1 " + text + " $2"))

  err = f.Truncate(0)
  if err != nil {
    log.Fatal(err)
  }

  _, err = f.Seek(0, 0)
  if err != nil {
    log.Fatal(err)
  }

  _,err = f.Write(b_out)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("File successfully updated.")
}

func text() string {
  res, err := http.Get("https://uselessfacts.jsph.pl/api/v2/facts/random")
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if res.StatusCode != http.StatusOK {
    log.Fatal("API Request failed with status:", res.Status)
  }

  body, err := io.ReadAll(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  var data map[string]interface{}
  if err := json.Unmarshal([]byte(body), &data); err != nil {
    log.Fatal(err)
  }

  return data["text"].(string)
}
