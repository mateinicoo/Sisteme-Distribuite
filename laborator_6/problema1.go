package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// am ales ca nodurile esuate sa fie alese aleatoriu, in functia "crashNode"
// am schimbat functia "diametru" pentru a ignora nodurile esuate cand calculeaza diametru

func diametru(adj [][]int, crashed []int) int {
	n := len(adj)
	diam := 0
	for src := 0; src < n; src++ {
		// BFS din nodul src
		dist := make([]int, n)
		for i := range dist {
			dist[i] = -1
		}
		dist[src] = 0
		queue := []int{src}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					queue = append(queue, v)
				}
			}
		}
		if slices.Contains(crashed, src) == false {
			for v, d := range dist {
				if slices.Contains(crashed, v) {
					continue
				}
				if d == -1 {
					return -1 // avem un nod la care nu putem ajunge din src
				}
				if d > diam {
					diam = d
				}
			}
		}
	}
	return diam
}

// Info reprezinta o pereche (id, valoare) retinuta in memoria unui nod.
type Info struct {
	ID    int
	Value int
}

// Memory este o multime de token -uri (reuniunea se face prin inserare).
type Memory map[Info]struct{}

func (m Memory) Union(other Memory) {
	for info := range other {
		m[info] = struct{}{}
	}
}

func (m Memory) Clone() Memory {
	cp := make(Memory, len(m))
	for info := range m {
		cp[info] = struct{}{}
	}
	return cp
}

type Node struct {
	ID        int
	Neighbors []int
	M         Memory
}

// Structura grafului:

type GraphConfig struct {
	N      int
	S      int
	Tokens []int
	Adj    [][]int
}

// Citire:

func ReadGraphConfig(filename string) (*GraphConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("eroare la deschiderea fișierului: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Helper pentru a citi urmatoarea linie si a o sparge dupa spatii albe
	nextFields := func() ([]string, error) {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("fisier incomplet (EOF neasteptat)")
		}
		return strings.Fields(scanner.Text()), nil
	}

	// 1. Linia 1: n si s
	line1, err := nextFields()
	if err != nil {
		return nil, err
	}
	if len(line1) < 2 {
		return nil, fmt.Errorf("linia 1 trebuie sa contina n si s")
	}
	n, err := strconv.Atoi(line1[0])
	if err != nil {
		return nil, fmt.Errorf("n invalid: %w", err)
	}
	s, err := strconv.Atoi(line1[1])
	if err != nil {
		return nil, fmt.Errorf("s invalid: %w", err)
	}

	// 2. Linia 2: n valori pentru token-uri
	line2, err := nextFields()
	if err != nil {
		return nil, err
	}
	if len(line2) != n {
		return nil, fmt.Errorf("asteptat %d token-uri pe linia 2, gasit %d", n, len(line2))
	}
	tokens := make([]int, n)
	for i := 0; i < n; i++ {
		val, err := strconv.Atoi(line2[i])
		if err != nil {
			return nil, fmt.Errorf("token invalid la indexul %d: %w", i, err)
		}
		tokens[i] = val
	}

	// 3. Urmatoarele n linii: lista de adiacenta a nodului i
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		fields, err := nextFields()
		if err != nil {
			return nil, fmt.Errorf("eroare la citirea vecinilor nodului %d: %w", i, err)
		}

		adj[i] = make([]int, 0, len(fields))
		for _, f := range fields {
			neighbor, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("vecin invalid '%s' pentru nodul %d: %w", f, i, err)
			}
			adj[i] = append(adj[i], neighbor)
		}
	}

	return &GraphConfig{
		N:      n,
		S:      s,
		Tokens: tokens,
		Adj:    adj,
	}, nil
}

// Calculam cate noduri intra intr-un anumit nod

func inCount(nodeID int, graph GraphConfig) int {
	var cnt int = 0
	for i := range graph.N {
		if slices.Contains(graph.Adj[i], nodeID) {
			cnt++
		}
	}
	return cnt
}

func crashNode(graph *GraphConfig) []int {
	s := graph.S
	crashedNodes := []int{}
	for s > 0 { // simulam esuarea aleatorie a nodurilor
		crash := rand.IntN(graph.N - 1)

		if slices.Contains(crashedNodes, crash) == false {
			s--
			crashedNodes = append(crashedNodes, crash)

			graph.Adj[crash] = nil
		}
	}

	for _, node := range crashedNodes {
		for i := range graph.Adj {
			val := node

			graph.Adj[i] = slices.DeleteFunc(graph.Adj[i], func(x int) bool {
				return x == val
			})
		}
	}

	fmt.Println("Au esuat ", graph.S, " noduri: ")
	for _, node := range crashedNodes {
		fmt.Println(node)
	}

	return crashedNodes
}

func goNode(node Node, leader int, nodeCnt int, iterations int, inCnt int, in <-chan Memory, out []chan<- Memory, wg *sync.WaitGroup) {
	defer wg.Done()

	for range iterations {
		for _, channelOut := range out {
			channelOut <- node.M.Clone()
		}

		for range inCnt {
			otherMemory := <-in
			node.M.Union(otherMemory)
		}

	}

	if node.ID == leader { // doar un nod afiseaza valoarea de consens
		fmt.Printf("Nodul %d raporteaza:", node.ID)
		fmt.Printf("\n")

		msgCnt := len(node.M)

		if msgCnt == nodeCnt {
			fmt.Println("Valoarea de consens a fost atinsa.")
			fmt.Println("Numar de mesaje primite: ", msgCnt)
		} else {
			fmt.Println("Valoarea de consens NU a fost atinsa.")

			fmt.Println("Au fost primite doar ", msgCnt, " mesaje")
		}

		for info := range node.M {
			fmt.Printf("(ID:%d, Val:%d) ", info.ID, info.Value)
			fmt.Printf("\n")
		}
	}

}

func FloodSetAlgorithm(graph GraphConfig) {
	var wg sync.WaitGroup

	crashedNodes := crashNode(&graph) // simularea esuarii nodurilor

	nodeCnt := graph.N

	wg.Add(nodeCnt - graph.S)
	channelsIn := make([]chan Memory, 50)
	for i := 0; i < nodeCnt; i++ {
		channelsIn[i] = make(chan Memory, 50)
	}
	diam := diametru(graph.Adj, crashedNodes)
	fmt.Println("diametru: ", diam)
	iterations := (graph.S + 1) * diam

	leader := -1
	for i := range nodeCnt {
		if slices.Contains(crashedNodes, i) == false {
			if leader == -1 { // alegerea leader-ului, cel care va afisa mesajul
				leader = i
			}
			var channelsOut []chan<- Memory
			// cream canalele de iesire
			for _, neighborId := range graph.Adj[i] {
				channelsOut = append(channelsOut, channelsIn[neighborId])
			}
			initialMem := Memory{
				Info{ID: i, Value: graph.Tokens[i]}: {},
			}

			go goNode(Node{ID: i, Neighbors: graph.Adj[i], M: initialMem}, leader, nodeCnt-graph.S, iterations, inCount(i, graph), channelsIn[i], channelsOut, &wg)
		}
	}

	wg.Wait()

}

func main() {
	graph, err := ReadGraphConfig("input.txt")
	if err != nil {
		fmt.Println("Eroare:", err)
		os.Exit(1)
	}
	FloodSetAlgorithm(*graph)
}
