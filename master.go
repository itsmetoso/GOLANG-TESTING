package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)

	m := make(map[string]int)

	for _, word := range words {
		m[word] += 1
	}

	return m
}

func doParse(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(rw, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(rw, r, "master.html")

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(rw, "ParseForm() err: %v", err)
			return
		}

		data := r.FormValue("words")
		fmt.Fprintf(rw, "Words = %s\n", data)
		fmt.Fprintf(rw, "Top ten = %v", WordCount(data))
	default:
		fmt.Fprint(rw, "Sorry, only GET and POST methods are supported")
	}

}

func main() {
	http.HandleFunc("/", doParse)

	fmt.Printf("Starting server for testing HTTP POST ... \n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
