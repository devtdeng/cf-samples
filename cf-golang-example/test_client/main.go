package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// refer to test code at https://github.com/golang/go/issues/34902
// https://gist.github.com/KauzClay/2d84e0ce1d1884a7bf3ca952ebb57630
func main() {
	var n int
	// 10 goroutines is a good number to see the panic happen quickly on
	// our test environments, though this may not apply to all machines
	flag.IntVar(&n, "n", 1, "number of goroutines")
	flag.Parse()

	// Request needs a larger body to see the error happen more quickly
	// It is reproducable with smaller body sizes, but it takes longer to fail
	bodyString := strings.Repeat("a", 2048)
	client := &http.Client{}

	for i := 0; i < n; i++ {
		go func() {
			for {
				buf := bytes.NewBufferString(bodyString)
				req, _ := http.NewRequest("POST", os.Getenv("SERVER_URL"), buf)
				req.Header.Add("Expect", "100-continue")

				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Request Failed: %s\n", err.Error())
				} else {
					resp.Body.Close()
				}
			}
		}()
	}

	select {}
}
