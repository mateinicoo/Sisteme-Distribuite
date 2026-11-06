package main

import (
	"fmt"
	"sync"
)

type Msg struct {
	Value int
	Done  bool // liderul anunta ca a fost ales
}

// realizeaza operatiile interne nodului
func node(id int, in <-chan Msg, out chan<- Msg, isLeader *int, wg *sync.WaitGroup) { // am schimbat isLeader din bool in int pentru a returna in el id-ul leader-ului
	defer wg.Done()
	var mx Msg = Msg{Value: id, Done: false}

	out <- mx

	for primit := range in {
		// daca am gasit leaderul, inchidem toate canalele
		if primit.Done == true {
			mx.Done = true
			out <- mx
			close(out)
			return

		} else if mx.Value < primit.Value {
			mx = primit
			out <- mx
		} else if mx == primit {
			*isLeader = mx.Value
			mx.Done = true
			out <- mx
		}
	}

}

// returneaza id-ul liderului
func RingFlooding(ids []int) int {
	wg := sync.WaitGroup{}

	// crearea canalelor
	n := 6
	channels := make([]chan Msg, n)

	for index := range ids {
		channels[index] = make(chan Msg, 2)
	}

	var isLeader int = -1
	wg.Add(6)

	for index, val := range ids {
		go node(val, channels[index], channels[(index+1)%n], &isLeader, &wg)
	}

	wg.Wait()

	return isLeader
}

func main() {
	ids := []int{4, 6, 19, 11, 1, 3}
	Leader := RingFlooding(ids)

	fmt.Println("ID-ul leader-ului este: ", Leader)
}
