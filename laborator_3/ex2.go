package main

import (
	"fmt"
	"sort"
	"sync"
)

type LamportClock struct{ counter int }

// incrementeaza contorul

func (l *LamportClock) Tick() {
	l.counter++
}
func (l *LamportClock) Send() {
	l.counter++
}
func (l *LamportClock) Receive(ts int) {
	l.counter = max(l.counter, ts)
	l.counter++
}

type Message struct{ sender, ts int }
type Event struct {
	process          int
	internal_counter int
	kind             string // "internal", "send", "recv"
}

// TODO: retuneaza true daca a precede b
func less(a, b Event) bool {
	if a.internal_counter != b.internal_counter {
		return a.internal_counter < b.internal_counter
	}
	return a.process < b.process
}

func main() {
	var ch12 = make(chan Message)
	var ch23 = make(chan Message)
	var ch31 = make(chan Message)

	eventsChan := make(chan Event, 100)
	var wg sync.WaitGroup

	// ==========================================
	// Procesul P1
	// ==========================================
	wg.Add(1)
	go func() {
		defer wg.Done()
		clock := LamportClock{}

		// 1. Eveniment intern
		clock.Tick()
		eventsChan <- Event{process: 1, internal_counter: clock.counter, kind: "internal"}

		// 2. Trimitere mesaj catre P2
		clock.Send()
		eventsChan <- Event{process: 1, internal_counter: clock.counter, kind: "send"}
		ch12 <- Message{sender: 1, ts: clock.counter}

		// 3. Primire mesaj inapoi de la P3
		msg := <-ch31
		clock.Receive(msg.ts)
		eventsChan <- Event{process: 1, internal_counter: clock.counter, kind: "recv"}
	}()

	// ==========================================
	// Procesul P2
	// ==========================================
	wg.Add(1)
	go func() {
		defer wg.Done()
		clock := LamportClock{}

		// 1. Primire mesaj de la P1
		msg := <-ch12
		clock.Receive(msg.ts)
		eventsChan <- Event{process: 2, internal_counter: clock.counter, kind: "recv"}

		// 2. Eveniment intern
		clock.Tick()
		eventsChan <- Event{process: 2, internal_counter: clock.counter, kind: "internal"}

		// 3. Trimitere mesaj catre P3
		clock.Send()
		eventsChan <- Event{process: 2, internal_counter: clock.counter, kind: "send"}
		ch23 <- Message{sender: 2, ts: clock.counter}
	}()

	// ==========================================
	// Procesul P3
	// ==========================================
	wg.Add(1)
	go func() {
		defer wg.Done()
		clock := LamportClock{}

		// 1. Primire mesaj de la P2
		msg := <-ch23
		clock.Receive(msg.ts)
		eventsChan <- Event{process: 3, internal_counter: clock.counter, kind: "recv"}

		// 2. Eveniment intern
		clock.Tick()
		eventsChan <- Event{process: 3, internal_counter: clock.counter, kind: "internal"}

		// 3. Trimitere mesaj la P1
		clock.Send()
		eventsChan <- Event{process: 3, internal_counter: clock.counter, kind: "send"}
		ch31 <- Message{sender: 3, ts: clock.counter}
	}()

	wg.Wait()
	close(eventsChan)

	var events []Event
	for ev := range eventsChan {
		events = append(events, ev)
	}

	sort.Slice(events, func(i, j int) bool {
		return less(events[i], events[j])
	})

	fmt.Printf("%-10s | %-16s | %-10s\n", "Proces", "Contor Lamport", "Tip Eveniment")
	fmt.Println("----------------------------------------------")
	for _, e := range events {
		fmt.Printf("P%-9d | %-16d | %-10s\n", e.process, e.internal_counter, e.kind)
	}
}
