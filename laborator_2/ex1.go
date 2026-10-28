package main

import "fmt"

func channel(ch chan<- string) {
	ch <- "salut"
}

func main() {
	ch := make(chan string)

	go channel(ch)

	mesaj := <-ch
	fmt.Println(mesaj)
}
