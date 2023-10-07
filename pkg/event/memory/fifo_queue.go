package memory

type fifoQueue[T any] struct {
	array []T
	head  int
	tail  int
}

func createFifoQueue[T any]() fifoQueue[T] {
	return fifoQueue[T]{
		array: make([]T, 4),
		head:  0,
		tail:  0,
	}
}

func (queue *fifoQueue[T]) enqueue(element T) {
	if queue.head == cap(queue.array) {
		queue.resize()
	}

	queue.array[queue.head] = element
	queue.head += 1
}

func (queue *fifoQueue[T]) dequeue() (T, bool) {
	var element T

	if queue.isEmpty() {
		return element, false
	}

	element = queue.array[queue.tail]
	queue.tail += 1

	return element, true
}

func (queue *fifoQueue[T]) count() int {
	return queue.head - queue.tail
}

func (queue *fifoQueue[T]) isEmpty() bool {
	return queue.count() == 0
}

func (queue *fifoQueue[T]) resize() {
	new_array := make([]T, len(queue.array)*2)
	copy(new_array, queue.array)
	queue.array = new_array
}
