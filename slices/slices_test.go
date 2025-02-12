package slices

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

// go test -v slices_test.go

type CircularQueue struct {
	values []int
	size   int
	front  int
	rear   int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{values: make([]int, size)}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	q.values[q.rear] = value
	q.rear = (q.rear + 1) % cap(q.values)
	q.size++
	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	q.front = (q.front + 1) % cap(q.values)
	q.size--
	return true
}

func (q *CircularQueue) Front() int {
	if !q.Empty() {
		return q.values[q.front]
	}
	return -1
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	if q.rear == 0 {
		return q.values[cap(q.values)-1]
	}
	return q.values[q.rear-1]
}

func (q *CircularQueue) Empty() bool {
	return q.size == 0
}

func (q *CircularQueue) Full() bool {
	return q.size == cap(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
