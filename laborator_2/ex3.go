package main

import "fmt"

func channel(ch chan<- string) {
	ch <- "mesaj1"
	ch <- "mesaj2"

}

func main() {
	ch := make(chan string, 2)

	go channel(ch)

	mesajPrimit1 := <-ch
	mesajPrimit2 := <-ch
	fmt.Println(mesajPrimit1)
	fmt.Println(mesajPrimit2)
}
