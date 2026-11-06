package main

import (
	"fmt"
	"sync"
)

// Execuția unui nod pe parcursul celor 'totalRounds' runde sincrone
func nodeSync(
	id int,
	totalRounds int,
	in []<-chan int,
	out []chan<- int,
	startRound <-chan struct{},
	endRound chan<- struct{},
	isLeader *int,
	mu *sync.Mutex,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	maxID := id

	for r := 0; r < totalRounds; r++ {
		// Așteptăm semnalul coordonatorului pentru a începe Runda r
		<-startRound

		// 1. FAZA DE TRIMITERE: Scriem maxID pe toate canalele de ieșire
		for _, chOut := range out {
			chOut <- maxID
		}

		// 2. FAZA DE PRIMIRE: Citim exact 1 mesaj de pe fiecare canal de intrare
		for _, chIn := range in {
			receivedVal := <-chIn
			if receivedVal > maxID {
				maxID = receivedVal
			}
		}

		// 3. Confirmăm coordonatorului că am încheiat Runda r
		endRound <- struct{}{}
	}

	// La finalul tuturor rundelor, nodul al cărui ID coincide cu maxID este LIDERUL!
	if maxID == id {
		mu.Lock()
		*isLeader = id
		mu.Unlock()
	}
}

func RingFloodingSync(ids []int, outTopology map[int][]int) int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	n := len(ids)

	// Canalele de rețea au buffer = 1 (suficient pentru 1 mesaj per rundă)
	channelsIn := make([]chan int, n)
	for i := 0; i < n; i++ {
		channelsIn[i] = make(chan int, 1)
	}

	// Canale pentru sincronizarea rundelor
	startRound := make(chan struct{})
	endRound := make(chan struct{})

	isLeader := -1
	wg.Add(n)

	// Pornim gorutinele pentru fiecare nod
	for i := 0; i < n; i++ {
		var channelsOut []chan<- int
		for _, neighborIdx := range outTopology[i] {
			channelsOut = append(channelsOut, channelsIn[neighborIdx])
		}

		// n runde sunt suficiente pentru ca informația să străbată orice graf cu n noduri
		go nodeSync(ids[i], n, channelsIn[i], channelsOut, startRound, endRound, &isLeader, &mu, &wg)
	}

	// ==========================================================
	// BARIERA COORDONATOARE (Sincronizatorul de Runde)
	// ==========================================================
	for r := 0; r < n; r++ {
		// Pasul A: Dăm startul Rundei r pentru TOATE cele n noduri
		for i := 0; i < n; i++ {
			startRound <- struct{}{}
		}

		// Pasul B: Așteptăm ca TOATE cele n noduri să confirme finalizarea Rundei r
		for i := 0; i < n; i++ {
			<-endRound
		}
	}

	wg.Wait()
	return isLeader
}

func main() {
	ids := []int{3, 2, 10, 4, 1, 5, 6, 7}

	// Topologia grafului conex cu două inele
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

	Leader := RingFloodingSync(ids, outTopology)
	fmt.Println("ID-ul liderului ales în mod sincron este:", Leader)
}
