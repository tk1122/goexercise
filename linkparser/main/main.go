package main

import (
	"fmt"
	"goercises/linkparser"
	"log"
	"strings"
)

var exampleHtml = `
<html>
	<body>
	<h1>Hello!</h1>
	<a href="/other-page">
		A link to another page
		<a href="/bad">
			A bad link
		</a>
	</a>
	<div>
	<a href="/another-page">
		Outside span ...
		<span>
			Inside span
			<a href="/another-bad">
				A another bad link
			</a>
		</span>
	</a>
	</div>
	</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml)
	links, err := linkparser.Parse(r)
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range links {
		fmt.Printf("%+v\n", *l)
	}
}
