package main

import (
	"fmt"
	"sync"
	"time"
)

type Coordonata struct {
	X, Y int
}

func robot(id int, traseu []Coordonata, out chan<- Coordonata, in <-chan Coordonata, wgObserver *sync.WaitGroup) {
	defer wgObserver.Done()

	var wgCitire sync.WaitGroup
	wgCitire.Add(1)

	go func() {
		defer wgCitire.Done()
		for poz := range in {
			fmt.Printf("Robotul %d a primit pozitia (%d, %d) de la celalalt robot.\n", id, poz.X, poz.Y)
		}
		fmt.Printf("Robotul %d: Celalalt robot s-a oprit\n", id)
	}()

	for _, poz := range traseu {
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("Robotul %d a ajuns la (%d, %d)\n", id, poz.X, poz.Y)
		out <- poz
	}

	fmt.Printf("Robotul %d a terminat traseul\n", id)
	close(out)

	wgCitire.Wait()
}

func main() {

	chRob1ToRob2 := make(chan Coordonata)
	chRob2ToRob1 := make(chan Coordonata)

	var wgObserver sync.WaitGroup
	wgObserver.Add(2)

	traseu1 := []Coordonata{{0, 1}, {0, 2}, {0, 3}}
	traseu2 := []Coordonata{{5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}}

	fmt.Println("Robotii incep")
	go robot(1, traseu1, chRob1ToRob2, chRob2ToRob1, &wgObserver)

	go robot(2, traseu2, chRob2ToRob1, chRob1ToRob2, &wgObserver)

	wgObserver.Wait()
	fmt.Println("Final cursa")
}
