package main

import (
	"flag"
	"fmt"
)

func main() {
	urlFlag := flag.String("url", "https://gophercies.com", "the url used to build site map")
	flag.Parse()
	/*
		1.Get the web page
		2.parse all links
		3.build proper urls with our links
	*/
	fmt.Println(*urlFlag)
}
