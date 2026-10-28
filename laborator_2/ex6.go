package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func rob1(ch chan<- string) {
	timp := rand.IntN(100)
	time.Sleep(time.Duration(timp) * time.Millisecond)
	ch <- "finishROBOT1"
}
func rob2(ch chan<- string) {
	timp := rand.IntN(100)
	time.Sleep(time.Duration(timp) * time.Millisecond)
	ch <- "finishROBOT2"
}

func main() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go rob1(ch1)
	go rob2(ch2)

	select {
	case msg := <-ch1:
		fmt.Println("Robotul 1 a terminat primul cu mesajul: ", msg)
	case msg := <-ch2:
		fmt.Println("Robotul 2 a terminat primul cu mesajul: ", msg)
	case <-time.After(200 * time.Millisecond):
		fmt.Println("Cursa s-a terminat fara un castigator")
	}
}
