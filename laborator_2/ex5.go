package main

import (
	"fmt"
)

func send(ch chan<- string) {
	ch <- "salut"
}

func receive(ch <-chan string) string {
	mesaj := <-ch
	return mesaj
}

func main() {
	ch := make(chan string)
	go send(ch)

	rezultat := receive(ch)

	fmt.Println(rezultat)
}
