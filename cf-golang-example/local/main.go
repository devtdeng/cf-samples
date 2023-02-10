package main

import (
    "log"
    "net/http"
)

func main() {
    client := &http.Client{
    }
    _, err := client.Get("http://spring-music.cfapps-06.slot-59.pez.vmware.com")
    req, err := http.NewRequest("GET", "http://spring-music.cfapps-06.slot-59.pez.vmware.com", nil)
    req.Header.Add("X-Test-header", "Test header value\r")
    _, err = client.Do(req)
    if err != nil {
      log.Fatalln(err)
    }
}
