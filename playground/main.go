package main

import (
	"fmt"
)

func main() {
	const cap = 5
	ch := make(chan string, cap)

	go func() {
		for p := range ch {
			fmt.Println("employee : received :", p)
		}
	}()

	const work = 20
	for w := 0; w < work; w++ {
		select {
		case ch <- "paper":
			fmt.Println("manager : send ack")
		default:
			fmt.Println("manager : drop")
		}
	}

	close(ch)
}