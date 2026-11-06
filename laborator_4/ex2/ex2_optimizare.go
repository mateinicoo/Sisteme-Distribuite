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
func nodeGen(id int, inDegree int, diametru int, in <-chan MsgGen, out []chan<- MsgGen, wg *sync.WaitGroup) { // am schimbat isLeader din bool in int pentru a returna in el id-ul leader-ului
	defer wg.Done()
	var Msg MsgGen = MsgGen{Value: id, Done: false}

	for _, channel := range out {
		channel <- Msg
	}

	var i int = 0

	for i < diametru {
		i++

		for j := 0; j < inDegree; j++ {
			primit := <-in

			if primit.Value > Msg.Value {
				Msg.Value = primit.Value
			}

		}
		for _, channel := range out {
			channel <- Msg
		}
	}

	fmt.Println("Leaderul este: ", Msg.Value)

}

func RingFloodingGen(ids []int) {
	var wg sync.WaitGroup

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

	var diametru int = n + 5

	wg.Add(n)

	inDegrees := make([]int, n)
	for _, neighbors := range outTopology {
		for _, target := range neighbors {
			inDegrees[target]++
		}
	}

	for i := 0; i < n; i++ {
		var channelsOut []chan<- MsgGen
		// cream canalele de iesire
		for _, neighborId := range outTopology[i] {
			channelsOut = append(channelsOut, channelsIn[neighborId])
		}

		go nodeGen(ids[i], inDegrees[i], diametru, channelsIn[i], channelsOut, &wg)
	}

	wg.Wait()

}

func main() {
	ids := []int{3, 2, 10, 4, 1, 5, 6, 7}
	RingFloodingGen(ids)

}
