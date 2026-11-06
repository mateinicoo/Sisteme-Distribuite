package main

import (
	"fmt"
	"math"
	"sync"
)

// pentru aceasta problema am transmis si numarul de noduri de intrare al fiecarui nod ca parametru

type Barrier struct {
	n          int
	count      int
	generation int
	mutex      sync.Mutex
	cond       *sync.Cond
}

func NewBarrier(n int) *Barrier {
	b := &Barrier{n: n}
	b.cond = sync.NewCond(&b.mutex)
	return b
}

func (b *Barrier) Wait() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	gen := b.generation
	b.count++

	if b.count == b.n {
		b.count = 0
		b.generation++
		b.cond.Broadcast()
	} else {
		for gen == b.generation {
			b.cond.Wait()
		}
	}
}

func node(vot int, inDegrees int, in <-chan float64, out []chan<- float64, wg *sync.WaitGroup, barrier *Barrier) {
	defer wg.Done()

	var avg float64
	avg = float64(vot)

	for _, channel := range out {
		channel <- float64(vot)
	}
	i := 0
	for i < 70 {
		i++
		var sum float64 = avg // propria valoare
		for j := 0; j < inDegrees; j++ {
			msg := <-in
			sum = sum + msg
		}
		barrier.Wait()

		avg = sum / float64(inDegrees+1)

		for _, channel := range out {
			channel <- float64(avg)
		}
	}

	votFinal := math.Floor(avg + 0.5)
	fmt.Println("Valoare vot:", votFinal)

}

func MajorityVoting(votes []int, outTopology map[int][]int) {
	var wg sync.WaitGroup
	n := len(votes)
	channelsIn := make([]chan float64, n)
	for v := range votes {
		channelsIn[v] = make(chan float64, 50)
	}
	inDegrees := make([]int, n)
	for _, neighbors := range outTopology {
		for _, target := range neighbors {
			inDegrees[target]++
		}
	}
	barrier := NewBarrier(n)
	wg.Add(n)

	for v := 0; v < n; v++ {
		var channelsOut []chan<- float64
		for _, neighbor := range outTopology[v] {
			channelsOut = append(channelsOut, channelsIn[neighbor])
		}
		go node(votes[v], inDegrees[v], channelsIn[v], channelsOut, &wg, barrier)
	}

	wg.Wait()
}

func main() {
	votes := []int{1, 1, 0, 0, 0, 1, 0, 0}
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
	MajorityVoting(votes, outTopology)
}
