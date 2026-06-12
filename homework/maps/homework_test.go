package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Ordered interface {
	int8 | int16 | int32 | int64 | int |
		uint8 | uint16 | uint32 | uint64 | uint |
		~string
}

type node[T Ordered] struct {
	key   T
	left  *node[T]
	right *node[T]
	value any
}

type OrderedMap[T Ordered] struct {
	root *node[T]
	size int
}

func NewOrderedMap[T Ordered]() OrderedMap[T] {
	return OrderedMap[T]{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap[T]) Insert(key T, value any) {
	// Храним указатель на node, чтобы избежать лишней рекурсии
	current := &m.root

	for *current != nil {
		switch {
		case key < (*current).key:
			current = &(*current).left
		case key > (*current).key:
			current = &(*current).right
		default:
			// Ключ уже есть — перезаписываем значение
			(*current).value = value
			return
		}
	}

	*current = &node[T]{
		key:   key,
		value: value,
	}
	m.size++
}

func (m *OrderedMap[T]) Erase(key T) {
	current := &m.root
	for *current != nil && (*current).key != key {
		if key < (*current).key {
			current = &(*current).left
			continue
		}

		current = &(*current).right
	}

	// Ключ не найден
	if *current == nil {
		return
	}

	target := *current
	switch {
	case target.left == nil:
		*current = target.right
	case target.right == nil:
		*current = target.left
	default:
		// Два ребенка - ищем минимального в правом поддерве
		successors := &target.right
		for (*successors).left != nil {
			successors = &(*successors).left
		}
		target.key = (*successors).key
		target.value = (*successors).value
		// У successor'а нет левого ребёнка
		*successors = (*successors).right
	}
	m.size--
}

func (m *OrderedMap[T]) Contains(key T) bool {
	current := m.root

	for current != nil {
		switch {
		case key < current.key:
			current = current.left
		case key > current.key:
			current = current.right
		default:
			return true
		}
	}

	return false
}

func (m *OrderedMap[T]) Size() int {
	return m.size
}

func (m *OrderedMap[T]) ForEach(action func(T, any)) {
	var recursive func(n *node[T])

	recursive = func(n *node[T]) {
		if n == nil {
			return
		}

		recursive(n.left)
		action(n.key, n.value)
		recursive(n.right)
	}

	recursive(m.root)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key int, _ any) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
