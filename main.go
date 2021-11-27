package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	urlFlag := flag.String("url", "https://courses.calhoun.io/courses/cor_gophercises", "the url used to build site map")
	flag.Parse()
	/*
		1.Get the web page
		2.Parse all links
		3.Build proper urls with our links
		4.Filer out links with other domains
		5.Find all the pages(BFS)
		6.Print XML
	*/
	resp, err := http.Get(*urlFlag)
	if err != nil {
		fmt.Println(err)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	htmlFile := string(data)
	fmt.Println(htmlFile)
	r := strings.NewReader(htmlFile)
	links, _ := link.Parse(r)
	fmt.Println(links)
	// defer resp.Body.Close()

	fmt.Println(*urlFlag)
}
