package main

import (
	"fmt"
	"log"

	"github.com/ericywl/gophercises/sitemapbuilder/sitemap"
)

func main() {
	output, err := sitemap.BuildSiteMap("http://calhoun.io/")
	if err != nil {
		log.Printf("Error building sitemap: %v", err)
		return
	}

	fmt.Println(output)
}
