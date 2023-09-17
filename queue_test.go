package gsync

import (
	"sync"
	"testing"
)

func TestQueue_Push(t *testing.T) {
	q := NewQueue[int]()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				q.Push(i)
			}
		}()
	}
	wg.Wait()
	if q.Len() != 1000 {
		t.Error("queue len error")
	}
}

func TestQueue_Pop(t *testing.T) {
	q := NewQueue[int]()
	for i := 0; i < 1000; i++ {
		q.Push(i)
	}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				q.MustPop()
			}
		}()
	}
	wg.Wait()
	if q.Len() != 0 {
		t.Error("queue len error")
	}
}
