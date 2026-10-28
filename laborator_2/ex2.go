package main

import (
	"fmt"
	"sync"
)

func channel(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	mesaj := <-ch
	fmt.Println(mesaj)

}

func main() {
	var wg sync.WaitGroup
	ch := make(chan string)
	wg.Add(1)
	go channel(ch, &wg)
	ch <- "salut"

	wg.Wait()

}
