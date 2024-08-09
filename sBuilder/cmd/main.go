package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap")

	fmt.Println(*urlFlag)
	r, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	io.Copy(os.Stdout, r.Body)

	// data, _ := sbuilder.Parse(r.Body)
	// _ = data

}
