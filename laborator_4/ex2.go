package main

import (
	"fmt"
	"sync"
)

type MsgGen struct {
	Value int
	Done  bool // liderul anunta ca a fost ales
}

// realizeaza operatiile interne nodului
func nodeGen(id int, in <-chan MsgGen, out []chan<- MsgGen, isLeader *int, wg *sync.WaitGroup, mu *sync.Mutex) { // am schimbat isLeader din bool in int pentru a returna in el id-ul leader-ului
	defer wg.Done()
	var initialMsg MsgGen = MsgGen{Value: id, Done: false}

	for _, channel := range out {
		channel <- initialMsg
	}

	maxValue := id
	terminated := false

	for msg := range in {
		if msg.Done == true {
			if !terminated {
				terminated = true

				for _, channel := range out {
					channel <- msg
				}
			}

			return

		} else if msg.Value > maxValue {
			maxValue = msg.Value

			for _, channel := range out {
				channel <- msg
			}

		} else if msg.Value == id && !terminated {
			terminated = true
			mu.Lock()
			*isLeader = msg.Value
			mu.Unlock()

			msg.Done = true

			for _, channel := range out {
				channel <- msg
			}

			return

		}
	}

	return

}

// returneaza id-ul liderului
func RingFloodingGen(ids []int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex

	n := len(ids)

	// cream canalele de intrare
	channelsIn := make([]chan MsgGen, n)
	for i := 0; i < n; i++ {
		channelsIn[i] = make(chan MsgGen, 50)
	}

	outTopology := map[int][]int{
		0: {1, 4},
		1: {2},
		2: {3},
		3: {0},
		4: {5},
		5: {6},
		6: {7},
		7: {0},
	}

	isLeader := -1
	wg.Add(n)

	for i := 0; i < n; i++ {
		var channelsOut []chan<- MsgGen
		// cream canalele de iesire
		for _, neighborId := range outTopology[i] {
			channelsOut = append(channelsOut, channelsIn[neighborId])
		}

		go nodeGen(ids[i], channelsIn[i], channelsOut, &isLeader, &wg, &mu)
	}

	wg.Wait()
	return isLeader

}

func main() {
	ids := []int{3, 2, 10, 4, 1, 5, 6, 7}
	Leader := RingFloodingGen(ids)

	fmt.Println("ID-ul leader-ului este: ", Leader)
}
