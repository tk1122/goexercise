package main

import (
	"fmt"
	"goercises/sitemapbuilder"
)

func main() {
	//n := sitemapbuilder.Links("/", "https://golang.org")
	//
	//for _, l := range n {
	//	fmt.Printf("%+v\n", l)
	//}

	result := sitemapbuilder.ParseSiteMap("/", "https://golang.org")

	fmt.Printf("%+v\n", result)
}
