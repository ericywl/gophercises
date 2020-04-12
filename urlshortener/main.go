package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"

	"github.com/ericywl/gophercises/urlshortener/urlshort"
)

func main() {
	mux := defaultMux()
	// Parse flags
	filePtr := flag.String("file", "", "file to load URL mapping from")
	dbPtr := flag.String("db", "", "BoltDB file for fallback URL mapping")
	flag.Parse()
	if *filePtr == "" {
		flag.Usage()
		return
	}
	// Set up fallback handler
	var fallbackHandler http.HandlerFunc
	if *dbPtr != "" {
		// Check if database file exist
		if _, err := os.Stat(*dbPtr); os.IsNotExist(err) {
			panic(fmt.Errorf("%s does not exist", *dbPtr))
		} else if err != nil {
			panic(err)
		}
		// Open database file
		db, err := bolt.Open(*dbPtr, 0600, nil)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		fallbackHandler = urlshort.DBHandler(db, []byte("URLMap"), mux)
	} else {
		fallbackHandler = defaultFallback(mux)
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
		handler, err = urlshort.YAMLHandler([]byte(contents), fallbackHandler)
	case ".json":
		handler, err = urlshort.JSONHandler([]byte(contents), fallbackHandler)
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

func defaultFallback(mux *http.ServeMux) http.HandlerFunc {
	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	return urlshort.MapHandler(pathsToUrls, mux)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
