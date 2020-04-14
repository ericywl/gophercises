package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ericywl/gophercises/recover/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(":3000", middleware.RecoverFromPanic(mux)))
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "<h1>Hello!</h1>")
	if err != nil {
		log.Fatal(err)
	}
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "<h1>Hello!</h1>")
	if err != nil {
		log.Fatal(err)
	}
}
