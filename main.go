package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Go")
		d, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "ooops", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Hello %s", d)
	})

	http.HandleFunc("/hello", func(http.ResponseWriter, *http.Request) {
		log.Println("Hello Go")
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request) {
		log.Println("Goodbye Go")
	})

	http.ListenAndServe(":9090", nil)
}
