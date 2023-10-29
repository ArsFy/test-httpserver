package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
		tempFile, _ := os.CreateTemp("temp", handler.Filename)
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		fileBytes, err := io.ReadAll(file)
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

	go func() {
		_, crt := os.Stat("server.crt")
		_, key := os.Stat("server.key")
		if crt == nil && key == nil {
			cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
			if err != nil {
				panic(err)
			}
			server := &http.Server{
				Addr:      ":8443",
				TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}},
			}
			server.ListenAndServeTLS("", "")
		}
	}()

	http.ListenAndServe(":8000", nil)
}
