package memory

import (
	"math/rand"
	"testing"
)

func Test_Sequence1(t *testing.T) {
	queue := createFifoQueue[int]()
	queue.enqueue(0)
	queue.enqueue(1)
	d1, e1 := queue.dequeue()
	queue.enqueue(2)
	queue.enqueue(3)
	d2, e2 := queue.dequeue()
	d3, e3 := queue.dequeue()
	queue.enqueue(4)
	d4, e4 := queue.dequeue()
	d5, e5 := queue.dequeue()
	_, e6 := queue.dequeue()

	if !e1 {
		t.Error("Dequeue 1 did not exists")
	}
	if d1 != 0 {
		t.Error("Dequeue 1 was", d1, "and not", 0)
	}

	if !e2 {
		t.Error("Dequeue 2 did not exists")
	}
	if d2 != 1 {
		t.Error("Dequeue 2 was", d2, "and not", 1)
	}

	if !e3 {
		t.Error("Dequeue 3 did not exists")
	}
	if d3 != 2 {
		t.Error("Dequeue 3 was", d3, "and not", 2)
	}

	if !e4 {
		t.Error("Dequeue 4 did not exists")
	}
	if d4 != 3 {
		t.Error("Dequeue 4 was", d4, "and not", 3)
	}

	if !e5 {
		t.Error("Dequeue 5 did not exists")
	}
	if d5 != 4 {
		t.Error("Dequeue 5 was", d4, "and not", 4)
	}

	if e6 {
		t.Error("Dequeue 6 did exist")
	}

	if !queue.isEmpty() {
		t.Error("Queue is not empty but has", queue.count(), "elements")
	}
}

func Test_Enqueue(t *testing.T) {
	tests := []struct {
		name   string
		amount int
	}{
		{
			name:   "Zero enqueues zero count",
			amount: 0,
		},
		{
			name:   "One enqueue one count",
			amount: 1,
		},
		{
			name:   "Five enqueues five count",
			amount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := createFifoQueue[int]()
			count := 0

			for count < tt.amount {
				queue.enqueue(0)
				count += 1
			}

			if tt.amount != queue.count() {
				t.Error("Performed", count, "enqueues, queue count is", queue.count())
			}
		})
	}
}

func Test_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		enqueues int
		dequeues int
		expected bool
	}{
		{
			name:     "Initial queue is empty",
			enqueues: 0,
			dequeues: 0,
			expected: true,
		},
		{
			name:     "Dequeue and enqueue once results in an empty queue",
			enqueues: 1,
			dequeues: 1,
			expected: true,
		},
		{
			name:     "Dequeue and enqueue more than initial capacity results in an empty queue",
			enqueues: 10,
			dequeues: 10,
			expected: true,
		},
		{
			name:     "Dequeueing an empty queue is still empty",
			enqueues: 0,
			dequeues: 1,
			expected: true,
		},
		{
			name:     "Enqueuing once and no dequeue does not result in an empty queue",
			enqueues: 1,
			dequeues: 0,
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := createFifoQueue[int]()

			performedEnqueues := 0
			for performedEnqueues < tt.enqueues {
				queue.enqueue(0)
				performedEnqueues += 1
			}

			performedDequeues := 0
			for performedDequeues < tt.dequeues {
				queue.dequeue()
				performedDequeues += 1
			}

			if queue.isEmpty() != tt.expected {
				if tt.expected == true {
					t.Error("Queue was not empty - head:", queue.head, "tail:", queue.tail)
				} else {
					t.Error("Queue was empty - head:", queue.head, "tail:", queue.tail)
				}
			}
		})
	}
}

func Test_Count(t *testing.T) {
	tests := []struct {
		name     string
		enqueues int
		dequeues int
		expected int
	}{
		{
			name:     "Equal amount of enqueues and dequese count should be zero",
			enqueues: 10,
			dequeues: 5,
			expected: 5,
		},
		{
			name:     "Equal huge amount of enqueues and dequese count should be zero",
			enqueues: 1000,
			dequeues: 1000,
			expected: 0,
		},
		{
			name:     "Unequal huge amount of enqueues and dequese count should be 900",
			enqueues: 100000,
			dequeues: 10000,
			expected: 90000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := createFifoQueue[int]()

			performedEnqueues := 0
			performedDequeues := 0

			for {
				missesEnqueues := performedEnqueues < tt.enqueues
				missesDequeues := performedDequeues < tt.dequeues
				queueShouldBeEmpty := performedEnqueues == performedDequeues

				if !missesDequeues && !missesEnqueues {
					break
				}

				shouldEnqueue := missesEnqueues && (rand.Intn(2) == 0 || !missesDequeues || queueShouldBeEmpty)

				if shouldEnqueue {
					queue.enqueue(0)
					performedEnqueues += 1
				} else {
					queue.dequeue()
					performedDequeues += 1
				}
			}

			if queue.count() != tt.expected {
				t.Error("Queue count was", queue.count(), "and not", tt.expected)
			}
		})
	}
}
