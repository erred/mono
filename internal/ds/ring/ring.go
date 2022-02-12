package ring

type Buffer[T any] struct {
	ring []T
	head int
	tail int
	full bool
}

func Of[T any](size int) *Buffer[T] {
	return &Buffer[T]{
		ring: make([]T, size),
	}
}

func (b *Buffer[T]) Put(item T) {
	b.ring[b.head] = item
	if b.full && (b.tail == b.head) {
		b.tail = (b.tail + 1) % len(b.ring)
	}
	b.head = (b.head + 1) % len(b.ring)
	if b.head == b.tail {
		b.full = true
	}
}

func (b *Buffer[T]) Get() (T, bool) {
	if !b.full && (b.tail == b.head) {
		var t T
		return t, false
	}
	t := b.ring[b.tail]
	b.tail = (b.tail + 1) % cap(b.ring)
	if b.tail == b.head {
		b.full = false
	}
	return t, true
}

func (b *Buffer[T]) Dump() []T {
	if b.tail < b.head {
		out := make([]T, b.head-b.tail)
		copy(out, b.ring[b.tail:b.head])
		return out
	} else if b.tail == b.head && !b.full {
		return nil
	}
	out := make([]T, (len(b.ring)-b.tail)+b.head)
	copy(out, b.ring[b.tail:])
	copy(out[len(b.ring)-b.tail:], b.ring[:b.head])
	return out
}
