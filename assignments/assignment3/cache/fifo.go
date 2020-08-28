package cache

import (
	"container/list"
)

// An FIFO is a fixed-size in-memory cache with first-in first-out eviction
type FIFO struct {
	// whatever fields you want here
	Limit           int
	OrderedElements *list.List
	ElementsMap     map[string]*list.Element
	CurrentStorage  int
	Statistics      *Stats
}

type ListElement struct {
	Value []byte
	Key   string
}

// NewFIFO returns a pointer to a new FIFO with a capacity to store limit bytes
func NewFifo(limit int) *FIFO {
	return &FIFO{Limit: limit, OrderedElements: list.New(), ElementsMap: make(map[string]*list.Element), Statistics: &Stats{}}
}

// MaxStorage returns the maximum number of bytes this FIFO can store
func (fifo *FIFO) MaxStorage() int {
	return fifo.Limit
}

// RemainingStorage returns the number of unused bytes available in this FIFO
func (fifo *FIFO) RemainingStorage() int {
	return fifo.Limit - fifo.CurrentStorage
}

// Get returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise.
func (fifo *FIFO) Get(key string) (value []byte, ok bool) {
	listElement, ok := fifo.ElementsMap[key]
	if !ok {
		fifo.Statistics.Misses += 1
		return nil, false
	} else {
		fifo.Statistics.Hits += 1
		return listElement.Value.(ListElement).Value, true
	}
}

// Remove removes and returns the value associated with the given key, if it exists.
// ok is true if a value was found and false otherwise
func (fifo *FIFO) Remove(key string) (value []byte, ok bool) {
	listElement, ok := fifo.ElementsMap[key]
	if !ok {
		return nil, false
	} else {
		fifo.CurrentStorage -= len(listElement.Value.(ListElement).Key) - len(listElement.Value.(ListElement).Value)
		fifo.OrderedElements.Remove(listElement)
		delete(fifo.ElementsMap, key)
		return listElement.Value.(ListElement).Value, true
	}
}

// Set associates the given value with the given key, possibly evicting values
// to make room. Returns true if the binding was added successfully, else false.
func (fifo *FIFO) Set(key string, value []byte) bool {
	sizeOfInsert := len(key) + len(value)
	if fifo.Limit < sizeOfInsert {
		return false
	}
	// pop elements one at time until we have enough storage
	for (fifo.RemainingStorage() - sizeOfInsert) < 0 {
		frontElement := fifo.OrderedElements.Front()
		fifo.OrderedElements.Remove(frontElement)
		fifo.CurrentStorage -= len(frontElement.Value.(ListElement).Key) - len(frontElement.Value.(ListElement).Value)
		delete(fifo.ElementsMap, frontElement.Value.(ListElement).Key)
	}

	listElement := ListElement{Value: value, Key: key}
	fifo.ElementsMap[key] = fifo.OrderedElements.PushBack(listElement)
	fifo.CurrentStorage += sizeOfInsert
	return true
}

// Len returns the number of bindings in the FIFO.
func (fifo *FIFO) Len() int {
	return fifo.OrderedElements.Len()
}

// Stats returns statistics about how many search hits and misses have occurred.
func (fifo *FIFO) Stats() *Stats {
	return fifo.Statistics
}
