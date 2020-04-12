package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ericywl/gophercises/linkparser/link"
)

func main() {
	bytes, err := ioutil.ReadFile("./linkparser/ex3.html")
	if err != nil {
		log.Println("Error reading file.")
		return
	}
	links := link.ParseHTML(bytes)
	printLinks(links)
}

func printLinks(links []link.Link) {
	for i, l := range links {
		fmt.Printf("Link{\n\tHref: \"%v\",\n\tText: \"%v\",\n}", l.Href, l.Text)
		if i != len(links) - 1 {
			fmt.Print(",")
		}
		fmt.Print("\n")
	}
}
