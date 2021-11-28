package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	linkparser "github.com/chethan0707/linkparser"
)

const xlmns = "http://www.sitemaps.org/schemas/sitemap/0.9"

/*
	1.Get the web page
	2.Parse all links
	3.Build proper urls with our links
	4.Filer out links with other domains
	5.Find all the pages(BFS)
	6.Print XML
*/
type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url used to build site map")
	maxDepth := flag.Int("depth", 3, "Maximum depth to traverse")
	flag.Parse()
	finalLinks := bfs(*urlFlag, *maxDepth)
	// for _, link := range finalLinks {
	// 	fmt.Println(link)
	// }

	// pages := get(*urlFlag)
	// for _, page := range pages {
	// 	fmt.Println(page)
	// }
	toXML := urlset{
		Xmlns: xlmns,
	}
	for _, page := range finalLinks {
		toXML.Urls = append(toXML.Urls, loc{page})
	}
	fmt.Println(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()

}
func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		fmt.Println("i = ", i)
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				// if _, ok := seen[link]; !ok {
				nq[link] = struct{}{}
				// }
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println(err)
	}

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	defer resp.Body.Close()
	base := baseURL.String()
	// fmt.Println(reqURL)
	// fmt.Println(base)
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

func hrefs(r io.Reader, base string) []string {
	var ret []string
	pages, _ := linkparser.Parse(r)
	for _, l := range pages {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}
