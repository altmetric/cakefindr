package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	var backends []string

	if len(os.Getenv("BACKENDS")) > 0 {
		backends = strings.Split(strings.Trim(os.Getenv("BACKENDS"), ","), ",")
	} else {
		panic("No backend servers specified")
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte("PONG"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)

		backend := backends[rand.Int()%len(backends)]
		url := fmt.Sprintf("http://%s", backend)
		fmt.Printf("Requesting payload from backend URL=%s\n", url)

		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		text, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		w.Write(text)
	})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
