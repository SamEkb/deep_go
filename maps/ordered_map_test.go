package maps

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type node struct {
	key   int
	value int
	left  *node
	right *node
}

type OrderedMap struct {
	root *node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root: nil,
		size: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.size++
		m.root = &node{key: key, value: value}
		return
	}

	nodeToInsert := &node{key: key, value: value}

	m.insert(m.root, nodeToInsert)
}

func (m *OrderedMap) insert(root, nodeToInsert *node) *node {
	if root == nil {
		m.size++
		return nodeToInsert
	}

	if root.key > nodeToInsert.key {
		root.left = m.insert(root.left, nodeToInsert)
	} else {
		root.right = m.insert(root.right, nodeToInsert)
	}

	return root
}

func (m *OrderedMap) Erase(key int) {
	m.root = m.erase(key, m.root)
}

func (m *OrderedMap) erase(key int, root *node) *node {
	if root == nil {
		return nil
	}

	if key < root.key {
		root.left = m.erase(key, root.left)
	} else if key > root.key {
		root.right = m.erase(key, root.right)
	} else {
		m.size--

		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		}

		minNode := m.findMin(root.right)
		root.key = minNode.key
		root.value = minNode.value
		root.right = m.erase(minNode.key, root.right)
	}

	return root
}

func (m *OrderedMap) findMin(root *node) *node {
	for root.left != nil {
		root = root.left
	}
	return root
}

func (m *OrderedMap) Contains(key int) bool {
	return m.contains(key, m.root)
}

func (m *OrderedMap) contains(key int, root *node) bool {
	if root == nil {
		return false
	}

	if root.key == key {
		return true
	}

	if key < root.key {
		return m.contains(key, root.left)
	}

	return m.contains(key, root.right)
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.forEach(m.root, action)
}

func (m *OrderedMap) forEach(root *node, action func(int, int)) {
	if root == nil {
		return
	}

	m.forEach(root.left, action)
	action(root.key, root.value)
	m.forEach(root.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
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
	data.ForEach(func(key, _ int) {
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
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
