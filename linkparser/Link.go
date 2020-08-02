package linkparser

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

// Link represent a link (<a href="..." > text </a>) in a HTML document
type link struct {
	href string
	text string
}

func NewLink(href, text string) *link {
	return &link{
		href: href,
		text: text,
	}
}

// Parse will take a HTML document and return a slice of links
// parsed from it
func Parse(r io.Reader) ([]*link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return findLink(root), nil
}

func findLink(n *html.Node) []*link{
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				return []*link{NewLink(a.Val, findLinkText(n))}
			}
		}
	}
	var links []*link
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
		if c.Data == "a" {
			return ""
		}
		text += findLinkText(c)
	}

	return strings.Join(strings.Fields(text), " ")
}
