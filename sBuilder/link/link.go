package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Links struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Links, error) {
	doc, err := html.Parse(r)
	if err != nil {
		panic(err)
	}
	nodes := linkNodes(doc)
	links := []Links{}
	for _, val := range nodes {
		links = append(links, buildLink(val))
	}

	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Links {
	var ret Links

	for _, att := range n.Attr {
		if att.Key == "href" {
			ret.Href = att.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}

	return strings.Join(strings.Fields(ret), " ")
}
