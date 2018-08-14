package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte("PONG"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte("Hello world!"))
	})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
