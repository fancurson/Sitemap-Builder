package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/fancurson/sitemap/sBuilder/link"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap")

	// get sorted slice of links
	links := get(*urlFlag)

	for _, el := range links {
		fmt.Println(el)
	}

}

func get(requestUrl string) []string {
	// site script
	r, err := http.Get(requestUrl)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// base = Scheme + Host
	base := fmt.Sprint(r.Request.URL.Scheme + "/" + r.Request.URL.Host)

	/*
		urlBody := r.Request.URL
		baseUrl := &url.URL{
			Scheme: urlBody.Scheme,
			Host:   urlBody.Host,
			}
		base := baseUrl.String()
	*/

	return filter(hrefs(r.Body, base), giveMeTrue(base))
}

func hrefs(r io.Reader, urlStr string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, el := range links {
		switch {
		case strings.HasPrefix(el.Href, "/"):
			ret = append(ret, urlStr+el.Href)
		case strings.HasPrefix(el.Href, "http"):
			ret = append(ret, el.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, el := range links {
		if keepFn(el) {
			ret = append(ret, el)
		}
	}
	return ret
}

func giveMeTrue(prx string) func(string) bool {
	return func(linkUrl string) bool {
		if strings.HasPrefix(linkUrl, prx) {
			return true
		}
		return false
	}
}
