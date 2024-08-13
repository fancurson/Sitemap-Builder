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
	urlFlag := flag.String("url", "https://go.dev", "the url that you want to build a sitemap")
	maxDepth := flag.Int("Depth", 3, "Depth of searching")

	sliceOfLinks := bfs(*urlFlag, *maxDepth)

	for _, el := range sliceOfLinks {
		fmt.Println(el)
	}

}

func bfs(urlFlag string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlFlag: {},
	}

	for i := 0; i < maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			 if _, ok := seen[url]; ok {
				continue
			 }
			 seen[url] = struct{}{}
			 for _, link := range get(url){
				nq[link] = struct{}{}
			 }
		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlFlag string) []string {
	// site script
	siteScript, err := http.Get(urlFlag)
	if err != nil {
		return []string{}
	}
	defer siteScript.Body.Close()

	// base = Scheme + Host
	base := fmt.Sprint(siteScript.Request.URL.Scheme + "/" + siteScript.Request.URL.Host)

	/*
		urlBody := siteScript.Request.URL
		baseUrl := &url.URL{
			Scheme: urlBody.Scheme,
			Host:   urlBody.Host,
			}
		base := baseUrl.String()
	*/

	return filter(hrefs(siteScript.Body, base), filterSetting(base))
}

func hrefs(siteBody io.Reader, base string) []string {
	links, _ := link.Parse(siteBody)
	var ret []string
	for _, el := range links {
		switch {
		case strings.HasPrefix(el.Href, "/"):
			ret = append(ret, base+el.Href)
		case strings.HasPrefix(el.Href, "http"):
			ret = append(ret, el.Href)
		}
	}
	return ret
}

func filter(sliceOfLinks []string, settingFn func(string) bool) []string {
	var ret []string
	for _, el := range sliceOfLinks {
		if settingFn(el) {
			ret = append(ret, el)
		}
	}
	return ret
}

func filterSetting(prx string) func(string) bool {
	return func(linkUrl string) bool {
		return strings.HasPrefix(linkUrl, prx)
	}
}
