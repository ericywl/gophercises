package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"./urlshort"
)

func main() {
	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// Parse flags
	filePtr := flag.String("file", "", "file to load URL mapping from")
	flag.Parse()
	if *filePtr == "" {
		flag.Usage()
		return
	}
	// Read file
	contents, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		panic(err)
	}
	// Build the handler using the mapHandler as the fallback
	var handler http.HandlerFunc
	fileExt := filepath.Ext(*filePtr)
	switch fileExt {
	case ".yaml":
		handler, err = urlshort.YAMLHandler([]byte(contents), mapHandler)
	case ".json":
		handler, err = urlshort.JSONHandler([]byte(contents), mapHandler)
	default:
		fmt.Println("Unsupported file format", fileExt)
		return
	}

	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", handler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
