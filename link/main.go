package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	bytes, err := ioutil.ReadFile("./link/ex3.html")
	if err != nil {
		log.Println("Error reading file.")
		return
	}
	links := parseHTML(bytes)
	printLinks(links)
}

func printLinks(links []Link) {
	for i, link := range links {
		fmt.Printf("Link{\n\tHref: \"%v\",\n\tText: \"%v\",\n}", link.Href, link.Text)
		if i != len(links) - 1 {
			fmt.Print(",")
		}
		fmt.Print("\n")
	}
}