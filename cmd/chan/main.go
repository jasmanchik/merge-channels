package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-mergeChannels(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(6*time.Second),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))

}

func mergeChannels(channels ...<-chan any) <-chan any {

	var once sync.Once
	s := make(chan any)
	for _, ch := range channels {
		go func(ch <-chan any, s chan any) {
			for {
				val, ok := <-ch
				if !ok {
					once.Do(func() {
						close(s)
					})
					return
				}
				s <- val
			}
		}(ch, s)
	}

	return s
}
