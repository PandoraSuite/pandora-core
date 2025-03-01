package main

import (
	"sync"
	"time"
)

func main() {
	t := new(test)

	go t.function1()
	go t.function2()

	time.Sleep(time.Minute * 3)
}

type test struct {
	mu sync.Mutex
}

func (t *test) function1() {
	println("Start test function1")
	t.mu.Lock()
	println("middle test function1")
	time.Sleep(time.Second * 10)
	t.mu.Unlock()
	println("end test function1")
}

func (t *test) function2() {
	time.Sleep(time.Second * 2)
	println("Start test function2")
	t.mu.Lock()
	println("middle test function2")
	time.Sleep(time.Second * 120)
	t.mu.Unlock()
	println("end test function2")
}
