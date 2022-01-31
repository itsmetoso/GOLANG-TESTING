package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func doUpload(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(rw, "404 not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		http.ServeFile(rw, r, "client.html")

	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprint(rw, "ParseForm() err: %v", err)
			return
		}

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("filename")
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()
		fmt.Fprintf(rw, "Uploaded File: %v\n", handler.Filename)
		fmt.Fprintf(rw, "File Size: %v\n", handler.Size)
		fmt.Fprintf(rw, "MIME Header: %v\n", handler.Header)

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		text := string(fileBytes)

		resp, err := http.PostForm("http://localhost:8080", url.Values{"words": {text}})

		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()

		//Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)
		//log.Printf(sb)
		fmt.Fprintf(rw, sb)

	default:
		fmt.Fprint(rw, "Sorry, only GET and POST methods are supported")
	}

}

func main() {
	http.HandleFunc("/", doUpload)

	fmt.Printf("Starting client for testing HTTP POST ... \n")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
