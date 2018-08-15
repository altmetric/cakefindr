package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var text string
	if len(os.Getenv("SECRET_KEY")) > 0 {
		text = os.Getenv("SECRET_KEY")
	} else {
		panic("SECRET_KEY was not found")
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte("PONG"))
	})

	http.HandleFunc("/key-3", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte(text))
	})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
