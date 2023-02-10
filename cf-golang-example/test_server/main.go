package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}




func redirectResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Location", "https://cf-golang-server.cfapps-06.slot-59.pez.vmware.com")
				w.Header().Set("Set-Cookie", "YYYYY0=xxx; Max-Age=28800")
				w.WriteHeader(http.StatusFound)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.WriteHeader(http.StatusFound)
				fmt.Fprint(w, string("Response from GET"))
			}
		})
}

func normalResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from GET"))
			}
		})
}

func largeResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from GET"))
			}
			io.WriteString(w, randStringBytes(10240000))
		})
}

func delayResponseHandler() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			second, _ := strconv.Atoi(os.Getenv("DELAY_SECONDS"))
			time.Sleep(time.Second * time.Duration(second))
			switch r.Method {
			case http.MethodPost:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from POST"))
			case http.MethodGet:
				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, string("Response from GET"))
			}
		})
}

func setupGracefulShutdown() {
	cancelChan := make(chan os.Signal, 1)
	// done := make(chan bool, 1)

	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-cancelChan
		log.Printf("Caught SIGTERM %v", sig)
		// done <- true
	}()
	// <-done
}

func main() {
	setupGracefulShutdown()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This is home page!")
	})
	http.HandleFunc("/redirect", redirectResponseHandler())

	http.HandleFunc("/normal", normalResponseHandler())
	http.HandleFunc("/large", largeResponseHandler())
	http.HandleFunc("/delay", delayResponseHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
