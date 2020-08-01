package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	res int
	ready chan struct{}
}

type Memo struct {
	mu sync.Mutex
	cache map[int]*entry
	f func(int) int
}

func mockWork(key int) int {
	time.Sleep(200*time.Millisecond)
	return key * 2
}

func (m *Memo) Get (key int) int {
	m.mu.Lock()
	e, ok := m.cache[key]
	if !ok {
		e = &entry{ready: make(chan struct{})}

		// mark entry != nil
		// all subsequence goroutine will stuck waiting
		// for ready channel to be closed
		m.mu.Unlock()
		e.res = m.f(key)

		// signal
		close(e.ready)
	} else {
		m.mu.Unlock()

		// wait for entry to be ready
		<-e.ready
		return e.res
	}
	return e.res
}


func main() {
	memo := Memo{cache: make(map[int]*entry), f: mockWork}
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			wg.Add(1)
			i := i
			go func() {
				defer wg.Done()
				fmt.Println("Key ", i, " value ", memo.Get(i))
			}()
		}
	}

	wg.Wait()
	fmt.Println("elapsed in ", time.Since(start))
}