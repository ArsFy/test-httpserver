package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func o3[T string | int](b bool, v1 T, v2 T) T {
	if b {
		return v1
	} else {
		return v2
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// r.ParseMultipartForm(4096)

	var header, form string
	header = fmt.Sprintf("Host: %s\n", r.Host)
	for i, j := range r.Header {
		header += fmt.Sprintf("%s: %s\n", i, strings.Join(j, ""))
	}
	for i, j := range r.Form {
		form += fmt.Sprintf("%s: %s\n", i, strings.Join(j, ""))
	}

	file, handler, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
	}

	fmt.Println()
	fmt.Println(r.RemoteAddr, r.Method, r.URL.Path, r.Proto, o3(header == "", "", "\n"))
	fmt.Println(header)
	fmt.Print(form, o3(form == "", "", "\n"))
	if err == nil {
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		hb, _ := json.Marshal(handler.Header)
		fmt.Printf("MIME Header: %s\n", string(hb))
		tempFile, _ := ioutil.TempFile("temp", handler.Filename)
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		tempFile.Write(fileBytes)
		fmt.Println()
	}
	fmt.Println("--------------------------------------")
	fmt.Fprintf(w, "ok")
}

func main() {
	http.HandleFunc("/", indexHandler)

	fmt.Println("Server Starting...")
	http.ListenAndServe(":8000", nil)
}
