package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func Crypt(input string, key string) (output string) {
	for i := range input {
		output += string(input[i] ^ key[i%len(key)])
	}

	return output
}

var Secret = []byte{32, 13, 26, 7, 17, 68, 7, 69, 18, 84, 70, 72, 7, 17, 83, 25, 86, 7, 22, 18, 19, 84}

func main() {
	keyserver1 := os.Getenv("KEYSERVER_1")
	keyserver2 := os.Getenv("KEYSERVER_2")
	keyserver3 := os.Getenv("KEYSERVER_3")

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		w.Write([]byte("PONG"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		url1 := fmt.Sprintf("http://%s/key-1", keyserver1)
		url2 := fmt.Sprintf("http://%s/key-2", keyserver2)
		url3 := fmt.Sprintf("http://%s/key-3", keyserver3)
		fmt.Printf("Loading keys from %s, %s, %s\n", url1, url2, url3)

		text1 := []byte("")
		text2 := []byte("")
		text3 := []byte("")

		if res1, err := http.Get(url1); err == nil {
			text1, _ = ioutil.ReadAll(res1.Body)
		}
		if res2, err := http.Get(url2); err == nil {
			text2, _ = ioutil.ReadAll(res2.Body)
		}
		if res3, err := http.Get(url3); err == nil {
			text3, _ = ioutil.ReadAll(res3.Body)
		}

		fmt.Printf("Got keys %s, %s, %s\n", text1, text2, text3)
		key := fmt.Sprintf("%s-%s-%s", text1, text2, text3)

		decrypted := Crypt(string(Secret), key)

		w.Write([]byte(fmt.Sprintf("Decoded message: %s", decrypted)))
	})

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
