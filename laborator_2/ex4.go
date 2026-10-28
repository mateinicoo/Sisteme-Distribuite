package main

import "fmt"

func channel(ch chan<- string) {
	ch <- "salut"
	close(ch)
}

func main() {
	ch := make(chan string)

	go channel(ch)

	mesaj1, ok1 := <-ch
	fmt.Println(mesaj1, ok1)
	mesaj2, ok2 := <-ch
	fmt.Println(mesaj2, ok2)
}
