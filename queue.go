package gsync

import (
	"errors"
	"sync"
)

var (
	ErrQueueEmpty = errors.New("queue is empty")
	ErrQueueLess  = errors.New("queue len less than n")
)

type Queue[T any] struct {
	queue []T
	rw    sync.RWMutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{queue: make([]T, 0)}
}

func (q *Queue[T]) Push(v T) {
	q.rw.Lock()
	defer q.rw.Unlock()
	q.queue = append(q.queue, v)
}

func (q *Queue[T]) Pop() (T, error) {
	q.rw.Lock()
	defer q.rw.Unlock()
	if len(q.queue) == 0 {
		return GetZero[T](), ErrQueueEmpty
	}
	v := q.queue[0]
	q.queue = q.queue[1:]
	return v, nil
}

func (q *Queue[T]) MustPop() T {
	v, err := q.Pop()
	if err != nil {
		panic(err)
	}
	return v
}

func (q *Queue[T]) Len() int {
	q.rw.RLock()
	defer q.rw.RUnlock()
	return len(q.queue)
}

func (q *Queue[T]) IsEmpty() bool {
	return q.Len() == 0
}

func (q *Queue[T]) Clear() {
	q.rw.Lock()
	defer q.rw.Unlock()
	q.queue = nil
}

func (q *Queue[T]) Peek() (T, error) {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if len(q.queue) == 0 {
		return GetZero[T](), ErrQueueEmpty
	}
	return q.queue[0], nil
}

func (q *Queue[T]) MustPeek() T {
	v, err := q.Peek()
	if err != nil {
		panic(err)
	}
	return v
}

func (q *Queue[T]) PeekN(n int) ([]T, error) {
	q.rw.RLock()
	defer q.rw.RUnlock()
	if len(q.queue) < n {
		return nil, ErrQueueLess
	}
	return q.queue[:n], nil
}

func (q *Queue[T]) MustPeekN(n int) []T {
	v, err := q.PeekN(n)
	if err != nil {
		panic(err)
	}
	return v
}

func (q *Queue[T]) PopN(n int) ([]T, error) {
	q.rw.Lock()
	defer q.rw.Unlock()
	if len(q.queue) < n {
		return nil, ErrQueueLess
	}
	v := q.queue[:n]
	q.queue = q.queue[n:]
	return v, nil
}

func (q *Queue[T]) MustPopN(n int) []T {
	v, err := q.PopN(n)
	if err != nil {
		panic(err)
	}
	return v
}

func (q *Queue[T]) PopAll() []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	v := q.queue
	q.queue = nil
	return v
}

func (q *Queue[T]) PopWhile(f func(T) bool) []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	var i int
	for i = 0; i < len(q.queue); i++ {
		if !f(q.queue[i]) {
			break
		}
	}
	v := q.queue[:i]
	q.queue = q.queue[i:]
	return v
}

func (q *Queue[T]) PopUntil(f func(T) bool) []T {
	q.rw.Lock()
	defer q.rw.Unlock()
	var i int
	for i = 0; i < len(q.queue); i++ {
		if f(q.queue[i]) {
			break
		}
	}
	v := q.queue[:i]
	q.queue = q.queue[i:]
	return v
}
