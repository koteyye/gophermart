package queue

import (
	"context"
	"sync"
)

// FIFO определяет очередь по принципу FIFO.
//
// Структура потоко-безопасна.
type FIFO[T any] struct {
	once    sync.Once
	list    chan *linkedList[T]
	waiters chan struct{}
}

func (fifo *FIFO[T]) lazyInit() {
	fifo.once.Do(func() {
		fifo.list = make(chan *linkedList[T], 1)
		fifo.waiters = make(chan struct{}, 1)
		fifo.list <- new(linkedList[T])
	})
}

// notify уведомляет первого ждущего потребителя о наличии элементов в очереди.
func (fifo *FIFO[T]) notify() {
	select {
	case fifo.waiters <- struct{}{}:
	default:
	}
}

// Size возвращает количество элементов в очереди.
func (fifo *FIFO[T]) Size() int {
	fifo.lazyInit()

	list := <-fifo.list
	size := list.size
	fifo.list <- list

	return size
}

// Enqueue добавляет элемент в конец очереди.
func (fifo *FIFO[T]) Enqueue(ctx context.Context, value T) error {
	fifo.lazyInit()

	var list *linkedList[T]

	select {
	case <-ctx.Done():
		return ctx.Err()
	case list = <-fifo.list:
	}

	fifo.notify()

	list.Put(value)
	fifo.list <- list

	return nil
}

// Dequeue возвращает первый элемент из очереди.
func (fifo *FIFO[T]) Dequeue(ctx context.Context) (value T, err error) {
	fifo.lazyInit()

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	case <-fifo.waiters:
	}

	var list *linkedList[T]

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	case list = <-fifo.list:
	}

	value = list.Pop()

	if list.size > 0 {
		fifo.notify()
	}

	fifo.list <- list

	return value, nil
}

// linkedNode определяет узел односвязанного списка.
type linkedNode[T any] struct {
	value T
	next  *linkedNode[T]
}

// linkedList определяет односвязанный список.
type linkedList[T any] struct {
	head, tail *linkedNode[T]
	size       int
}

// Put добавляет элемент в конец односвязанного списка.
func (list *linkedList[T]) Put(value T) {
	node := &linkedNode[T]{value: value}
	switch {
	case list.head == nil:
		list.head = node
	case list.tail == nil:
		list.tail = node
		list.head.next = node
	default:
		list.tail.next = node
		list.tail = node
	}
	list.size++
}

// Pop возвращает первый элемент из односвязанного списка.
func (list *linkedList[T]) Pop() T {
	node := list.head
	list.head = node.next
	if list.head != nil && list.head == list.tail {
		list.head.next = nil
		list.tail = nil
	}
	list.size--
	return node.value
}
