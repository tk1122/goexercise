package linkparser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func NewLink(href, text string) *Link {
	return &Link{
		Href: href,
		Text: text,
	}
}

// Parse will take a HTML document and return a slice of links
// parsed from it
func Parse(r io.Reader) ([]*Link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return findLink(root), nil
}

func findLink(n *html.Node) []*Link {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				return []*Link{NewLink(a.Val, findLinkText(n))}
			}
		}
	}
	var links []*Link
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, findLink(c)...)
	}
	return links
}

func findLinkText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	if n.Type != html.ElementNode {
		return ""
	}

	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += findLinkText(c)
	}

	return strings.Join(strings.Fields(text), " ")
}
