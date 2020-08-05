package sitemapbuilder

import (
	"goercises/linkparser"
	"net/http"
	"strings"
)

type node struct {
	Href string
	Children []*node
}

/*
	node { href: / children:
		[
			node {href: /child-1, children:
				[
					node {href: /child-1/grand-1, children: nil}
				]
			}
			node {href: /child-2, children: nil}
		]
	}
 */

func ParseSiteMap(url, root string) *node{
	n := &node{}
	links := Links(url, root)
	for _, l := range links {
		n.Children = append(n.Children, ParseSiteMap(l.Href, root))
	}

	return n
}

func Links(url string, root string) []*linkparser.Link{
	res, err := http.Get(root + url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	links, err := linkparser.Parse(res.Body)
	if err != nil {
		return nil
	}

	var filLinks []*linkparser.Link
	for _, l := range links {
		if strings.HasPrefix(l.Href, "/") {
			filLinks = append(filLinks, l)
		}
	}

	return filLinks
}