package main

import (
	"fmt"
	"math/rand/v2"
	"time"
)

type Clock struct {
	offset time.Duration
}

func (c *Clock) Now() time.Time {
	return time.Now().Add(c.offset)
}

func estimate(t1, t2, t3, t4 time.Time) (theta, delta time.Duration) {
	theta = (t2.Sub(t1) + t3.Sub(t4)) / 2
	delta = t4.Sub(t1) + t2.Sub(t3)
	return theta, delta
}

func measure(c_clock, s_clock *Clock) (theta, delta time.Duration) {
	// Request
	t1 := c_clock.Now()
	time.Sleep(3*time.Millisecond + time.Duration(rand.IntN(10))*time.Millisecond)
	t2 := s_clock.Now()

	// Procesare intre 5 si 25 ms
	time.Sleep(5*time.Millisecond + time.Duration(rand.IntN(20))*time.Millisecond)

	t3 := s_clock.Now()

	// Response
	time.Sleep(3*time.Millisecond + time.Duration(rand.IntN(10))*time.Millisecond)
	t4 := c_clock.Now()
	return estimate(t1, t2, t3, t4)
}

func main() {
	// sereverul reprezinta timpul real, offset=0
	var s_clock Clock = Clock{offset: time.Duration(0)}
	// ceasul clientului este decalat cu +100ms
	var c_clock Clock = Clock{offset: 100 * time.Millisecond}

	fmt.Println("Offset client Initial: ", c_clock.offset)

	theta, delta := measure(&c_clock, &s_clock)

	// aplicam corectia asupra ceasului
	c_clock.offset += theta

	fmt.Println("theta", theta)
	fmt.Println("delta", delta)

	fmt.Println("Offset client dupa corectie: ", c_clock.offset)

	// Punctul B:

	// retinem minimul:
	var best_theta time.Duration
	var delta_min time.Duration = 1<<63 - 1

	// iteram de 10 ori
	for range 10 {
		// ceasul clientului este decalat cu +100ms
		c_clock = Clock{offset: 100 * time.Millisecond}

		theta, delta := measure(&c_clock, &s_clock)

		if delta < delta_min {
			delta_min = delta
			best_theta = theta
		}
	}

	fmt.Println("theta minim", best_theta)
	fmt.Println("delta minim", delta_min)

	// aplicam corectia asupra ceasului
	c_clock.offset += best_theta

	fmt.Println("Offset client minim dupa corectie: ", c_clock.offset)

}
