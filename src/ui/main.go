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

var Secret = []byte{3, 22, 14, 39, 91, 88, 76, 62, 39, 123, 86, 82, 83, 13, 93, 70, 25, 74, 72, 88, 94, 92, 77, 13, 42, 91, 86, 85, 56, 10, 102, 33, 18, 88, 94, 88, 72, 20, 66, 95, 97, 35, 86, 55, 20, 94, 88, 94, 35, 100, 95, 93, 33, 81, 66, 13, 86, 92, 68, 78, 121, 92, 93, 69, 25, 124, 49, 20, 118, 42, 94, 93, 18, 73, 90, 85, 20, 83, 88, 44, 78, 75, 12}

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
